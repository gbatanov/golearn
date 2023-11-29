package winapi

import (
	"errors"
	"image"
	"log"
	"sync"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

type Stage uint8

const (
	// StagePaused is the stage for windows that have no on-screen representation.
	// Paused windows don't receive FrameEvent.
	StagePaused Stage = iota
	// StageInactive is the stage for windows that are visible, but not active.
	// Inactive windows receive FrameEvent.
	StageInactive
	// StageRunning is for active and visible
	// Running windows receive FrameEvent.
	StageRunning
)

// winMap maps win32 HWNDs to *
var winMap sync.Map

type WindowMode uint8

const (
	// Windowed is the normal window mode with OS specific window decorations.
	Windowed WindowMode = iota
	// Fullscreen is the full screen window mode.
	Fullscreen
	// Minimized is for systems where the window can be minimized to an icon.
	Minimized
	// Maximized is for systems where the window can be made to fill the available monitor area.
	Maximized
)

type Config struct {
	Position image.Point
	Size     image.Point
	MinSize  image.Point
	MaxSize  image.Point
	Mode     WindowMode
	//	Decorated bool
	Title      string
	EventChan  chan Event
	BorderSize image.Point
}

type Window struct {
	Id          int32 // 0 у главного, начиная с 1 у дочерних
	Hwnd        syscall.Handle
	Hdc         syscall.Handle
	HInst       syscall.Handle
	Focused     bool
	Stage       Stage
	Config      *Config
	Cursor      syscall.Handle
	PointerBtns Buttons //Кнопки мыши
	Parent      *Window
	Childrens   []*Window
	// cursorIn tracks whether the cursor was inside the window according
	// to the most recent WM_SETCURSOR.
	CursorIn bool
}

// iconID is the ID of the icon in the resource file.
const iconID = 1

var resources struct {
	once sync.Once
	// handle is the module handle from GetModuleHandle.
	handle syscall.Handle
	// class is window class from RegisterClassEx.
	class string
	// cursor is the arrow cursor resource.
	cursor syscall.Handle
}

// initResources initializes the resources global.
func initResources(child bool) error {
	SetProcessDPIAware()
	hInst, err := GetModuleHandle()
	if err != nil {
		return err
	}

	c, err := LoadCursor(IDC_ARROW)
	if err != nil {
		return err
	}

	icon, _ := LoadImage(hInst, iconID, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE|LR_SHARED)
	wcls := WndClassEx{
		CbSize:    uint32(unsafe.Sizeof(WndClassEx{})),
		HInstance: hInst,
	}

	wcls.Style = CS_HREDRAW | CS_VREDRAW | CS_OWNDC
	wcls.HIcon = icon
	wcls.LpszClassName = syscall.StringToUTF16Ptr("GsbWindow")

	wcls.LpfnWndProc = syscall.NewCallback(windowProc)
	_, err = RegisterClassEx(&wcls)
	if err != nil {
		return err
	}
	resources.handle = hInst
	resources.cursor = c
	resources.class = "GsbWindow"

	return nil
}

const dwExStyle = WS_EX_APPWINDOW | WS_EX_WINDOWEDGE

// Создание основного окна программы
func CreateNativeMainWindow(config *Config) error {

	var resErr error
	resources.once.Do(func() {
		resErr = initResources(false)
	})
	if resErr != nil {
		return resErr
	}
	// WS_CAPTION включает в себя WS_BORDER
	var dwStyle uint32 = 0 | WS_CAPTION | WS_SYSMENU //WS_THICKFRAME
	//	if config.Decorated {
	//		dwStyle = dwStyle | WS_SYSMENU | WS_CLIPSIBLINGS | WS_CLIPCHILDREN
	//	}

	hwnd, err := CreateWindowEx(
		dwExStyle,
		"GsbWindow",                                        //	resourceMain.class,                                 //lpClassame
		config.Title,                                       // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		0,                //hWndParent
		0,                // hMenu
		resources.handle, //hInstance
		0)                // lpParam
	if err != nil {
		return err
	}
	w := &Window{
		Id:        0,
		Hwnd:      hwnd,
		HInst:     resources.handle,
		Config:    config,
		Parent:    nil,
		Childrens: make([]*Window, 0),
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return err
	}

	winMap.Store(w.Hwnd, w)
	defer winMap.Delete(w.Hwnd)

	childConfig := &Config{
		Title:      "Child",
		EventChan:  w.Config.EventChan,
		Size:       image.Pt(int(80), int(40)),
		MinSize:    w.Config.MinSize,
		MaxSize:    w.Config.MaxSize,
		Position:   image.Pt(int(10), int(10)),
		Mode:       Windowed,
		BorderSize: image.Pt(0, 0),
	}
	chWin, err := CreateChildWindow(w, childConfig)
	if err != nil {
		panic(err.Error())
	}

	winMap.Store(chWin.Hwnd, chWin)
	defer winMap.Delete(chWin.Hwnd)

	SetForegroundWindow(w.Hwnd)
	SetFocus(w.Hwnd)
	w.SetCursor(CursorDefault)
	ShowWindow(w.Hwnd, SW_SHOWNORMAL)

	msg := new(Msg)
	for {
		ret := GetMessage(msg, 0, 0, 0)
		switch ret {
		case -1:
			return errors.New("GetMessage failed")
		case 0:
			// WM_QUIT received.
			return nil
		}

		TranslateMessage(msg)
		DispatchMessage(msg)
	}

}

func coordsFromlParam(lParam uintptr) (int, int) {
	x := int(int16(lParam & 0xffff))
	y := int(int16((lParam >> 16) & 0xffff))
	return x, y
}

func (w *Window) draw(sync bool) {
	if w.Config.Size.X == 0 || w.Config.Size.Y == 0 {
		return
	}

	// Fill the region
	if w.Id == 0 {
		r2 := GetClientRect(w.Hwnd)
		r2.Left += 120
		r2.Top += 20
		r2.Bottom -= 20
		r2.Right -= 20

		FillRect(w.Hdc, &r2, GetStockObject(1)) // 0,5-белый, 1 - серый, 2-темно-серый, 4 - черный
	}

}

// update() handles changes done by the user, and updates the configuration.
// It reads the window style and size/position and updates w.config.
// If anything has changed it emits a ConfigEvent to notify the application.
func (w *Window) update() {

	cr := GetClientRect(w.Hwnd)
	w.Config.Size = image.Point{
		X: int(cr.Right - cr.Left),
		Y: int(cr.Bottom - cr.Top),
	}

	w.Config.BorderSize = image.Pt(
		GetSystemMetrics(SM_CXSIZEFRAME),
		GetSystemMetrics(SM_CYSIZEFRAME),
	)
	//		w.w.Event(ConfigEvent{Config: w.config})

}

func (w *Window) SetCursor(cursor Cursor) {
	c, err := loadCursor(cursor)
	if err != nil {
		c = resources.cursor
	}
	w.Cursor = c
	SetCursor(w.Cursor) // Win32 API function
}

func loadCursor(cursor Cursor) (syscall.Handle, error) {
	switch cursor {
	case CursorDefault:
		return resources.cursor, nil
	case CursorNone:
		return 0, nil
	default:
		return LoadCursor(windowsCursor[cursor])
	}
}

func (w *Window) pointerButton(btn Buttons, press bool, lParam uintptr, kmods Modifiers) {
	if !w.Focused {
		SetFocus(w.Hwnd)
	}
	log.Println("pointerButton", btn, press)
	var kind Kind
	if press {
		kind = Press
		if w.PointerBtns == 0 {
			SetCapture(w.Hwnd) // Захват событий мыши окном
		}
		w.PointerBtns |= btn
	} else {
		kind = Release
		w.PointerBtns &^= btn
		if w.PointerBtns == 0 {
			ReleaseCapture() // Освобождение событий мыши окном
		}
	}

	x, y := coordsFromlParam(lParam)
	p := image.Point{X: (x), Y: (y)}
	w.Config.EventChan <- Event{
		SWin:      w,
		Kind:      kind,
		Source:    Mouse,
		Position:  p,
		Buttons:   w.PointerBtns,
		Time:      GetMessageTime(),
		Modifiers: kmods,
	}

}

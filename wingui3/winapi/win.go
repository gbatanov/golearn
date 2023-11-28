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
	Position  image.Point
	Size      image.Point
	MinSize   image.Point
	MaxSize   image.Point
	Mode      WindowMode
	Decorated bool
	Title     string
	Wg        *sync.WaitGroup
	EventChan chan Event
}

type Window struct {
	Hwnd        syscall.Handle
	Hdc         syscall.Handle
	Focused     bool
	Stage       Stage
	Config      *Config
	BorderSize  image.Point
	Cursor      syscall.Handle
	PointerBtns Buttons //Кнопки мыши
	Parent      *Window
	Childrens   []*Window
}

// iconID is the ID of the icon in the resource file.
const iconID = 1

var resources struct {
	once sync.Once
	// handle is the module handle from GetModuleHandle.
	handle syscall.Handle
	// class is window class from RegisterClassEx.
	class uint16
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
	resources.handle = hInst
	c, err := LoadCursor(IDC_ARROW)
	if err != nil {
		return err
	}
	resources.cursor = c
	icon, _ := LoadImage(hInst, iconID, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE|LR_SHARED)
	wcls := WndClassEx{
		CbSize:      uint32(unsafe.Sizeof(WndClassEx{})),
		Style:       CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
		LpfnWndProc: syscall.NewCallback(windowProc),
		HInstance:   hInst,
		HIcon:       icon,
		//		LpszClassName: syscall.StringToUTF16Ptr("GsbWindow"),
	}
	if child {
		wcls.LpszClassName = syscall.StringToUTF16Ptr("GsbChildWindow")
	} else {
		wcls.LpszClassName = syscall.StringToUTF16Ptr("GsbWindow")
	}
	cls, err := RegisterClassEx(&wcls)
	if err != nil {
		return err
	}
	resources.class = cls
	return nil
}

const dwExStyle = WS_EX_APPWINDOW | WS_EX_WINDOWEDGE

// Создание основоного окна программы
func CreateNativeMainWindow(config *Config) error {

	var resErr error
	resources.once.Do(func() {
		resErr = initResources(false)
	})
	if resErr != nil {
		return resErr
	}
	const dwStyle = WS_OVERLAPPED | WS_CAPTION | WS_THICKFRAME | WS_SYSMENU

	hwnd, err := CreateWindowEx(
		dwExStyle,
		resources.class, //lpClassame
		config.Title,    // lpWindowName
		dwStyle|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,            //dwStyle
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
		Hwnd:      hwnd,
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
	SetForegroundWindow(w.Hwnd)
	SetFocus(w.Hwnd)
	// Since the window class for the cursor is null,
	// set it here to show the cursor.
	w.SetCursor(CursorDefault)

	_, err = CreateChildWindow(w, 10, 10, 80, 40)
	if err != nil {
		log.Println(err)
	}

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
	return nil
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
	/*
		dpi := GetWindowDPI(w.Hwnd)
		cfg := configForDPI(dpi)
		w.w.Event(frameEvent{
			FrameEvent: system.FrameEvent{
				Now:    time.Now(),
				Size:   w.config.Size,
				Metric: cfg,
			},
			Sync: sync,
		})
	*/
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

	w.BorderSize = image.Pt(
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

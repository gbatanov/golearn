package winapi

import (
	"image"
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

// winMap maps win32 HWNDs to *Window
var WinMap sync.Map

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
	ID         uintptr
	Position   image.Point
	Size       image.Point
	MinSize    image.Point
	MaxSize    image.Point
	Mode       WindowMode
	SysMenu    bool
	Title      string
	EventChan  chan Event
	BorderSize image.Point
	TextColor  uint32
	FontSize   int32
	BgColor    uint32
	Class      string
}

type Window struct {
	Hwnd        syscall.Handle
	Hdc         syscall.Handle
	HInst       syscall.Handle
	Focused     bool
	Stage       Stage
	Config      Config
	Cursor      syscall.Handle
	PointerBtns Buttons //Кнопки мыши
	Parent      *Window
	Childrens   map[int]*Window
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
func initResources() error {
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

const dwExStyle = 0 //| WS_EX_WINDOWEDGE // WS_EX_APPWINDOW

// Создание основного окна программы
func CreateNativeMainWindow(config Config) (*Window, error) {

	var resErr error
	resources.once.Do(func() {
		resErr = initResources()
	})
	if resErr != nil {
		return nil, resErr
	}
	// WS_CAPTION включает в себя WS_BORDER
	var dwStyle uint32 = 0
	if config.SysMenu {
		dwStyle = dwStyle | WS_SYSMENU | WS_CAPTION | WS_SIZEBOX
	} else {
		dwStyle = dwStyle | WS_POPUP
	}

	hwnd, err := CreateWindowEx(
		dwExStyle,
		config.Class,                                       //	resourceMain.class,                                 //lpClassame
		config.Title,                                       // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		0,                //hWndParent
		0,                // hMenu
		resources.handle, //hInstance
		0)                // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Hwnd:      hwnd,
		HInst:     resources.handle,
		Config:    config,
		Parent:    nil,
		Childrens: make(map[int]*Window, 0),
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}

	WinMap.Store(w.Hwnd, w)

	SetForegroundWindow(w.Hwnd)
	SetFocus(w.Hwnd)
	w.SetCursor(CursorDefault)
	ShowWindow(w.Hwnd, SW_SHOWNORMAL)
	return w, nil
}

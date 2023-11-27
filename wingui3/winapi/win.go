package winapi

import (
	"errors"
	"fmt"
	"image"
	"sync"
	"unicode"
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
	Size      image.Point
	MinSize   image.Point
	MaxSize   image.Point
	Mode      WindowMode
	Decorated bool
	Title     string
	Wg        *sync.WaitGroup
}

type Window struct {
	Hwnd       syscall.Handle
	Hdc        syscall.Handle
	focused    bool
	stage      Stage
	config     Config
	borderSize image.Point
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
func initResources() error {
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
		CbSize:        uint32(unsafe.Sizeof(WndClassEx{})),
		Style:         CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
		LpfnWndProc:   syscall.NewCallback(windowProc),
		HInstance:     hInst,
		HIcon:         icon,
		LpszClassName: syscall.StringToUTF16Ptr("GsbWindow"),
	}
	cls, err := RegisterClassEx(&wcls)
	if err != nil {
		return err
	}
	resources.class = cls
	return nil
}

const dwExStyle = WS_EX_APPWINDOW | WS_EX_WINDOWEDGE

func CreateNativeWindow(config Config) (*Window, error) {

	var resErr error
	resources.once.Do(func() {
		resErr = initResources()
	})
	if resErr != nil {
		return nil, resErr
	}
	const dwStyle = WS_OVERLAPPEDWINDOW

	hwnd, err := CreateWindowEx(
		dwExStyle,
		resources.class,              //lpClassame
		config.Title,                 // lpWindowName
		dwStyle,                      //WS_CLIPSIBLINGS|WS_CLIPCHILDREN, //dwStyle
		CW_USEDEFAULT, CW_USEDEFAULT, //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		0,                //hWndParent
		0,                // hMenu
		resources.handle, //hInstance
		0)                // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Hwnd:   hwnd,
		config: config,
	}
	winMap.Store(w.Hwnd, w)
	defer winMap.Delete(w.Hwnd)
	SetForegroundWindow(w.Hwnd)
	SetFocus(w.Hwnd)
	// Since the window class for the cursor is null,
	// set it here to show the cursor.
	//		win.SetCursor(pointer.CursorDefault)
	ShowWindow(w.Hwnd, SW_SHOWNORMAL)

	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}
	if err := w.Loop(); err != nil {
		panic(err)
	}

	return w, nil
}

// Adapted from https://blogs.msdn.microsoft.com/oldnewthing/20060126-00/?p=32513/
func (w *Window) Loop() error {
	msg := new(Msg)
loop:
	for {
		//		fmt.Println(msg)
		anim := false // w.animating
		if anim && !PeekMessage(msg, 0, 0, 0, PM_NOREMOVE) {
			w.draw(false)
			continue
		}
		ret := GetMessage(msg, 0, 0, 0)
		//		fmt.Println(ret)
		switch ret {
		case -1:
			return errors.New("GetMessage failed")
		case 0:
			// WM_QUIT received.
			break loop
		}
		TranslateMessage(msg)
		DispatchMessage(msg)
	}
	return nil
}

func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	win, exists := winMap.Load(hwnd)
	if !exists {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	w := win.(*Window)

	switch msg {
	case WM_UNICHAR:
		if wParam == UNICODE_NOCHAR {
			// Tell the system that we accept WM_UNICHAR messages.
			return TRUE
		}
		fallthrough
	case WM_CHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			//			w.w.EditorInsert(string(r))
		}
		// The message is processed.
		return TRUE
	case WM_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return TRUE
	case WM_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return TRUE
	case WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP:
		/*
			if n, ok := convertKeyCode(wParam); ok {
				e := Event{
					Name: n,
					//				Modifiers: getModifiers(),
					State: Press,
				}
				if msg == WM_KEYUP || msg == WM_SYSKEYUP {
					e.State = Release
				}

				//			w.w.Event(e)

				if (wParam == VK_F10) && (msg == WM_SYSKEYDOWN || msg == WM_SYSKEYUP) {
					// Reserve F10 for ourselves, and don't let it open the system menu. Other Windows programs
					// such as cmd.exe and graphical debuggers also reserve F10.
					return 0
				}
			}*/
	case WM_LBUTTONDOWN:
		//		w.pointerButton(pointer.ButtonPrimary, true, lParam, getModifiers())
	case WM_LBUTTONUP:
		//		w.pointerButton(pointer.ButtonPrimary, false, lParam, getModifiers())
	case WM_RBUTTONDOWN:
		//		w.pointerButton(pointer.ButtonSecondary, true, lParam, getModifiers())
	case WM_RBUTTONUP:
		//		w.pointerButton(pointer.ButtonSecondary, false, lParam, getModifiers())
	case WM_MBUTTONDOWN:
		//		w.pointerButton(pointer.ButtonTertiary, true, lParam, getModifiers())
	case WM_MBUTTONUP:
		//		w.pointerButton(pointer.ButtonTertiary, false, lParam, getModifiers())
	case WM_CANCELMODE:
		//		w.w.Event(pointer.Event{
		//			Kind: pointer.Cancel,
		//		})
	case WM_SETFOCUS:
		w.focused = true
		//		w.w.Event(key.FocusEvent{Focus: true})
	case WM_KILLFOCUS:
		w.focused = false
		//		w.w.Event(key.FocusEvent{Focus: false})
	case WM_NCACTIVATE:
		if w.stage >= StageInactive {
			if wParam == TRUE {
				//				w.setStage(system.StageRunning)
			} else {
				//				w.setStage(system.StageInactive)
			}
		}
	case WM_NCHITTEST:
		if w.config.Decorated {
			// Let the system handle it.
			break
		}
		//		x, y := coordsFromlParam(lParam)
		x := 10.0
		y := 20.0
		np := Point{X: int32(x), Y: int32(y)}
		ScreenToClient(w.Hwnd, &np)
		//		return w.hitTest(int(np.X), int(np.Y))
	case WM_MOUSEMOVE:

		x, y := coordsFromlParam(lParam)
		fmt.Println(x, y)
		/*
			p := f32.Point{X: float32(x), Y: float32(y)}

				w.w.Event(pointer.Event{
					Kind:      pointer.Move,
					Source:    pointer.Mouse,
					Position:  p,
					Buttons:   w.pointerBtns,
					Time:      GetMessageTime(),
					Modifiers: getModifiers(),
				})
		*/
	case WM_MOUSEWHEEL:
		//		w.scrollEvent(wParam, lParam, false, getModifiers())
	case WM_MOUSEHWHEEL:
		//		w.scrollEvent(wParam, lParam, true, getModifiers())
	case WM_DESTROY:
		//		w.w.Event(ViewEvent{})
		//		w.w.Event(system.DestroyEvent{})
		if w.Hdc != 0 {
			ReleaseDC(w.Hdc)
			w.Hdc = 0
		}
		// The system destroys the HWND for us.
		w.Hwnd = 0
		PostQuitMessage(0)
	case WM_NCCALCSIZE:
		if w.config.Decorated {
			// Let Windows handle decorations.
			break
		}
		// No client areas; we draw decorations ourselves.
		if wParam != 1 {
			return 0
		}
		// lParam contains an NCCALCSIZE_PARAMS for us to adjust.
		place := GetWindowPlacement(w.Hwnd)
		if !place.IsMaximized() {
			// Nothing do adjust.
			return 0
		}
		// Adjust window position to avoid the extra padding in maximized
		// state. See https://devblogs.microsoft.com/oldnewthing/20150304-00/?p=44543.
		// Note that trying to do the adjustment in WM_GETMINMAXINFO is ignored by
		szp := (*NCCalcSizeParams)(unsafe.Pointer(uintptr(lParam)))
		mi := GetMonitorInfo(w.Hwnd)
		szp.Rgrc[0] = mi.WorkArea
		return 0
	case WM_PAINT:
		w.draw(true)
	case WM_SIZE:
		w.update()
		switch wParam {
		case SIZE_MINIMIZED:
			w.config.Mode = Minimized
			//			w.setStage(system.StagePaused)
		case SIZE_MAXIMIZED:
			w.config.Mode = Maximized
			//			w.setStage(system.StageRunning)
		case SIZE_RESTORED:
			if w.config.Mode != Fullscreen {
				w.config.Mode = Windowed
			}
			//			w.setStage(system.StageRunning)
		}
	case WM_GETMINMAXINFO:
		mm := (*MinMaxInfo)(unsafe.Pointer(uintptr(lParam)))
		var bw, bh int32
		if w.config.Decorated {
			r := GetWindowRect(w.Hwnd)
			cr := GetClientRect(w.Hwnd)
			bw = r.Right - r.Left - (cr.Right - cr.Left)
			bh = r.Bottom - r.Top - (cr.Bottom - cr.Top)
		}
		if p := w.config.MinSize; p.X > 0 || p.Y > 0 {
			mm.PtMinTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		if p := w.config.MaxSize; p.X > 0 || p.Y > 0 {
			mm.PtMaxTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		return 0
	case WM_SETCURSOR:
		/*
			w.cursorIn = (lParam & 0xffff) == HTCLIENT
			if w.cursorIn {
				SetCursor(w.cursor)
				return TRUE
			}
		*/
		/*
			case _WM_WAKEUP:
			w.w.Event(wakeupEvent{})

					case WM_IME_STARTCOMPOSITION:
						imc := ImmGetContext(w.hwnd)
						if imc == 0 {
							return TRUE
						}
						defer ImmReleaseContext(w.hwnd, imc)
						sel := w.w.EditorState().Selection
						caret := sel.Transform.Transform(sel.Caret.Pos.Add(f32.Pt(0, sel.Caret.Descent)))
						icaret := image.Pt(int(caret.X+.5), int(caret.Y+.5))
						ImmSetCompositionWindow(imc, icaret.X, icaret.Y)
						ImmSetCandidateWindow(imc, icaret.X, icaret.Y)
		*/
		/*
			case  WM_IME_COMPOSITION:
				imc :=  ImmGetContext(w.hwnd)
				if imc == 0 {
					return  TRUE
				}
				defer  ImmReleaseContext(w.hwnd, imc)
				state := w.w.EditorState()
				rng := state.compose
				if rng.Start == -1 {
					rng = state.Selection.Range
				}
				if rng.Start > rng.End {
					rng.Start, rng.End = rng.End, rng.Start
				}
				var replacement string
				switch {
				case lParam& GCS_RESULTSTR != 0:
					replacement =  ImmGetCompositionString(imc,  GCS_RESULTSTR)
				case lParam& GCS_COMPSTR != 0:
					replacement =  ImmGetCompositionString(imc,  GCS_COMPSTR)
				}
				end := rng.Start + utf8.RuneCountInString(replacement)
				w.w.EditorReplace(rng, replacement)
				state = w.w.EditorState()
				comp := Range{
					Start: rng.Start,
					End:   end,
				}
				if lParam& GCS_DELTASTART != 0 {
					start :=  ImmGetCompositionValue(imc,  GCS_DELTASTART)
					comp.Start = state.RunesIndex(state.UTF16Index(comp.Start) + start)
				}
				w.w.SetComposingRegion(comp)
				pos := end
				if lParam& GCS_CURSORPOS != 0 {
					rel :=  ImmGetCompositionValue(imc,  GCS_CURSORPOS)
					pos = state.RunesIndex(state.UTF16Index(rng.Start) + rel)
				}
				w.w.SetEditorSelection(key.Range{Start: pos, End: pos})
				return  TRUE
			case  WM_IME_ENDCOMPOSITION:
				w.w.SetComposingRegion(key.Range{Start: -1, End: -1})
				return  TRUE
		*/
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

func convertKeyCode(code uintptr) (string, bool) {
	if '0' <= code && code <= '9' || 'A' <= code && code <= 'Z' {
		return string(rune(code)), true
	}
	var r string

	switch code {
	case VK_ESCAPE:
		r = NameEscape
	case VK_LEFT:
		r = NameLeftArrow
	case VK_RIGHT:
		r = NameRightArrow
	case VK_RETURN:
		r = NameReturn
	case VK_UP:
		r = NameUpArrow
	case VK_DOWN:
		r = NameDownArrow
	case VK_HOME:
		r = NameHome
	case VK_END:
		r = NameEnd
	case VK_BACK:
		r = NameDeleteBackward
	case VK_DELETE:
		r = NameDeleteForward
	case VK_PRIOR:
		r = NamePageUp
	case VK_NEXT:
		r = NamePageDown
	case VK_F1:
		r = NameF1
	case VK_F2:
		r = NameF2
	case VK_F3:
		r = NameF3
	case VK_F4:
		r = NameF4
	case VK_F5:
		r = NameF5
	case VK_F6:
		r = NameF6
	case VK_F7:
		r = NameF7
	case VK_F8:
		r = NameF8
	case VK_F9:
		r = NameF9
	case VK_F10:
		r = NameF10
	case VK_F11:
		r = NameF11
	case VK_F12:
		r = NameF12
	case VK_TAB:
		r = NameTab
	case VK_SPACE:
		r = NameSpace
	case VK_OEM_1:
		r = ";"
	case VK_OEM_PLUS:
		r = "+"
	case VK_OEM_COMMA:
		r = ","
	case VK_OEM_MINUS:
		r = "-"
	case VK_OEM_PERIOD:
		r = "."
	case VK_OEM_2:
		r = "/"
	case VK_OEM_3:
		r = "`"
	case VK_OEM_4:
		r = "["
	case VK_OEM_5, VK_OEM_102:
		r = "\\"
	case VK_OEM_6:
		r = "]"
	case VK_OEM_7:
		r = "'"
	case VK_CONTROL:
		r = NameCtrl
	case VK_SHIFT:
		r = NameShift
	case VK_MENU:
		r = NameAlt
	case VK_LWIN, VK_RWIN:
		r = NameSuper
	default:
		return "", false
	}
	return r, true
}
func coordsFromlParam(lParam uintptr) (int, int) {
	x := int(int16(lParam & 0xffff))
	y := int(int16((lParam >> 16) & 0xffff))
	return x, y
}

// Modifiers
type Modifiers uint32

const (
	// ModCtrl is the ctrl modifier
	ModCtrl Modifiers = 1 << iota
	// ModCommand is the command modifier key
	// found on Apple keyboards.
	ModCommand
	// ModShift is the shift modifier
	ModShift
	// ModAlt is the alt modifier key, or the option
	// key on Apple keyboards.
	ModAlt
	// ModSuper is the "logo" modifier key, often
	// represented by a Windows logo.
	ModSuper
)

func getModifiers() Modifiers {
	var kmods Modifiers
	if GetKeyState(VK_LWIN)&0x1000 != 0 || GetKeyState(VK_RWIN)&0x1000 != 0 {
		kmods |= ModSuper
	}
	if GetKeyState(VK_MENU)&0x1000 != 0 {
		kmods |= ModAlt
	}
	if GetKeyState(VK_CONTROL)&0x1000 != 0 {
		kmods |= ModCtrl
	}
	if GetKeyState(VK_SHIFT)&0x1000 != 0 {
		kmods |= ModShift
	}
	return kmods
}

func (w *Window) draw(sync bool) {
	if w.config.Size.X == 0 || w.config.Size.Y == 0 {
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
	w.config.Size = image.Point{
		X: int(cr.Right - cr.Left),
		Y: int(cr.Bottom - cr.Top),
	}

	w.borderSize = image.Pt(
		GetSystemMetrics(SM_CXSIZEFRAME),
		GetSystemMetrics(SM_CYSIZEFRAME),
	)
	//		w.w.Event(ConfigEvent{Config: w.config})

}

const (
	// Names for special keys.
	NameLeftArrow      = "←"
	NameRightArrow     = "→"
	NameUpArrow        = "↑"
	NameDownArrow      = "↓"
	NameReturn         = "⏎"
	NameEnter          = "⌤"
	NameEscape         = "⎋"
	NameHome           = "⇱"
	NameEnd            = "⇲"
	NameDeleteBackward = "⌫"
	NameDeleteForward  = "⌦"
	NamePageUp         = "⇞"
	NamePageDown       = "⇟"
	NameTab            = "Tab"
	NameSpace          = "Space"
	NameCtrl           = "Ctrl"
	NameShift          = "Shift"
	NameAlt            = "Alt"
	NameSuper          = "Super"
	NameCommand        = "⌘"
	NameF1             = "F1"
	NameF2             = "F2"
	NameF3             = "F3"
	NameF4             = "F4"
	NameF5             = "F5"
	NameF6             = "F6"
	NameF7             = "F7"
	NameF8             = "F8"
	NameF9             = "F9"
	NameF10            = "F10"
	NameF11            = "F11"
	NameF12            = "F12"
	NameBack           = "Back"
)

package winapi

import (
	"image"
	"unicode"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

// Основной обработчик событий главного окна
func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	win, exists := winMap.Load(hwnd)
	if !exists {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}

	w := win.(*Window)

	switch msg {
	case WM_NCCREATE:
		panic("Main WM_NCCREATE")
	case WM_CREATE:
		panic("Main WM_CREATE")

	case WM_UNICHAR:
		if wParam == UNICODE_NOCHAR {
			// Tell the system that we accept WM_UNICHAR messages.
			return TRUE
		}
		fallthrough
	//	Оператор fallthrough используется в предложении case switch.
	// Он должен использоваться в конце предложения case.
	// Он используется для выполнения следующего предложения case без проверки выражения.
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

		if n, ok := convertKeyCode(wParam); ok {
			e := Event{
				Name:      n,
				Modifiers: getModifiers(),
				State:     Press,
			}
			if msg == WM_KEYUP || msg == WM_SYSKEYUP {
				e.State = Release
			}

			w.Config.EventChan <- (e)

			if (wParam == VK_F10) && (msg == WM_SYSKEYDOWN || msg == WM_SYSKEYUP) {
				// Reserve F10 for ourselves, and don't let it open the system menu. Other Windows programs
				// such as cmd.exe and graphical debuggers also reserve F10.
				return 0
			}
		}
	case WM_LBUTTONDOWN:
		w.pointerButton(ButtonPrimary, true, lParam, getModifiers())
	case WM_LBUTTONUP:
		w.pointerButton(ButtonPrimary, false, lParam, getModifiers())
	case WM_RBUTTONDOWN:
		w.pointerButton(ButtonSecondary, true, lParam, getModifiers())
	case WM_RBUTTONUP:
		w.pointerButton(ButtonSecondary, false, lParam, getModifiers())
	case WM_MBUTTONDOWN:
		w.pointerButton(ButtonTertiary, true, lParam, getModifiers())
	case WM_MBUTTONUP:
		w.pointerButton(ButtonTertiary, false, lParam, getModifiers())
	case WM_CANCELMODE:
		//		w.w.Event( Event{
		//			Kind:  Cancel,
		//		})
	case WM_SETFOCUS:
		// Это щелчок в окне
		w.Focused = true
		x, y := coordsFromlParam(lParam)
		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Enter,
			Source:    Mouse,
			Position:  image.Point{X: x, Y: y},
			Buttons:   w.PointerBtns,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}
	case WM_KILLFOCUS:
		// Щелчок вне нашего окна
		w.Focused = false
		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Leave,
			Source:    Mouse,
			Position:  image.Point{X: -1, Y: -1},
			Buttons:   w.PointerBtns,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}

	case WM_MOUSEMOVE:
		// Это событие будет, даже если наше окно не в фокусе
		// и может быть даже частично перекрыто другим окном
		x, y := coordsFromlParam(lParam)
		p := image.Point{X: x, Y: y}

		w.Config.EventChan <- Event{
			SWin:      w,
			Kind:      Move,
			Source:    Mouse,
			Position:  p,
			Buttons:   w.PointerBtns,
			Time:      GetMessageTime(),
			Modifiers: getModifiers(),
		}

	case WM_MOUSEWHEEL:
		//		w.scrollEvent(wParam, lParam, false, getModifiers())
	case WM_MOUSEHWHEEL:
		//		w.scrollEvent(wParam, lParam, true, getModifiers())
	case WM_NCACTIVATE:
		if w.Stage >= StageInactive {
			if wParam == TRUE {
				w.Stage = StageRunning
			} else {
				w.Stage = StageInactive
			}
		}

	case WM_NCHITTEST:
		//		if w.Config.Decorated {
		//			// Let the system handle it.
		//			break
		//		}
		x, y := coordsFromlParam(lParam)

		np := Point{X: int32(x), Y: int32(y)}
		ScreenToClient(w.Hwnd, &np)
		return w.hitTest(int(np.X), int(np.Y))

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
		//		if w.Config.Decorated {
		//			// Let Windows handle decorations.
		//			break
		//		}
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
			w.Config.Mode = Minimized
			w.Stage = StagePaused
		case SIZE_MAXIMIZED:
			w.Config.Mode = Maximized
			w.Stage = StageRunning
		case SIZE_RESTORED:
			if w.Config.Mode != Fullscreen {
				w.Config.Mode = Windowed
			}
			w.Stage = StageRunning
		}
	case WM_GETMINMAXINFO:
		mm := (*MinMaxInfo)(unsafe.Pointer(uintptr(lParam)))
		var bw, bh int32 = 0, 0
		//		if w.Config.Decorated {
		// Этот код дает косяки в отрисовке окна при перемещении
		//	r := GetWindowRect(w.Hwnd)
		//	cr := GetClientRect(w.Hwnd)
		//	bw = r.Right - r.Left - (cr.Right - cr.Left)
		//	bh = r.Bottom - r.Top - (cr.Bottom - cr.Top)
		//		}
		if p := w.Config.MinSize; p.X > 0 || p.Y > 0 {
			mm.PtMinTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		if p := w.Config.MaxSize; p.X > 0 || p.Y > 0 {
			mm.PtMaxTrackSize = Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		return 0
	case WM_SETCURSOR:

		w.CursorIn = (lParam & 0xffff) == HTCLIENT
		if w.CursorIn {
			SetCursor(w.Cursor)
			return TRUE
		}
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

// hitTest returns the non-client area hit by the point, needed to
// process WM_NCHITTEST.
func (w *Window) hitTest(x, y int) uintptr {
	if w.Config.Mode == Fullscreen {
		return HTCLIENT
	}
	if w.Config.Mode != Windowed {
		// Only windowed mode should allow resizing.
		return HTCLIENT
	}
	// Check for resize handle before system actions; otherwise it can be impossible to
	// resize a custom-decorations window when the system move area is flush with the
	// edge of the window.
	top := y <= w.Config.BorderSize.Y
	bottom := y >= w.Config.Size.Y-w.Config.BorderSize.Y
	left := x <= w.Config.BorderSize.X
	right := x >= w.Config.Size.X-w.Config.BorderSize.X
	switch {
	case top && left:
		return HTTOPLEFT
	case top && right:
		return HTTOPRIGHT
	case bottom && left:
		return HTBOTTOMLEFT
	case bottom && right:
		return HTBOTTOMRIGHT
	case top:
		return HTTOP
	case bottom:
		return HTBOTTOM
	case left:
		return HTLEFT
	case right:
		return HTRIGHT
	}
	/*
		p := f32.Pt(float32(x), float32(y))

		if a, ok := w.w.ActionAt(p); ok && a == system.ActionMove {
			return  HTCAPTION
		}
	*/
	return HTCLIENT
}

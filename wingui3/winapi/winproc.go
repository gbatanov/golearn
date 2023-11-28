package winapi

import (
	"fmt"
	"image"
	"unicode"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

// Основной обработчик событий окна
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
		w.Focused = true
		//		w.w.Event(FocusEvent{Focus: true})
	case WM_KILLFOCUS:
		w.Focused = false
		//		w.w.Event(FocusEvent{Focus: false})
	case WM_NCACTIVATE:
		if w.Stage >= StageInactive {
			if wParam == TRUE {
				w.Stage = StageRunning
			} else {
				w.Stage = StageInactive
			}
		}
		/*
			case WM_NCHITTEST:
				if w.Config.Decorated {
					// Let the system handle it.
					break
				}
				//		x, y := coordsFromlParam(lParam)
				x := 10.0
				y := 20.0
				np := Point{X: int32(x), Y: int32(y)}
				ScreenToClient(w.Hwnd, &np)
				//		return w.hitTest(int(np.X), int(np.Y))
		*/
	case WM_MOUSEMOVE:

		x, y := coordsFromlParam(lParam)

		fmt.Println(x, y)

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
		if w.Config.Decorated {
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
			w.Config.Mode = Minimized
			//			w.setStage(system.StagePaused)
		case SIZE_MAXIMIZED:
			w.Config.Mode = Maximized
			//			w.setStage(system.StageRunning)
		case SIZE_RESTORED:
			if w.Config.Mode != Fullscreen {
				w.Config.Mode = Windowed
			}
			//			w.setStage(system.StageRunning)
		}
	case WM_GETMINMAXINFO:
		mm := (*MinMaxInfo)(unsafe.Pointer(uintptr(lParam)))
		var bw, bh int32
		if w.Config.Decorated {
			r := GetWindowRect(w.Hwnd)
			cr := GetClientRect(w.Hwnd)
			bw = r.Right - r.Left - (cr.Right - cr.Left)
			bh = r.Bottom - r.Top - (cr.Bottom - cr.Top)
		}
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
				w.w.SetEditorSelection(Range{Start: pos, End: pos})
				return  TRUE
			case  WM_IME_ENDCOMPOSITION:
				w.w.SetComposingRegion(Range{Start: -1, End: -1})
				return  TRUE
		*/
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

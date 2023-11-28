package winapi

import (
	"image"
	"log"

	syscall "golang.org/x/sys/windows"
)

const ID_FIRSTCHILD = 100

var ChildId = 1

func CreateChildWindow(parent *Window, x, y, width, height int32) (*Window, error) {
	var resErr error
	resources.once.Do(func() {
		resErr = initResources(true)
	})
	if resErr != nil {
		return nil, resErr
	}
	const dwStyle = WS_THICKFRAME

	hwnd, err := CreateWindowEx(
		dwExStyle,
		resources.class,    //lpClassame
		"Child",            // lpWindowName
		WS_CHILD|WS_BORDER, //dwStyle
		x, y,               //x, y
		width, height, //w, h
		parent.Hwnd,      //hWndParent
		0,                // hMenu
		resources.handle, //hInstance
		0)                // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Hwnd: hwnd,
		Config: &Config{
			Title:     "Child",
			EventChan: parent.Config.EventChan,
			Size:      image.Pt(int(width), int(height)),
			MinSize:   parent.Config.MinSize,
			MaxSize:   parent.Config.MaxSize,
			Position:  image.Pt(int(x), int(y)),
		},
		Parent:    parent,
		Childrens: nil,
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}
	parent.Childrens = append(parent.Childrens, w)
	ShowWindow(hwnd, SW_SHOW)
	SetForegroundWindow(hwnd)
	SetFocus(hwnd)
	EnableWindow(hwnd, int32(1))
	ChildId += 1
	return w, nil
}

func windowChildProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	win, exists := winMap.Load(hwnd)
	if !exists {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	w := win.(*Window)
	if w.Config.Title == "GsbTest" {
		return windowProc(hwnd, msg, wParam, lParam)
	}

	switch msg {
	case WM_MOUSEMOVE:

		x, y := coordsFromlParam(lParam)

		log.Println(x, y)
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
	}

	return DefWindowProc(hwnd, msg, wParam, lParam)
}

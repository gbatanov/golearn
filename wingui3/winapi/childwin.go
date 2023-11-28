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
	resourceChild.once.Do(func() {
		resErr = initResources(true)
	})
	if resErr != nil {
		return nil, resErr
	}
	const dwStyle = /*WS_THICKFRAME |*/ WS_CHILD | WS_VISIBLE | WS_BORDER

	hwnd, err := CreateWindowEx(
		0,
		"Static", // resourceChild.class, //lpClassName
		"Child",  // lpWindowName
		dwStyle,  //dwStyle
		x, y,     //x, y
		width, height, //w, h
		parent.Hwnd,          //hWndParent
		0,                    // hMenu
		resourceChild.handle, //hInstance
		0)                    // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Id:   int32(ChildId),
		Hwnd: hwnd,
		Config: &Config{
			Decorated: false,
			Title:     "Childer",
			EventChan: parent.Config.EventChan,
			Size:      image.Pt(int(width), int(height)),
			MinSize:   parent.Config.MinSize,
			MaxSize:   parent.Config.MaxSize,
			Position:  image.Pt(int(x), int(y)),
			Mode:      Windowed,
		},
		Parent:    parent,
		Childrens: nil,
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}
	w.SetCursor(CursorDefault)
	parent.Childrens = append(parent.Childrens, w)
	//	ShowWindow(hwnd, SW_SHOWNORMAL)
	//	SetForegroundWindow(hwnd)
	//	SetFocus(hwnd)
	//	EnableWindow(hwnd, int32(1))

	ChildId += 1
	return w, nil
}

func windowChildProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	_, exists := winMap.Load(hwnd)
	if !exists {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	if msg == WM_CREATE {
		panic("WM_CREATE")
	}
	if msg == WM_NCCREATE {
		panic("WM_NCCREATE")
	}
	if msg == WM_CHILDACTIVATE {
		panic("WM_CHILDACTIVATE")
	}
	log.Printf("Child 0x%04x", msg)
	return windowProc(hwnd, msg, wParam, lParam)

}

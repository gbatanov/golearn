package winapi

const ID_FIRSTCHILD = 100

var ChildId = 1

func CreateChildWindow(parent *Window, config *Config) (*Window, error) {

	var dwStyle uint32 = WS_CHILD | WS_VISIBLE
	if config.BorderSize.X > 0 {
		dwStyle |= WS_BORDER
	}

	hwnd, err := CreateWindowEx(
		0,
		"Static",                                           // resourceChild.class, //lpClassName
		"Child",                                            // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		parent.Hwnd,  //hWndParent
		0,            // hMenu
		parent.HInst, //hInstance
		0)            // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Id:        int32(ChildId),
		Hwnd:      hwnd,
		HInst:     parent.HInst,
		Config:    config,
		Parent:    parent,
		Childrens: nil,
	}
	w.Hdc, err = GetDC(hwnd)
	if err != nil {
		return nil, err
	}
	w.SetCursor(CursorDefault)
	parent.Childrens = append(parent.Childrens, w)

	ChildId += 1
	return w, nil
}

/*
// Обработку событий в дочерних окнах перенаправляем в основное
func windowChildProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	win, exists := winMap.Load(hwnd)
	if !exists {
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	w := win.(*Window)
	if msg == WM_CREATE {
		panic("Child WM_CREATE")
	}
	if msg == WM_NCCREATE {
		panic("Child WM_NCCREATE")
	}
	if msg == WM_CHILDACTIVATE {
		panic("Child WM_CHILDACTIVATE")
	}
	//	log.Printf("Child 0x%04x", msg)
	switch msg {
	case WM_PAINT:
		log.Println("Draw child")
		w.childDraw(true)
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	return windowProc(hwnd, msg, wParam, lParam)

}

func (w *Window) childDraw(sync bool) {

	// Fill the region
	r2 := GetClientRect(w.Hwnd)

	FillRect(w.Hdc, &r2, GetStockObject(4)) // 0,5-белый, 1 - серый, 2-темно-серый, 4 - черный

}
*/

package winapi

var ChildId = 1

// Label
func CreateLabel(parent *Window, config Config) (*Window, error) {
	return CreateChildWindow(parent, config, "Static")
}

// Создаем статическое окно
func CreateChildWindow(parent *Window, config Config, class string) (*Window, error) {

	var dwStyle uint32 = WS_CHILD | WS_VISIBLE
	if config.BorderSize.X > 0 {
		dwStyle |= WS_BORDER
	}

	hwnd, err := CreateWindowEx(
		0,
		class,                                              // resourceChild.class, //lpClassName
		config.Title,                                       // lpWindowName
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
	parent.Childrens[ChildId] = w

	ChildId += 1
	return w, nil
}

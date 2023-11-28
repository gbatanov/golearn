package winapi

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
		Hwnd:      hwnd,
		Config:    nil,
		Parent:    parent,
		Childrens: nil,
	}
	parent.Childrens = append(parent.Childrens, w)
	ShowWindow(hwnd, SW_SHOW)
	ChildId += 1
	return w, nil
}

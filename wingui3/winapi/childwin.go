package winapi

import "golang.org/x/sys/windows"

const IDC_BUTTON_OK = 100
const IDC_BUTTON_CANCEL = 101

// Label
func CreateLabel(parent *Window, config Config, id int) (*Window, error) {
	return CreateChildWindow(parent, config, id)
}

// Button
func CreateButton(parent *Window, config Config, id int) (*Window, error) {
	return CreateChildWindow(parent, config, id)
}

// Создаем статическое окно
func CreateChildWindow(parent *Window, config Config, id int) (*Window, error) {

	var dwStyle uint32 = WS_CHILD | WS_VISIBLE
	if config.BorderSize.X > 0 {
		dwStyle |= WS_BORDER
	}
	var hMenu windows.Handle = 0
	if config.Class == "Button" {
		if config.Title == "Ok" {
			hMenu = windows.Handle(IDC_BUTTON_OK)
		} else if config.Title == "Cancel" {
			hMenu = windows.Handle(IDC_BUTTON_CANCEL)
		}
	}
	hwnd, err := CreateWindowEx(
		0,
		config.Class,                                       // standard static class,
		config.Title,                                       // lpWindowName
		dwStyle,                                            //dwStyle
		int32(config.Position.X), int32(config.Position.Y), //x, y
		int32(config.Size.X), int32(config.Size.Y), //w, h
		parent.Hwnd,  //hWndParent
		hMenu,        // hMenu
		parent.HInst, //hInstance
		0)            // lpParam
	if err != nil {
		return nil, err
	}
	w := &Window{
		Id:        int32(id),
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
	WinMap.Store(w.Hwnd, w)

	return w, nil
}

package window

type GetWindowInfo interface {
	Window() *WindowInfo
}

type WindowInfo struct {
	Title string
	Class string
}

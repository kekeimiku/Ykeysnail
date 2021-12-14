package window

type Darwin struct{}

func (w *Darwin) Window() *WindowInfo {
	return &WindowInfo{}
}

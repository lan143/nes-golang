package display

type Display interface {
	Init()
	Run()
	UpdateScreen()
	RenderPixel(x, y uint16, color uint32)
}

package display

type Display interface {
	Init()
	Run()
	UpdateScreen()
	RenderPixel(x, y int, color uint32)
}

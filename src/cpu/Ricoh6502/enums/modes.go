package enums

type Modes byte

const (
	ModeAcc  Modes = 1
	ModeZP         = 2
	ModeZPX        = 3
	ModeZPY        = 4
	ModeABS        = 5
	ModeABSX       = 6
	ModeABSY       = 7
	ModeIMM        = 8
	ModeINDX       = 9
	ModeINDY       = 10
	ModeIMP        = 11
	ModeIND        = 12
	ModeREL        = 13
)

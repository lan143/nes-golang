package enum

type VideoSystem uint8

const (
	VideoSystemPAL VideoSystem = iota
	VideoSystemNTSC
	VideoSystemDendy
)

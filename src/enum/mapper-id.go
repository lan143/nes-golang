package enum

type MapperId byte

const (
	MapperNROM    MapperId = 0
	MapperMMC1             = 1
	MapperUnROM            = 2
	INESMapper003          = 3
	MapperMMC3             = 4
	MapperMMC5             = 5
)

func (i MapperId) String() string {
	switch i {
	case MapperNROM:
		return "NROM"
	case MapperMMC1:
		return "MMC1"
	case MapperUnROM:
		return "UnROM"
	case MapperMMC3:
		return "MMC3"
	case MapperMMC5:
		return "MMC5"
	}

	return "UNKNOWN"
}

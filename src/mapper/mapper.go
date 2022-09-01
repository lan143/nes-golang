package mapper

type Mapper interface {
	LoadRom(data []byte)
	GetByte(address uint16) byte
	PutByte(address uint16, value byte)
}

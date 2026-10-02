package cartridge

import (
	"github.com/pkg/errors"
	"log"
	"main/src/enum"
	"main/src/mapper"
	"main/src/rom"
)

type Cartridge struct {
	rom    rom.Rom
	mapper mapper.Mapper

	mapperFactory *mapper.Factory

	chrRomOffset     uint32
	prgRomSizeBytes  uint32
	chrRomSizeBytes  uint32
	romMirroringType enum.MirroringType
	hasChrRom        bool
}

func (c *Cartridge) LoadRom(r rom.Rom) error {
	var err error
	c.rom = r

	log.Printf("Mapper: %s", enum.MapperId(c.rom.GetMapperId()))
	log.Printf("PRG ROM Size: %d", c.rom.GetPrgRomSize())
	log.Printf("CHR ROM Size: %d", c.rom.GetChrRomSize())

	c.prgRomSizeBytes = uint32(c.rom.GetPrgRomSize()) * 0x4000
	c.chrRomSizeBytes = uint32(c.rom.GetChrRomSize()) * 0x2000
	c.chrRomOffset = c.prgRomSizeBytes
	c.romMirroringType = c.rom.GetMirroringType()
	c.hasChrRom = c.rom.GetChrRomSize() > 0

	c.mapper, err = c.mapperFactory.GetMapper(enum.MapperId(c.rom.GetMapperId()))
	if err != nil {
		return errors.Wrap(err, "cartridge.load-rom.get-mapper")
	}

	err = c.mapper.Init(c.rom.GetPrgRomSize())
	if err != nil {
		return errors.Wrap(err, "cartridge.load-rom.init-mapper")
	}

	return nil
}

func (c *Cartridge) HasChrRom() bool {
	return c.hasChrRom
}

func (c *Cartridge) GetMirroringType() enum.MirroringType {
	if c.mapper.GetMirroringType() != 0 {
		return c.mapper.GetMirroringType()
	} else {
		return c.romMirroringType
	}
}

func (c *Cartridge) GetByte(address uint16) byte {
	if address < 0x2000 {
		mapped := c.mapper.MapChrRom(address)
		// Mapper bank registers may exceed the actual ROM in corner cases;
		// real hardware aliases, so keep the index in bounds.
		if c.chrRomSizeBytes > 0 {
			mapped %= c.chrRomSizeBytes
		}

		return c.rom.GetByte(c.chrRomOffset + mapped)
	}

	mapped := c.mapper.MapPrgRom(address)
	if c.prgRomSizeBytes > 0 {
		mapped %= c.prgRomSizeBytes
	}

	return c.rom.GetByte(mapped)
}

func (c *Cartridge) PutByte(address uint16, value byte) {
	c.mapper.PutByte(address, value)
}

func NewCartridge(mapperFactory *mapper.Factory) *Cartridge {
	return &Cartridge{
		mapperFactory: mapperFactory,
	}
}

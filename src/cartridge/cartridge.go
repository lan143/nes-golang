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
}

func (c *Cartridge) HasChrRom() bool {
	return c.rom.GetChrRomSize() > 0
}

func (c *Cartridge) GetMirroringType() enum.MirroringType {
	if c.mapper.GetMirroringType() != 0 {
		return c.mapper.GetMirroringType()
	} else {
		return c.rom.GetMirroringType()
	}
}

func (c *Cartridge) LoadRom(r rom.Rom) error {
	var err error
	c.rom = r

	log.Printf("Mapper: %s", enum.MapperId(c.rom.GetMapperId()))
	log.Printf("PRG ROM Size: %d", c.rom.GetPrgRomSize())
	log.Printf("CHR ROM Size: %d", c.rom.GetChrRomSize())

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

func (c *Cartridge) GetByte(address uint16) byte {
	if address < 0x2000 {
		return c.rom.GetByte(uint32(c.rom.GetPrgRomSize())*0x4000 + c.mapper.MapChrRom(address))
	} else {
		return c.rom.GetByte(c.mapper.MapPrgRom(address))
	}
}

func (c *Cartridge) PutByte(address uint16, value byte) {
	c.mapper.PutByte(address, value)
}

func NewCartridge(mapperFactory *mapper.Factory) *Cartridge {
	return &Cartridge{
		mapperFactory: mapperFactory,
	}
}

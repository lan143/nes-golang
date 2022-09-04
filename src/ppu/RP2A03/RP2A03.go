package RP2A03

import (
	"main/src/bus"
	"main/src/display"
	"main/src/mapper"
	"main/src/mapper/enum"
	"main/src/ppu/RP2A03/registers"
	"main/src/ppu/RP2A03/sprites"
)

type PPU struct {
	mapper  mapper.Mapper
	display display.Display
	bus     *bus.Bus
	palette Palette

	frame    uint64
	scanline uint16
	cycle    uint16

	vRam    [16384]byte
	oamRam  [256]byte
	oamRam2 [32]byte

	ppuCtrl   registers.PPUCtrlRegister   // 0x2000
	ppuMask   registers.PPUMaskRegister   // 0x2001
	ppuStatus registers.PPUStatusRegister // 0x2002
	oamAddr   byte                        // 0x2003
	oamData   byte                        // 0x2004
	ppuScroll byte                        // 0x2005
	ppuAddr   byte                        // 0x2006
	ppuData   byte                        // 0x2007
	oamDma    byte                        // 0x4014

	nameTableRegister          byte
	attributeTableLowRegister  registers.Uint16Register
	attributeTableHighRegister registers.Uint16Register
	patternTableLowRegister    registers.Uint16Register
	patternTableHighRegister   registers.Uint16Register

	nameTableLatch          byte
	attributeTableLowLatch  byte
	attributeTableHighLatch byte
	patternTableLowLatch    byte
	patternTableHighLatch   byte

	fineXScroll         uint16
	currentVRamAddress  uint16
	temporalVRamAddress uint16

	vRamReadBuffer     byte
	registerFirstStore bool

	spritePixels     [256]uint32
	spriteIds        [256]uint32
	spritePriorities [256]uint32

	spritesManager  sprites.Manager
	spritesManager2 sprites.Manager
}

func (p *PPU) Init(mapper mapper.Mapper, display display.Display) {
	p.mapper = mapper
	p.display = display
	p.palette.Init()

	p.nameTableLatch = 0
	p.attributeTableLowLatch = 0
	p.attributeTableHighLatch = 0
	p.patternTableLowLatch = 0
	p.patternTableHighLatch = 0
	p.fineXScroll = 0
	p.currentVRamAddress = 0
	p.temporalVRamAddress = 0
	p.vRamReadBuffer = 0
	p.registerFirstStore = true

	var i uint16
	for i = 0; i < 256; i++ {
		p.spritePixels[i] = 0x80000000
		p.spriteIds[i] = 0x80000000
		p.spritePriorities[i] = 0x80000000
	}

	p.spritesManager.Init(p.oamRam[:])
	p.spritesManager2.Init(p.oamRam2[:])

	p.bus.Subscribe(bus.Write2000, func() {
		p.ppuCtrl.SetValue(p.mapper.GetByte(0x2000))

		p.temporalVRamAddress &= ^uint16(0xC00)
		p.temporalVRamAddress |= (uint16(p.ppuCtrl.GetValue()) & 0x3) << 10
	})
	p.bus.Subscribe(bus.Write2001, func() {
		p.ppuMask.SetValue(p.mapper.GetByte(0x2001))
	})
	p.bus.Subscribe(bus.Write2003, func() {
		p.oamAddr = p.mapper.GetByte(0x2003)
	})
	p.bus.Subscribe(bus.Write2004, func() {
		p.oamData = p.mapper.GetByte(0x2004)
		p.oamRam[p.oamAddr] = p.mapper.GetByte(0x2004)
		p.oamAddr++
	})
	p.bus.Subscribe(bus.Write2005, func() {
		p.ppuScroll = p.mapper.GetByte(0x2005)

		if p.registerFirstStore {
			p.fineXScroll = uint16(p.ppuScroll) & 0x7
			p.temporalVRamAddress &= ^uint16(0x1F)
			p.temporalVRamAddress |= (uint16(p.mapper.GetByte(0x2005)) >> 3) & 0x1F
		} else {
			p.temporalVRamAddress &= ^uint16(0x73E0)
			p.temporalVRamAddress |= (uint16(p.mapper.GetByte(0x2005)) & 0xF8) << 2
			p.temporalVRamAddress |= (uint16(p.mapper.GetByte(0x2005)) & 0x7) << 12
		}

		p.registerFirstStore = !p.registerFirstStore
	})
	p.bus.Subscribe(bus.Write2006, func() {
		value := p.mapper.GetByte(0x2006)

		if p.registerFirstStore {
			p.temporalVRamAddress &= ^uint16(0x7F00)
			p.temporalVRamAddress |= (uint16(value) & 0x3F) << 8
			//p.temporalVRamAddress = uint16(value) << 8
		} else {
			p.ppuAddr = value
			p.temporalVRamAddress &= ^uint16(0xFF)
			p.temporalVRamAddress |= uint16(value) & 0xFF
			//p.temporalVRamAddress |= uint16(value)
			p.currentVRamAddress = p.temporalVRamAddress
		}

		p.registerFirstStore = !p.registerFirstStore
	})
	p.bus.Subscribe(bus.Write2007, func() {
		p.ppuData = p.mapper.GetByte(0x2007)
		p.store(p.currentVRamAddress, p.ppuData)

		var incAddr uint16
		if p.ppuCtrl.IsIncrementAddress() {
			incAddr = 32
		} else {
			incAddr = 1
		}

		p.currentVRamAddress += incAddr
		p.currentVRamAddress &= 0x7FFF
		p.ppuAddr = byte(p.currentVRamAddress) & 0xFF
	})
	p.bus.Subscribe(bus.Write4014, func() {
		p.oamDma = p.mapper.GetByte(0x4014)
		offset := uint16(p.oamDma) * 0x100
		var i uint16

		for i = uint16(p.oamAddr); i < 256; i++ {
			p.oamRam[i] = p.mapper.GetByte(offset + i)
		}
	})
}

func (p *PPU) Run() {
	p.ppuStatus.SetVBlank()

	for {
		p.runCycle()
	}
}

func (p *PPU) runCycle() {
	p.renderPixel()
	p.shiftRegisters()
	p.fetch()
	p.evaluateSprites()
	p.updateFlags()
	p.countUpScrollCounters()
	p.countUpCycle()
}

func (p *PPU) countUpCycle() {
	p.cycle++

	if p.cycle > 340 {
		p.cycle = 0
		p.scanline++

		if p.scanline > 261 {
			p.scanline = 0
			p.frame++
		}
	}
}

func (p *PPU) countUpScrollCounters() {
	if !p.ppuMask.IsBackgroundVisible() && !p.ppuMask.IsSpritesVisible() {
		return
	}

	if p.scanline >= 240 && p.scanline <= 260 {
		return
	}

	if p.scanline == 261 {
		if p.cycle >= 280 && p.cycle <= 304 {
			p.currentVRamAddress &= ^uint16(0x7BE0)
			p.currentVRamAddress |= p.temporalVRamAddress & 0x7BE0
		}
	}

	if p.cycle == 0 || (p.cycle >= 258 && p.cycle <= 320) {
		return
	}

	if (p.cycle % 8) == 0 {
		var v = p.currentVRamAddress

		if (v & 0x1F) == 31 {
			v &= ^uint16(0x1F)
			v ^= 0x400
		} else {
			v++
		}

		p.currentVRamAddress = v
	}

	if p.cycle == 256 {
		var v = p.currentVRamAddress

		if (v & 0x7000) != 0x7000 {
			v += 0x1000
		} else {
			v &= ^uint16(0x7000)
			var y = (v & 0x3E0) >> 5

			if y == 29 {
				y = 0
				v ^= 0x800
			} else if y == 31 {
				y = 0
			} else {
				y++
			}

			v = (v & ^uint16(0x3E0)) | (y << 5)
		}

		p.currentVRamAddress = v
	}

	if p.cycle == 257 {
		p.currentVRamAddress &= ^uint16(0x41F)
		p.currentVRamAddress |= p.temporalVRamAddress & 0x41F
	}
}

func (p *PPU) updateFlags() {
	if p.cycle == 1 {
		if p.scanline == 241 {
			p.ppuStatus.SetVBlank()
			p.display.UpdateScreen()
		} else if p.scanline == 261 {
			p.ppuStatus.ClearVBlank()
			p.ppuStatus.ClearSpriteZeroHit()
			p.ppuStatus.ClearSpriteOverflow()
		}
	}

	if p.cycle == 10 {
		if p.scanline == 241 {
			if p.ppuCtrl.IsNMIVBlank() {
				p.bus.PushEvent(bus.NMIInterrupt)
			}
		}
	}
}

func (p *PPU) evaluateSprites() {
	if p.scanline >= 240 {
		return
	}

	if p.cycle == 0 {
		p.processSpritePixels()

		var i int
		il := len(p.oamRam2)

		for i = 0; i < il; i++ {
			p.oamRam2[i] = 0xFF
		}
	} else if p.cycle == 65 {
		var height byte
		var n = 0

		if p.ppuCtrl.IsSpritesSize() {
			height = 16
		} else {
			height = 8
		}

		var i int
		il := p.spritesManager.GetCount()

		for i = 0; i < il; i++ {
			s := p.spritesManager.GetSprite(i)

			if s.On(p.scanline, height) {
				if n < 8 {
					p.spritesManager2.Copy(n, s)
					n++
				} else {
					p.ppuStatus.SetSpriteOverflow()
					break
				}
			}
		}
	}
}

func (p *PPU) processSpritePixels() {
	ay := p.scanline - 1
	var i int
	il := len(p.spritePixels)

	for i = 0; i < il; i++ {
		p.spritePixels[i] = 0x80000000
		p.spriteIds[i] = 0x80000000
		p.spritePriorities[i] = 0x80000000
	}

	var height uint16
	if p.ppuCtrl.IsSpritesSize() {
		height = 16
	} else {
		height = 8
	}

	il = p.spritesManager2.GetCount()

	for i = 0; i < il; i++ {
		var s = p.spritesManager2.GetSprite(i)

		if s.IsEmpty() {
			break
		}

		bx := uint16(s.GetXPosition())
		by := uint16(s.GetYPosition())
		j := ay - by

		var cy uint16
		if s.DoFlipVertically() {
			cy = height - j - 1
		} else {
			cy = j
		}

		horizontal := s.DoFlipHorizontally()
		var ptIndex uint16

		if height == 8 {
			ptIndex = uint16(s.GetTileIndex())
		} else {
			ptIndex = s.GetTileIndexForSize16()
		}

		msb := uint16(s.GetPalletNum())
		var k uint16

		for k = 0; k < 8; k++ {
			var cx uint16
			if horizontal {
				cx = 7 - k
			} else {
				cx = k
			}

			x := bx + k

			if x >= 256 {
				break
			}

			lsb := uint16(p.getPatternTableElement(ptIndex, cx, cy, height))

			if lsb != 0 {
				pIndex := (msb << 2) | lsb

				if p.spritePixels[x] == 0x80000000 {
					p.spritePixels[x] = p.palette.GetValue(p.load(0x3F10 + pIndex))
					p.spriteIds[x] = uint32(s.GetId())
					p.spritePriorities[x] = uint32(s.GetPriority())
				}
			}
		}
	}
}

func (p *PPU) getPatternTableElement(index, x, y, ySize uint16) byte {
	ax := x % 8
	var a, b byte

	if ySize == 8 {
		ay := y % 8
		var offset uint16

		if p.ppuCtrl.IsSpritesPatternTable() {
			offset = 0x1000
		} else {
			offset = 0
		}

		a = p.load(offset + index*0x10 + ay)
		b = p.load(offset + index*0x10 + 0x8 + ay)
	} else {
		ay := y % 8
		ay += (y >> 3) * 0x10
		a = p.load(index + ay)
		b = p.load(index + ay + 0x8)
	}

	return ((a >> (7 - byte(ax))) & 1) | (((b >> (7 - byte(ax))) & 1) << 1)
}

func (p *PPU) fetch() {
	if p.scanline >= 240 && p.scanline <= 260 {
		return
	}

	if p.cycle == 0 {
		return
	}

	if (p.cycle >= 257 && p.cycle <= 320) || p.cycle >= 337 {
		return
	}

	switch (p.cycle - 1) % 8 {
	case 0:
		p.fetchNameTable()
		break
	case 2:
		p.fetchAttributeTable()
		break
	case 4:
		p.fetchPatternTableLow()
		break
	case 6:
		p.fetchPatternTableHigh()
		break
	default:
		break
	}

	if p.cycle%8 == 1 {
		p.nameTableRegister = p.nameTableLatch
		p.attributeTableLowRegister.StoreLowerByte(p.attributeTableLowLatch)
		p.attributeTableHighRegister.StoreLowerByte(p.attributeTableHighLatch)
		p.patternTableLowRegister.StoreLowerByte(p.patternTableLowLatch)
		p.patternTableHighRegister.StoreLowerByte(p.patternTableHighLatch)
	}
}

func (p *PPU) fetchNameTable() {
	p.nameTableLatch = p.load(0x2000 | (p.currentVRamAddress & 0x0FFF))
}

func (p *PPU) fetchAttributeTable() {
	v := p.currentVRamAddress
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)

	b := p.load(address)

	coarseX := v & 0x1F
	coarseY := (v >> 5) & 0x1

	var topBottom, rightLeft byte
	if (coarseY % 4) >= 2 { // bottom, top
		topBottom = 1
	} else {
		topBottom = 0
	}

	if (coarseX % 4) >= 2 { // right, left
		rightLeft = 1
	} else {
		rightLeft = 0
	}

	position := (topBottom << 1) | rightLeft

	value := (b >> (position << 1)) & 0x3
	highBit := value >> 1
	lowBit := value & 1

	if highBit == 1 {
		p.attributeTableHighLatch = 0xff
	} else {
		p.attributeTableHighLatch = 0
	}

	if lowBit == 1 {
		p.attributeTableLowLatch = 0xff
	} else {
		p.attributeTableLowLatch = 0
	}
}

func (p *PPU) fetchPatternTableLow() {
	fineY := (p.currentVRamAddress >> 12) & 0x7
	index := uint16(p.ppuCtrl.GetBackgroundPatternTable())*0x1000 +
		uint16(p.nameTableRegister*0x10) + fineY

	p.patternTableLowLatch = p.load(index)
}

func (p *PPU) fetchPatternTableHigh() {
	fineY := (p.currentVRamAddress >> 12) & 0x7
	var index = uint16(p.ppuCtrl.GetBackgroundPatternTable())*0x1000 +
		uint16(p.nameTableRegister)*0x10 + fineY

	p.patternTableHighLatch = p.load(index + 0x8)
}

func (p *PPU) shiftRegisters() {
	if p.scanline >= 240 && p.scanline <= 260 {
		return
	}

	if (p.cycle >= 1 && p.cycle <= 256) ||
		(p.cycle >= 329 && p.cycle <= 336) {
		p.patternTableLowRegister.Shift(0)
		p.patternTableHighRegister.Shift(0)
		p.attributeTableLowRegister.Shift(0)
		p.attributeTableHighRegister.Shift(0)
	}
}

func (p *PPU) getNameTableAddressWithMirroring(address uint16) uint16 {
	address = address & 0x2FFF // just in case

	var baseAddress uint16

	switch p.mapper.GetMirroringType() {
	case enum.SingleScreen:
		baseAddress = 0x2000
		break
	case enum.Horizontal:
		if address >= 0x2000 && address < 0x2400 {
			baseAddress = 0x2000
		} else if address >= 0x2400 && address < 0x2800 {
			baseAddress = 0x2000
		} else if address >= 0x2800 && address < 0x2C00 {
			baseAddress = 0x2400
		} else {
			baseAddress = 0x2400
		}
		break
	case enum.Vertical:
		if address >= 0x2000 && address < 0x2400 {
			baseAddress = 0x2000
		} else if address >= 0x2400 && address < 0x2800 {
			baseAddress = 0x2400
		} else if address >= 0x2800 && address < 0x2C00 {
			baseAddress = 0x2000
		} else {
			baseAddress = 0x2400
		}
		break
	case enum.FourScreen:
		if address >= 0x2000 && address < 0x2400 {
			baseAddress = 0x2000
		} else if address >= 0x2400 && address < 0x2800 {
			baseAddress = 0x2400
		} else if address >= 0x2800 && address < 0x2C00 {
			baseAddress = 0x2800
		} else {
			baseAddress = 0x2C00
		}
		break
	}

	return baseAddress | (address & 0x3FF)
}

func (p *PPU) store(address uint16, value byte) {
	address = address & 0x3FFF // just in case

	// 0x0000 - 0x1FFF is mapped with cartridge's CHR-ROM if it exists

	if address < 0x2000 && p.mapper.HasChrRom() {
		p.mapper.PutByte(address, value)
		return
	}

	// 0x0000 - 0x0FFF: pattern table 0
	// 0x1000 - 0x1FFF: pattern table 1
	// 0x2000 - 0x23FF: nametable 0
	// 0x2400 - 0x27FF: nametable 1
	// 0x2800 - 0x2BFF: nametable 2
	// 0x2C00 - 0x2FFF: nametable 3
	// 0x3000 - 0x3EFF: Mirrors of 0x2000 - 0x2EFF
	// 0x3F00 - 0x3F1F: Palette RAM indices
	// 0x3F20 - 0x3FFF: Mirrors of 0x3F00 - 0x3F1F

	if address >= 0x2000 && address < 0x3F00 {
		p.vRam[p.getNameTableAddressWithMirroring(address&0x2FFF)] = value
		return
	}

	if address >= 0x3F00 && address < 0x4000 {
		address = address & 0x3F1F
	}

	// Addresses for palette
	// 0x3F10/0x3F14/0x3F18/0x3F1C are mirrors of
	// 0x3F00/0x3F04/0x3F08/0x3F0C.

	if address == 0x3F10 {
		address = 0x3F00
	}

	if address == 0x3F14 {
		address = 0x3F04
	}

	if address == 0x3F18 {
		address = 0x3F08
	}

	if address == 0x3F1C {
		address = 0x3F0C
	}

	p.vRam[address] = value
	p.mapper.PutByte(address, value)
}

func (p *PPU) load(address uint16) byte {
	address = address & 0x3FFF // just in case

	// 0x0000 - 0x1FFF is mapped with cartridge's CHR-ROM if it exists
	if address < 0x2000 && p.mapper.HasChrRom() {
		return p.mapper.GetChrByte(address)
	}

	// 0x0000 - 0x0FFF: pattern table 0
	// 0x1000 - 0x1FFF: pattern table 1
	// 0x2000 - 0x23FF: nametable 0
	// 0x2400 - 0x27FF: nametable 1
	// 0x2800 - 0x2BFF: nametable 2
	// 0x2C00 - 0x2FFF: nametable 3
	// 0x3000 - 0x3EFF: Mirrors of 0x2000 - 0x2EFF
	// 0x3F00 - 0x3F1F: Palette RAM indices
	// 0x3F20 - 0x3FFF: Mirrors of 0x3F00 - 0x3F1F
	if address >= 0x2000 && address < 0x3F00 {
		return p.vRam[p.getNameTableAddressWithMirroring(address&0x2FFF)]
	}

	if address >= 0x3F00 && address < 0x4000 {
		address = address & 0x3F1F
	}

	// Addresses for palette
	// 0x3F10/0x3F14/0x3F18/0x3F1C are mirrors of
	// 0x3F00/0x3F04/0x3F08/0x3F0C.
	if address == 0x3F10 {
		address = 0x3F00
	}

	if address == 0x3F14 {
		address = 0x3F04
	}

	if address == 0x3F18 {
		address = 0x3F08
	}

	if address == 0x3F1C {
		address = 0x3F0C
	}

	return p.vRam[address]
}

func (p *PPU) getBackgroundPixel() uint32 {
	offset := 15 - p.fineXScroll

	lsb := (p.patternTableHighRegister.GetValue() & offset << 1) |
		p.patternTableLowRegister.GetValue()&offset
	msb := (p.attributeTableHighRegister.GetValue() & offset << 1) |
		p.attributeTableLowRegister.GetValue()&offset
	index := (msb << 2) | lsb

	if p.ppuMask.IsGreyscale() {
		index = index & 0x30
	}

	return p.palette.GetValue(p.load(0x3F00 + index))
}

func (p *PPU) renderPixel() {
	if p.cycle >= 257 || p.scanline >= 240 || p.cycle == 0 {
		return
	}

	x := p.cycle - 1
	y := p.scanline

	backgroundVisible := p.ppuMask.IsBackgroundVisible()
	spritesVisible := p.ppuMask.IsSpritesVisible()

	backgroundPixel := p.getBackgroundPixel()
	spritePixel := p.spritePixels[x]
	spriteId := p.spriteIds[x]
	spritePriority := p.spritePriorities[x]

	c := p.palette.GetValue(p.load(0x3F00))

	if backgroundVisible && spritesVisible {
		if spritePixel == 0x80000000 {
			c = backgroundPixel
		} else {
			if backgroundPixel == c {
				c = spritePixel
			} else if spritePriority == 0 {
				c = spritePixel
			} else {
				c = backgroundPixel
			}
		}
	} else if backgroundVisible && !spritesVisible {
		c = backgroundPixel
	} else if !backgroundVisible && spritesVisible {
		if spritePixel != 0x80000000 {
			c = spritePixel
		}
	}

	if p.ppuMask.IsEmphasizeRed() {
		c = c | 0x00FF0000
	}

	if p.ppuMask.IsEmphasizeGreen() {
		c = c | 0x0000FF00
	}

	if p.ppuMask.IsEmphasizeBlue() {
		c = c | 0x000000FF
	}

	if backgroundVisible && spritesVisible &&
		spriteId == 0 && spritePixel != 0 && backgroundPixel != 0 {
		p.ppuStatus.SetSpriteZeroHit()
	}

	p.display.RenderPixel(x, y, c)
}

func NewPPU(b *bus.Bus) *PPU {
	return &PPU{
		bus: b,
	}
}

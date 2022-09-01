package main

import (
	"main/src/cpu/Ricoh6502"
	"main/src/mapper"
	"main/src/ppu"
	"main/src/rom"
	"os"
)

func main() {
	f, err := os.Open("battle_city.nes")
	if err != nil {
		panic(err)
	}

	nesRom := rom.INes{}
	err = nesRom.Load(f)
	if err != nil {
		panic(err)
	}

	simpleMapper := mapper.SimpleMapper{}
	simpleMapper.LoadRom(nesRom.Data)

	ppuUnit := ppu.PPU{}
	ppuUnit.Init(&simpleMapper)
	go ppuUnit.Run()

	mos6502 := Ricoh6502.Cpu{}
	mos6502.Init(&simpleMapper)
	mos6502.Reset()
	mos6502.Run()
}

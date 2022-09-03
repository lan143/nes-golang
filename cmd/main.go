package main

import (
	"main/src/bus"
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

	simpleMapper := mapper.NROMMapper{}
	simpleMapper.LoadRom(nesRom.Data)

	b := &bus.Bus{}
	b.Init()

	ppuUnit := ppu.P2C02{}
	ppuUnit.Init(&simpleMapper, b)
	go ppuUnit.Run()

	mos6502 := Ricoh6502.Cpu{}

	mos6502.Init(&simpleMapper, b)
	mos6502.Reset()
	mos6502.Run()
}

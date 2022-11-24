# NES EMU

This is multiplatform NES emulator written in Golang. Created just for fun.

# How to launch
1. Download build for your OS from release tab.
2. Execute command ./nes-emu-<platform> <rom_path>
3. ???
4. PROFIT!

# How to build
1. Install go 1.18 or above
2. Download source code
3. Change dir to download sources dir.
4. Run ```go mod download```
5. Run ```go build -o nes-emu ./cmd/main.go```
6. Go to launch section

# Config
On the first run, a file **config.yaml** will be created.
The file contains the emulator settings.
At the moment, you can change the settings of the buttons of the first gamepad mapping to keyboard.

# Supported mappers
The following mappers are currently supported:
 * NROM
 * UnROM
 * CNROM (iNes Mapper 003)
 * MMC1
 * MMC3 (with bugs)

# Libraries used
 * faiface/pixel - for video output
 * hajimehoshi/oto - for audio output
 * uber/dig - for DI
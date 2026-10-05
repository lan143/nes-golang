package pixelgl

import (
	"sync"
	"testing"

	"main/src/config"
)

// writePattern fills the whole back buffer with a single color through the
// public RenderPixel path (which applies the y-flip).
func writePattern(d *Display, pattern uint32) {
	for y := uint16(0); y < NesImageSizeY; y++ {
		for x := uint16(0); x < NesImageSizeX; x++ {
			d.RenderPixel(x, y, pattern)
		}
	}
}

// assertResizeBufferUniform verifies every pixel in the resized RGBA buffer
// expands little-endian from the expected 32-bit color.
func assertResizeBufferUniform(t *testing.T, d *Display, expected uint32) {
	t.Helper()

	if got, want := len(d.resizeBuffer), NesImageSizeX*NesImageSizeY*4; got != want {
		t.Fatalf("resizeBuffer length = %d, want %d", got, want)
	}

	for i := 0; i < len(d.resizeBuffer); i += 4 {
		if d.resizeBuffer[i] != uint8(expected) ||
			d.resizeBuffer[i+1] != uint8(expected>>8) ||
			d.resizeBuffer[i+2] != uint8(expected>>16) ||
			d.resizeBuffer[i+3] != uint8(expected>>24) {
			t.Fatalf("pixel at byte offset %d = [%02X %02X %02X %02X], want [%02X %02X %02X %02X]",
				i,
				d.resizeBuffer[i], d.resizeBuffer[i+1], d.resizeBuffer[i+2], d.resizeBuffer[i+3],
				uint8(expected), uint8(expected>>8), uint8(expected>>16), uint8(expected>>24))
		}
	}
}

func TestRenderFrameCopiesCompleteFrame(t *testing.T) {
	d := NewDisplay(config.NewConfig(), nil)
	d.Init()

	const pattern uint32 = 0xAABBCCDD

	writePattern(d, pattern)
	d.RenderFrame()
	d.resizeFrame()

	assertResizeBufferUniform(t, d, pattern)
}

func TestFrontBufferNotUpdatedBeforeRenderFrame(t *testing.T) {
	d := NewDisplay(config.NewConfig(), nil)
	d.Init()

	const patternA uint32 = 0x11223344
	const patternB uint32 = 0x55667788

	writePattern(d, patternA)
	d.RenderFrame()
	d.resizeFrame()
	assertResizeBufferUniform(t, d, patternA)

	// Write a new frame but never publish it: the front buffer (and therefore
	// the resized output) must still show the previously published frame.
	writePattern(d, patternB)
	d.resizeFrame()
	assertResizeBufferUniform(t, d, patternA)
}

func TestNoRaceBetweenRenderAndCopy(t *testing.T) {
	d := NewDisplay(config.NewConfig(), nil)
	d.Init()

	const pattern uint32 = 0xAABBCCDD
	const iterations = 200

	stop := make(chan struct{})
	var readerWg sync.WaitGroup
	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				d.resizeFrame()
			}
		}
	}()

	for i := 0; i < iterations; i++ {
		writePattern(d, pattern)
		d.RenderFrame()
	}

	close(stop)
	readerWg.Wait()

	d.resizeFrame()
	assertResizeBufferUniform(t, d, pattern)
}

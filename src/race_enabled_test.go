//go:build race

package src

// raceEnabled is true when the binary is built with the race detector. The
// real-time pacing test measures wall-clock throughput, which the race
// instrumentation slows well below real hardware speed, so it is skipped.
const raceEnabled = true

package clock

import "time"

// Sec returns the number of seconds for 24 pulses per quarter note
func Sec(bpm float64) float64 {
	return 1 / (bpm / 60 * 24)
}

func Dur(bpm float64) time.Duration {
	return time.Duration(Sec(bpm) * float64(time.Second))
}

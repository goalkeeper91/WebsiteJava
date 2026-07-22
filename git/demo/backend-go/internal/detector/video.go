package detector

// NewMotionSpikeDetector returns a SpikeDetector tuned for frame-diff motion
// scores (normalized to [0, 1]).
func NewMotionSpikeDetector() *SpikeDetector {
	return NewSpikeDetector(0.1, 2.0, 0.01)
}

// FrameDiff computes the mean absolute difference between two equally-sized
// raw pixel buffers (e.g. rgb24), normalized to [0, 1]. Higher = more motion
// or scene change between the two frames.
func FrameDiff(prev, curr []byte) float64 {
	if len(prev) == 0 || len(prev) != len(curr) {
		return 0
	}

	var sum int64
	for i := range curr {
		d := int(curr[i]) - int(prev[i])
		if d < 0 {
			d = -d
		}
		sum += int64(d)
	}

	return float64(sum) / float64(len(curr)) / 255.0
}

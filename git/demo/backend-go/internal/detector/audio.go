package detector

import (
	"encoding/binary"
	"math"
)

// NewAudioSpikeDetector returns a SpikeDetector tuned for RMS audio energy
// readings (normalized to [0, 1]).
func NewAudioSpikeDetector() *SpikeDetector {
	return NewSpikeDetector(0.05, 2.5, 0.02)
}

// RMS computes the root-mean-square energy of a buffer of little-endian
// signed 16-bit PCM samples, normalized to [0, 1].
func RMS(pcm []byte) float64 {
	sampleCount := len(pcm) / 2
	if sampleCount == 0 {
		return 0
	}

	var sumSquares float64
	for i := 0; i < sampleCount; i++ {
		sample := int16(binary.LittleEndian.Uint16(pcm[i*2 : i*2+2]))
		normalized := float64(sample) / 32768.0
		sumSquares += normalized * normalized
	}

	return math.Sqrt(sumSquares / float64(sampleCount))
}

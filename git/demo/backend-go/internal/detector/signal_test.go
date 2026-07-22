package detector

import (
	"math"
	"testing"
)

func TestRMS_Silence(t *testing.T) {
	pcm := make([]byte, 32000) // all zeros = silence
	if got := RMS(pcm); got != 0 {
		t.Errorf("RMS(silence) = %v, want 0", got)
	}
}

func TestRMS_FullScale(t *testing.T) {
	pcm := make([]byte, 4)
	// two samples at max positive amplitude (0x7FFF, little-endian)
	pcm[0], pcm[1] = 0xFF, 0x7F
	pcm[2], pcm[3] = 0xFF, 0x7F

	got := RMS(pcm)
	want := 0x7FFF / 32768.0
	if math.Abs(got-want) > 0.001 {
		t.Errorf("RMS(full scale) = %v, want ~%v", got, want)
	}
}

func TestRMS_EmptyBuffer(t *testing.T) {
	if got := RMS(nil); got != 0 {
		t.Errorf("RMS(nil) = %v, want 0", got)
	}
}

func TestFrameDiff_IdenticalFrames(t *testing.T) {
	frame := []byte{10, 20, 30, 40, 50, 60}
	if got := FrameDiff(frame, frame); got != 0 {
		t.Errorf("FrameDiff(identical) = %v, want 0", got)
	}
}

func TestFrameDiff_MaxDifference(t *testing.T) {
	prev := []byte{0, 0, 0}
	curr := []byte{255, 255, 255}

	got := FrameDiff(prev, curr)
	if math.Abs(got-1.0) > 0.001 {
		t.Errorf("FrameDiff(max) = %v, want ~1.0", got)
	}
}

func TestFrameDiff_MismatchedLength(t *testing.T) {
	if got := FrameDiff([]byte{1, 2}, []byte{1, 2, 3}); got != 0 {
		t.Errorf("FrameDiff(mismatched length) = %v, want 0", got)
	}
}

func TestSpikeDetector_NoSpikeOnFirstReading(t *testing.T) {
	d := NewSpikeDetector(0.1, 2.0, 0.01)
	if isSpike, _ := d.Observe(0.5); isSpike {
		t.Error("first reading should never be a spike (sets the baseline)")
	}
}

func TestSpikeDetector_DetectsSpikeAboveBaseline(t *testing.T) {
	d := NewSpikeDetector(0.1, 2.0, 0.01)

	// Establish a quiet baseline.
	for i := 0; i < 20; i++ {
		d.Observe(0.05)
	}

	isSpike, score := d.Observe(0.5) // 10x the baseline
	if !isSpike {
		t.Fatal("expected a spike for a reading far above baseline")
	}
	if score <= 1.0 {
		t.Errorf("expected score > 1.0 for a spike, got %v", score)
	}
}

func TestSpikeDetector_IgnoresBelowMinLevel(t *testing.T) {
	d := NewSpikeDetector(0.1, 2.0, 0.5)

	if isSpike, _ := d.Observe(0.01); isSpike {
		t.Error("first reading is never a spike")
	}
	// Even a "spike" relative to baseline shouldn't fire below minLevel.
	if isSpike, _ := d.Observe(0.05); isSpike {
		t.Error("reading below minLevel should never be flagged as a spike")
	}
}

func TestSpikeDetector_SpikeDoesNotPollutebaseline(t *testing.T) {
	d := NewSpikeDetector(0.5, 2.0, 0.01)

	for i := 0; i < 5; i++ {
		d.Observe(0.1)
	}

	d.Observe(1.0) // spike, should NOT move the baseline

	// A moderately elevated reading should still spike against the old baseline.
	isSpike, _ := d.Observe(0.25)
	if !isSpike {
		t.Error("baseline should not have been dragged up by the spike reading")
	}
}

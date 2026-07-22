package detector

// SpikeDetector maintains a rolling baseline for a signal (audio energy,
// video motion, ...) and flags a "spike" when a new reading exceeds the
// baseline by a multiplier. The baseline is only updated on non-spike
// readings, so a loud/action moment doesn't immediately raise the bar for
// the next one.
type SpikeDetector struct {
	baseline   float64
	alpha      float64 // baseline smoothing factor (0..1, higher = adapts faster)
	multiplier float64 // how far above baseline counts as a spike
	minLevel   float64 // ignore readings below this absolute level (silence/static)
}

func NewSpikeDetector(alpha, multiplier, minLevel float64) *SpikeDetector {
	return &SpikeDetector{
		alpha:      alpha,
		multiplier: multiplier,
		minLevel:   minLevel,
	}
}

// Observe feeds a new reading and reports whether it's a spike, plus a score
// (how many multiples of the baseline it is) usable as a weight.
func (d *SpikeDetector) Observe(value float64) (isSpike bool, score float64) {
	if d.baseline == 0 {
		d.baseline = value
		return false, 0
	}

	isSpike = value > d.baseline*d.multiplier && value > d.minLevel

	if !isSpike {
		d.baseline = d.baseline*(1-d.alpha) + value*d.alpha
		return false, 0
	}

	return true, value / d.baseline
}

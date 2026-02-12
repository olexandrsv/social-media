package tone

type Tone struct {
	positive float64
	negative float64
}

func New(positive, negative float64) *Tone {
	return &Tone{
		positive: positive,
		negative: negative,
	}
}

func (t *Tone) Positive() float64 {
	return t.positive
}

func (t *Tone) Negative() float64 {
	return t.negative
}

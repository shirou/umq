package umq

type ConsumeOption struct {
	isConsume bool
}

func (opt ConsumeOption) Target() []TransportType {
	return []TransportType{TransportAny}
}
func (opt ConsumeOption) Apply(q Queue) error {
	switch t := q.(type) {
	case *MemoryQueue:
		t.IsConsume = opt.isConsume
	case *SQSQueue:
		t.IsConsume = opt.isConsume
	}

	return nil
}
func NewConsumeOption(isConsume bool) Option {
	return &ConsumeOption{isConsume: isConsume}
}

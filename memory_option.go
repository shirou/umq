package umq

type MemoryDummyOption struct{}

func (opt MemoryDummyOption) Target() []TransportType { return []TransportType{TransportMemory} }
func (opt MemoryDummyOption) Apply(q Queue) error {
	return nil
}
func NewMemoryDummyOption() Option {
	return &MemoryDummyOption{}
}

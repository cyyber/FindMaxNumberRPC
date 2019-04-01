package maxnumber

type MaxNumber struct {
	maxNumber int64
	firstRun bool
}

func (m *MaxNumber) FindMaxNumber(value int64) (int64, bool) {
	if m.maxNumber < value || m.firstRun {
		m.firstRun = false
		m.maxNumber = value
		return m.maxNumber, true
	}
	return m.maxNumber, false
}

func (m *MaxNumber) GetMaxNumber() int64 {
	return m.maxNumber
}

func NewMaxNumber() *MaxNumber {
	m := &MaxNumber{0, true}
	return m
}
package jso

type WriterState int

const (
	AtStart WriterState = iota
	AtObjectStart
	AfterObjectKey
	AfterObjectValue
	AtArrayStart
	AfterArrayValue
	AtEnd
)

func (s WriterState) InObject() bool {
	switch s {
	case AtObjectStart, AfterObjectKey, AfterObjectValue:
		return true
	default:
		return false
	}
}
func (s WriterState) InArray() bool {
	switch s {
	case AtArrayStart, AfterArrayValue:
		return true
	default:
		return false
	}
}

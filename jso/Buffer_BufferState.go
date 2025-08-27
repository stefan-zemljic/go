package jso

type BufferState int

const (
	AtStart BufferState = iota
	AtObjectStart
	AfterObjectKey
	AfterObjectValue
	AtArrayStart
	AfterArrayValue
	AtEnd
)

func (s BufferState) InObject() bool {
	switch s {
	case AtObjectStart, AfterObjectKey, AfterObjectValue:
		return true
	default:
		return false
	}
}
func (s BufferState) InArray() bool {
	switch s {
	case AtArrayStart, AfterArrayValue:
		return true
	default:
		return false
	}
}

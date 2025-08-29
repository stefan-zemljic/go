package jso

type Kind int

const (
	EOF Kind = iota
	Null
	Bool
	Number
	String
	Array
	Object
	End
)

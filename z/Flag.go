package z

var errorPrefixFlag = "z.Flag: "
var True = Flag{true}

var False = Flag{false}

type Flag struct {
	V bool
}

func FlagOf(set bool) Flag {
	return Flag{set}
}

func (s Flag) Neg() Flag {
	return Flag{!s.V}
}

func (s Flag) Not() bool { return !s.V }

func (s Flag) Must() bool {
	if !s.V {
		panic(errorPrefixFlag + "flag not set")
	}
	return s.V
}

func (s Flag) MustNot() bool {
	if s.V {
		panic(errorPrefixFlag + "flag set")
	}
	return !s.V
}

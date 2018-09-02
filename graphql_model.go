package glexer

type IntOption interface {
	Get() int
	IsSet() bool
	GetOrElse(int) int
}

type intOptionStruct struct {
	value int
	isSet bool
}

func NewIntOption(v int) IntOption {
	return intOptionStruct{
		value: v,
		isSet: true,
	}
}

var EmptyIntOption = intOptionStruct{
	isSet: false,
}

func (o intOptionStruct) Get() int {
	if !o.isSet {
		panic("IntOption value not set.")
	}
	return o.value
}

func (o intOptionStruct) GetOrElse(v int) int {
	if o.isSet {
		return o.value
	}
	return v
}

func (o intOptionStruct) IsSet() bool {
	return o.isSet
}

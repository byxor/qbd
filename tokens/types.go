package tokens

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Assignment
	LocalReference
	Subtraction
	Addition
	Division
	Multiplication
	EqualityCheck
	LessThanCheck
	LessThanOrEqualCheck
	GreaterThanCheck
	GreaterThanOrEqualCheck
	StartOfStruct
	EndOfStruct
	StartOfArray
	EndOfArray
	StartOfSwitch
	EndOfSwitch
	SwitchCase
	StartOfFunction
	EndOfFunction
	Return
	Break
	StartOfIf
	OptimisedIf
	Else
	ElseIf
	EndOfIf
	Integer
	Float
	Name
	ShortJump
	ChecksumTableEntry
	NamespaceAccess
	Invalid
)

//go:generate stringer -type=Token

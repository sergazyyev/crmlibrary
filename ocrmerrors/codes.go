package ocrmerrors

const (
	EMPTY       Code = "EMPTY"
	INVALID     Code = "INVALID"
	INTERNAL    Code = "INTERNAL"
	ARGISNIL    Code = "ARGUMENTISNIL"
	NOCONFIGURE Code = "NOTCONFIGURED"
)

type Code string

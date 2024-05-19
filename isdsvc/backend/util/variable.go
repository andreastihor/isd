package util

type OptionalBool string

const (
	OptionalBool_UNKNOWN_OptionalBool OptionalBool = "UNKNOWN"
	OptionalBool_FALSE                OptionalBool = "FALSE"
	OptionalBool_TRUE                 OptionalBool = "TRUE"
)

func GetBool(val OptionalBool) OptionalBool {
	if val == OptionalBool_TRUE {
		return OptionalBool_TRUE
	} else if val == OptionalBool_FALSE {
		return OptionalBool_FALSE
	}

	return OptionalBool_UNKNOWN_OptionalBool
}

type Gender string

const (
	Gender_Male   Gender = "MALE"
	Gender_Female Gender = "FEMALE"
)

package utils

import "errors"

func GetArgument(arguments []string, arg string) (string, error) {
	idxArg := 0
	for i, a := range arguments {
		if "--"+arg == a {
			idxArg = i
		}
	}

	if idxArg == 0 {
		return "", errors.New("could not find argument")
	}

	return arguments[idxArg+1], nil
}

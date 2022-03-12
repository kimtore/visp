package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/visp/list"
)

func ErrMsgDataType(actual list.DataType, expected ...list.DataType) error {
	for _, kind := range expected {
		if kind == actual {
			return nil
		}
	}
	types := make([]string, len(expected))
	for i := range expected {
		types[i] = "'" + string(expected[i]) + "'"
	}
	msgPart := strings.Join(types, " or ")
	return fmt.Errorf("selected row is type '%s', need %s", actual, msgPart)
}

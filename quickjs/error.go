package quickjs

import (
	"strconv"
)

type Errno int

func (e Errno) Error() string {
	return "errno: " + strconv.Itoa(int(e))
}

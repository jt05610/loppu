package loppu

import (
	"strconv"
	"strings"
)

type Version string

func splitGet(s, sep string, i int) int {
	m := strings.Split(s, sep)[i]
	conv, err := strconv.Atoi(m)
	if err != nil {
		panic(err)
	}
	return conv
}

func (v Version) Major() int {
	return splitGet(string(v), ".", 0)
}

func (v Version) Minor() int {
	return splitGet(string(v), ".", 1)
}

func (v Version) Rev() int {
	return splitGet(string(v), ".", 2)
}

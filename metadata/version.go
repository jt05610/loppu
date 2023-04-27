package metadata

import (
	"strconv"
	"strings"
)

type Version string

type VersionType uint8

const (
	Major VersionType = 0
	Minor VersionType = 1
	Patch VersionType = 2
)

func splitGet(s, sep string, i int) int {
	m := strings.Split(s, sep)[i]
	conv, err := strconv.Atoi(m)
	if err != nil {
		panic(err)
	}
	return conv
}

func (v Version) Update(i VersionType) Version {
	tokens := strings.Split(string(v), ".")
	value, err := strconv.Atoi(tokens[i])
	if err != nil {
		panic(err)
	}
	tokens[i] = strconv.Itoa(value + 1)
	if int(i) != len(tokens)-1 {
		for j := int(i) + 1; j < len(tokens); j++ {
			tokens[j] = "0"
		}
	}
	return Version(strings.Join(tokens, "."))
}

func (v Version) Major() int {

	return splitGet(string(v), ".", 0)
}

func (v Version) Minor() int {

	return splitGet(string(v), ".", 1)
}

func (v Version) Patch() int {

	return splitGet(string(v), ".", 2)
}

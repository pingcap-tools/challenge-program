package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseStringSlice(input string) []string {
	return strings.Split(input, ",")
}

func ParseIntSlice(input string) []int {
	var s []int
	err := json.Unmarshal([]byte(fmt.Sprintf("[%s]", input)), &s)
	if err != nil {
		return []int{}
	}
	return s
}

func EncodeStringSlice(input []string) string {
	return strings.Join(input, ",")
}
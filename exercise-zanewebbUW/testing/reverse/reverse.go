package reverse

import "strings"

// TODO: Implement this method
//Reverse returns the reverse of the string passed as `s`.
//For example, if `s` is "abcd" this will return "dcba".
func Reverse(s string) string {
	if len(s) == 0 {
		return s
	}

	splitS := strings.Split(s, "")
	length := len(splitS)
	newSplitS := make([]string, length)
	for i, char := range splitS {
		newSplitS[(length-1)-i] = char
	}
	newS := strings.Join(newSplitS, "")
	return newS
}

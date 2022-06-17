package utils

import "testing"

func TestGenerateRandomString(t *testing.T) {
	println(GenerateRandomBytesToHex(3))
}

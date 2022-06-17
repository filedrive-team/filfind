package jwttoken

import "testing"

func TestGenerateTokenSecret(t *testing.T) {
	ts, err := GenerateTokenSecret(KTJwtHS256)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := TokenSecretEncode(ts)
	if err != nil {
		t.Fatal(err)
	}
	println(encoded)
}

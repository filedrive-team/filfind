package utils

import (
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math"
	"math/rand"
	"time"
)

// GenerateRandNumStr Generates a random numeric string
func GenerateRandNumStr(size int) string {
	if size <= 0 {
		return ""
	}
	return fmt.Sprint(RandInt(int(math.Pow10(size-1)), int(math.Pow10(size))))
}

// GenRandStr Generates a random string, space contains number, uppercase and lowercase letters
func GenRandStr(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

func GenerateClientId() string {
	return uuid.NewV4().String()
}

func GenerateRandomBytes(n int) ([]byte, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	_, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomBytesToHex(bytes int) (string, error) {
	b, err := GenerateRandomBytes(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

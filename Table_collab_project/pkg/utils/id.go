package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

func GenerateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 8

	result := make([]byte, idLength)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

func GenerateColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#FFD166", "#06D6A0",
		"#118AB2", "#073B4C", "#EF476F", "#7209B7",
	}

	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(colors))))
	return colors[num.Int64()]
}

func GenerateTimestamp() int64 {
	return time.Now().UnixMilli()
}

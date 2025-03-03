package xotp

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"time"
)

const otpString = "0123456789"

func genOTP(length int) string {
	byteSlice := make([]byte, length)
	otpStringLength := len(otpString)
	otpStringLengthBigInt := big.NewInt(int64(otpStringLength))
	mathRand := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		// use crypto to generate otp digit
		num, err := rand.Int(rand.Reader, otpStringLengthBigInt)
		if err == nil {
			byteSlice[i] = otpString[num.Int64()]
			continue
		}
		// if crypto fails, use math package to generate otp digit instead
		mathNum := mathRand.Intn(otpStringLength)
		byteSlice[i] = otpString[mathNum]
	}
	return string(byteSlice)
}

func Gen6DigitsOTP() string {
	return genOTP(6)
}

func Gen4DigitsOTP() string {
	return genOTP(4)
}

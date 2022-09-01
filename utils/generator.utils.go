package utils

// #nosec

import (
	cryptoRand "crypto/rand"
	"fmt"
	"gin-starter/common/constant"
	"io"
	"math/rand"
	"time"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// GenerateRandomNumber generate random number with threshold
func GenerateRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(max-min) + min // #nosec

	return num
}

// GenerateTrxID generate transaction ID
func GenerateTrxID(prefix string) string {
	res := prefix
	rand.Seed(time.Now().UnixNano())                                                     // #nosec
	num := rand.Intn(constant.NinetyNineHundred-constant.Hundred) + constant.TenThousand // #nosec
	res = res + time.Now().Format("20060102") + "/" + fmt.Sprint(num)

	return res
}

// GenerateExternalID generate external ID. Commonly used for third party payment
func GenerateExternalID(prefix string) string {
	res := prefix + fmt.Sprint(time.Now().Unix())

	return res
}

// GenerateOTP generate OTP number for user
func GenerateOTP(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(cryptoRand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes generate random string by bytes
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))] // #nosec
	}
	return string(b)
}

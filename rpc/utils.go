package rpc

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// IntToHex convert int to hexadecimal representation
func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}

// ParseBigInt parse hex string value to big.Int
func ParseBigInt(value string) (big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(value, &i)

	return i, err
}

// ParseInt parse hex string value to int
func ParseInt(value string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(value, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func newBigInt(s string) big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return *i
}

func ptrInt(i int) *int {
	return &i
}

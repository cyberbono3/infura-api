package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntToHex(t *testing.T) {
	tt := []struct {
		name     string
		input    int
		expected string
	}{
		{"10000", 100, "0x64"},
		{"100000", 100000, "0x186a0"},
		{"1000000000000000000", 1000000000000000000, "0xde0b6b3a7640000"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := IntToHex(tc.input)
			if actual != tc.expected {
				t.Fatalf("Unable to convert int to hex, expected: %s, got: %s", tc.expected, actual)
			}
		})
	}
}

func TestParseBigInt(t *testing.T) {
	i, err := ParseBigInt("0xabc")
	assert.Nil(t, err)
	assert.Equal(t, int64(2748), i.Int64())

	_, err = ParseBigInt("*1")
	assert.NotNil(t, err)

	_, err = ParseBigInt("")
	assert.NotNil(t, err)
}

func TestParseInt(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		expected int
		err      string
	}{
		{name: "valid input", value: "0xde0b6b3a7640000", expected: 1000000000000000000},
		{name: "not a hex string", value: "x", err: "not a hex string"},
		{name: "empty input", value: "", err: "empty input "},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			outputInt, err := ParseInt(tc.value)
			if tc.err != "" {
				assert.NotNil(t, err)
				assert.Equal(t, outputInt, 0)
				return
			}
			assert.Nil(t, err)
			assert.NotEqual(t, outputInt, 0)
			if outputInt != tc.expected {
				t.Fatalf("Unable to convert hex string to int, expected: %d, got: %d", tc.expected, outputInt)
			}
		})
	}
}

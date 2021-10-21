package data

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionDataMissingIndexReturnsErr(t *testing.T) {
	td := TransactionData{
		BlockNumber: 114996,
	}

	v := NewValidation()
	err := v.Validate(td)
	assert.Len(t, err, 1)
}

func TestTransactionDataNegativeIndexReturnsErr(t *testing.T) {
	td := TransactionData{
		BlockNumber:      114996,
		TransactionIndex: -1,
	}

	v := NewValidation()
	err := v.Validate(td)
	assert.Len(t, err, 1)
}

func TestTransactionDataReturnsNoErr(t *testing.T) {
	td := TransactionData{
		BlockNumber:      114996,
		TransactionIndex: 10,
	}

	v := NewValidation()
	err := v.Validate(td)
	assert.Len(t, err, 0)
}

func TestTransactionDataNegativeBlockReturnsErr(t *testing.T) {
	td := TransactionData{
		BlockNumber:      -137377,
		TransactionIndex: 10,
	}

	v := NewValidation()
	err := v.Validate(td)
	assert.Len(t, err, 1)
}

func TestTransactionDataMissingBlockReturnsErr(t *testing.T) {
	td := TransactionData{
		TransactionIndex: 10,
	}

	v := NewValidation()
	err := v.Validate(td)
	assert.Len(t, err, 1)
}

func TestTransactionDataToJSONReturnsNoErr(t *testing.T) {
	td := TransactionData{
		BlockNumber:      373773,
		TransactionIndex: 10,
	}

	b := bytes.NewBufferString("")
	err := ToJSON(td, b)
	assert.NoError(t, err)
}

func TestBlockDataMissingBlockReturnsErr(t *testing.T) {
	bd := BlockData{
		ShowTransFlag: false,
	}

	v := NewValidation()
	err := v.Validate(bd)
	assert.Len(t, err, 1)
}

func TestBlockDataNegativeBlockReturnsErr(t *testing.T) {
	bd := BlockData{
		BlockNumber:   -12837,
		ShowTransFlag: false,
	}

	v := NewValidation()
	err := v.Validate(bd)
	assert.Len(t, err, 1)
}

func TestBlockDataMissingShowTransFlagReturnsNoErr(t *testing.T) {
	bd := BlockData{
		BlockNumber: 12837,
	}

	v := NewValidation()
	err := v.Validate(bd)
	assert.Len(t, err, 0)

}

func TestBlockDataToJSONReturnsNoErr(t *testing.T) {
	bd := BlockData{
		BlockNumber:   373773,
		ShowTransFlag: false,
	}

	b := bytes.NewBufferString("")
	err := ToJSON(bd, b)
	assert.NoError(t, err)
}

package encoder

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultLen = 10

func FuzzIDEncoderEncodeID(f *testing.F) {
	sut := NewIDEncoder()
	results := make(map[uint]string)
	var mutex sync.Mutex
	checkUniqueEncodedFunc := func(t *testing.T, id uint, encodedID string) {
		t.Helper()

		mutex.Lock()
		defer mutex.Unlock()
		previousValue, contains := results[id]
		if !contains {
			results[id] = encodedID
			return
		}

		if previousValue != encodedID {
			assert.Fail(t, "Different values for one id")
		}
	}

	f.Fuzz(func(t *testing.T, id uint) {
		result := sut.EncodeID(id, defaultLen)
		assert.GreaterOrEqual(t, len(result), defaultLen, "Encoded string length is less then requested")
		checkUniqueEncodedFunc(t, id, result)
	})
}

func TestIdEncoder_EncodeID(t *testing.T) {
	sut := NewIDEncoder()

	tests := []struct {
		name                 string
		firstIDInput         uint
		secondIDInput        uint
		shouldResultsBeEqual bool
	}{
		{
			name:                 "same inputs must return same output value",
			firstIDInput:         5,
			secondIDInput:        5,
			shouldResultsBeEqual: true,
		},
		{
			name:                 "different inputs must return different output values",
			firstIDInput:         5,
			secondIDInput:        8,
			shouldResultsBeEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1 := sut.EncodeID(tt.firstIDInput, defaultLen)
			result2 := sut.EncodeID(tt.secondIDInput, defaultLen)

			if tt.shouldResultsBeEqual {
				assert.Equal(t, result1, result2)
			} else {
				assert.NotEqual(t, result1, result2)
			}

			resultLessThanDefaultMessage := "Encoded string length is less then requested"
			assert.GreaterOrEqual(t, len(result1), defaultLen, resultLessThanDefaultMessage)
			assert.GreaterOrEqual(t, len(result2), defaultLen, resultLessThanDefaultMessage)
		})
	}
}

func Test_charSet(t *testing.T) {
	expectedSet := []byte("_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := baseCharSet()
	assert.Equal(t, expectedSet, result)
}

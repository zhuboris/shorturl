// Package encoder provides a type for encoding integers into strings to generate unique identifiers for id values.
package encoder

// IDEncoder is an interface that describes type encoding integer into string.
// It can be used to mock that internal type.
type IDEncoder interface {
	EncodeID(id, minLen uint) string
}

// idEncoder is an encoder of integer id into string with base encoding.
// It is using symbols from his baseChars byte slice.
//
// Each id is mapped to one string value, it guarantees no collisions or extending length
// up to len(baseChars) in power of asked string length.
// For example, for 63 symbols and string length equal to 10, it will be able to generate
// 63^10 different encoded strings.
type idEncoder struct {
	baseChars []byte
}

// NewIDEncoder initializes instance of idEncoder and returns it as IDEncoder interface.
// It initializes encoder's set of symbols with baseCharSet function.
// The slice could be set manually, but this function can avoid some possible human mistakes.
func NewIDEncoder() IDEncoder {
	return idEncoder{
		baseChars: baseCharSet(),
	}
}

// EncodeID is a method that encodes id into string with selected minimum length.
// The length of the result can be more than asked, caller should handle it himself.
//
// The method accepts uint values to avoid passing incorrect values less than 0 and
// to increase the range of possible ids.
//
// Result string is reversed base value of inputted id. However, it does not matter because no decoding is planned.
//
// Base encoded string length can be less than the requested minimum, in this case it will add zero values
// at the beginning of base string (actually in the end because it is reversed).
func (e idEncoder) EncodeID(id, minLen uint) string {
	encodedID := e.baseEncode(id, minLen)
	encodedID = e.extendResultIfNeeded(minLen, encodedID)
	return string(encodedID)
}

func (e idEncoder) extendResultIfNeeded(minLen uint, encodedID []byte) []byte {
	missingLen := int(minLen) - len(encodedID)
	if missingLen <= 0 {
		return encodedID
	}

	toAdd := e.repeatedZeroChar(missingLen)
	encodedID = append(encodedID, toAdd...)
	return encodedID

}

func (e idEncoder) baseEncode(id, minLen uint) []byte {
	var baseSymbolsCount = uint(len(e.baseChars))

	result := make([]byte, 0, minLen)
	for id > 0 {
		remainder := id % baseSymbolsCount
		result = append(result, e.baseChars[remainder])
		id /= baseSymbolsCount
	}

	return result
}

func (e idEncoder) repeatedZeroChar(repeatsCount int) []byte {
	result := make([]byte, repeatsCount)
	for i := range result {
		result[i] = e.baseChars[0]
	}

	return result
}

// baseCharSet returns filled slice with desired symbols to encode in bytes.
// It is using ranges 0-9, A-Z, a-z and '_' symbol.
func baseCharSet() []byte {
	const charsOptionsCount = 63

	result := make([]byte, 0, charsOptionsCount)
	result = append(result, '_')
	result = addCharsFromRange(result, '0', '9')
	result = addCharsFromRange(result, 'a', 'z')
	result = addCharsFromRange(result, 'A', 'Z')

	return result
}

func addCharsFromRange(chars []byte, from, to byte) []byte {
	for i := from; i <= to; i++ {
		chars = append(chars, i)
	}

	return chars
}

package NeverScript

import (
	"bytes"
	"github.com/pkg/errors"
	goErrors "errors"
)

type ByteCode struct {
	content []byte
	length  int
}

func NewByteCode(content []byte) ByteCode {
	return ByteCode{
		content: content,
		length:  len(content),
	}
}

func NewEmptyByteCode() ByteCode {
	return NewByteCode([]byte{})
}

func (this *ByteCode) Push(bytes ...byte) {
	this.content = append(this.content, bytes...)
	this.length += len(bytes)
}

func (this ByteCode) GetSlice(startIndex, endIndex int) (ByteCode, error) {
	if startIndex < 0 {
		return nilByteCode, indexOutOfRange
	}

	if startIndex >= endIndex {
		return nilByteCode, indexOutOfRange
	}

	if endIndex > this.length {
		return nilByteCode, indexOutOfRange
	}

	content := this.content[startIndex:endIndex]
	return NewByteCode(content), nil
}

func (this ByteCode) ToBytes() []byte {
	return this.content
}

func (this ByteCode) IsLongerThan(other ByteCode) bool {
	return this.length > other.length
}

func (this ByteCode) IsShorterThan(other ByteCode) bool {
	return this.length < other.length
}

func (this ByteCode) IsSameLengthAs(other ByteCode) bool {
	return this.length == other.length
}

func (this ByteCode) IsEqualTo(other ByteCode) bool {
	return bytes.Equal(this.content, other.content)
}

func (this ByteCode) Contains(other ByteCode) (bool, error) {
	if other.IsLongerThan(this) {
		return false, nil
	}

	iterateUpTo := this.length - other.length + 1

	for i := 0; i < iterateUpTo; i++ {
		sliceOfThis, err := this.GetSlice(i, i + other.length)
		if err != nil {
			return false, errors.Wrap(err, "Failed to get slice of bytecode")
		}

		if sliceOfThis.IsEqualTo(other) {
			return true, nil
		}
	}

	return false, nil
}

func (this ByteCode) IsEqualTo_IgnoreByte(other ByteCode, byteToIgnore byte) bool {
	forceFutureComparisonsToPass(other, &this, byteToIgnore)
	return this.IsEqualTo(other)
}

func (this ByteCode) Contains_IgnoreByte(other ByteCode, byteToIgnore byte) (bool, error) {
	forceFutureComparisonsToPass(other, &this, byteToIgnore)
	return this.Contains(other)
}

func forceFutureComparisonsToPass(byteCodeToCheck ByteCode, byteCodeToModify *ByteCode, byteToIgnore byte) {
	length := min(byteCodeToCheck.length, byteCodeToModify.length)

	for i := 0; i < length; i++ {
		byteToCheck := byteCodeToCheck.content[i]

		if byteToCheck == byteToIgnore {
			// This will force the comparison checks to pass.
			byteCodeToModify.content[i] = byteToIgnore
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

var (
	indexOutOfRange = goErrors.New("Index is out of range")
	nilByteCode = NewEmptyByteCode()
)
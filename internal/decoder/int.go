package decoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type intDecoder struct {
	typ        *runtime.Type
	kind       reflect.Kind
	op         func(unsafe.Pointer, int64)
	structName string
	fieldName  string
}

func newIntDecoder(typ *runtime.Type, structName, fieldName string, op func(unsafe.Pointer, int64)) *intDecoder {
	return &intDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		op:         op,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *intDecoder) typeError(buf []byte, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   runtime.RType2Type(d.typ),
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

var (
	pow10i64 = [...]int64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
	pow10i64Len = len(pow10i64)
)

func (d *intDecoder) parseInt(b []byte) (int64, error) {
	isNegative := false
	if b[0] == '-' {
		b = b[1:]
		isNegative = true
	}
	maxDigit := len(b)
	if maxDigit > pow10i64Len {
		return 0, fmt.Errorf("invalid length of number")
	}
	sum := int64(0)
	for i := 0; i < maxDigit; i++ {
		c := int64(b[i]) - 48
		digitValue := pow10i64[maxDigit-i-1]
		sum += c * digitValue
	}
	if isNegative {
		return -1 * sum, nil
	}
	return sum, nil
}

var (
	numTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
)

var (
	numZeroBuf = []byte{'0'}
)

func (d *intDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	if char(b, cursor) != 'i' {
		return nil, cursor, errors.ErrExpected("int", cursor)
	}

	cursor++
	if char(b, cursor) != ':' {
		return nil, cursor, errors.ErrExpected("int sep ':'", cursor)
	}
	cursor++

	switch char(b, cursor) {
	case '0':
		cursor++
		if char(b, cursor) != ';' {
			return nil, cursor, errors.ErrExpected("';' end int", cursor)
		}
		return numZeroBuf, cursor + 1, nil
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start := cursor
		cursor++
		for numTable[char(b, cursor)] {
			cursor++
		}
		if char(b, cursor) != ';' {
			return nil, cursor, errors.ErrExpected("';' end int", cursor)
		}
		num := buf[start:cursor]
		return num, cursor + 1, nil
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil
	default:
		return nil, 0, d.typeError([]byte{char(b, cursor)}, cursor)
	}
}

func (d *intDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c

	return d.processBytes(bytes, cursor, p)
}

func (d *intDecoder) processBytes(bytes []byte, cursor int64, p unsafe.Pointer) (int64, error) {
	i64, err := d.parseInt(bytes)
	if err != nil {
		return 0, d.typeError(bytes, cursor)
	}
	switch d.kind {
	case reflect.Int8:
		if i64 < -1*(1<<7) || (1<<7) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Int16:
		if i64 < -1*(1<<15) || (1<<15) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Int32:
		if i64 < -1*(1<<31) || (1<<31) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	}

	d.op(p, i64)

	return cursor, nil
}

func readInt(buf []byte, cursor int64) (int, int64, error) {
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	if char(b, cursor) != 'i' {
		return 0, cursor, errors.ErrExpected("'i' to start a int", cursor)
	}

	cursor++
	if char(b, cursor) != ':' {
		return 0, cursor, errors.ErrExpected("int sep ':'", cursor)
	}
	cursor++

	switch char(b, cursor) {
	case '0':
		cursor++
		if char(b, cursor) != ';' {
			return 0, cursor, errors.ErrExpected("';' end int", cursor)
		}
		cursor++
		return 0, cursor, nil
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start := cursor
		cursor++
		for numTable[char(b, cursor)] {
			cursor++
		}

		if char(b, cursor) != ';' {
			return 0, cursor, errors.ErrExpected("';' end int", cursor)
		}
		value := parseByteStringInt(buf[start:cursor])
		cursor++
		return value, cursor, nil
	default:
		return 0, 0, errors.ErrExpected("int", cursor)
	}
}

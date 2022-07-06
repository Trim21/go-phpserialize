package decoder

import (
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type numberDecoder struct {
	stringDecoder *stringDecoder
	op            func(unsafe.Pointer, json.Number)
	structName    string
	fieldName     string
}

func newNumberDecoder(structName, fieldName string, op func(unsafe.Pointer, json.Number)) *numberDecoder {
	return &numberDecoder{
		stringDecoder: newStringDecoder(structName, fieldName),
		op:            op,
		structName:    structName,
		fieldName:     fieldName,
	}
}

func (d *numberDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if _, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&bytes)), 64); err != nil {
		return 0, errors.ErrSyntax(err.Error(), c)
	}
	cursor = c
	s := *(*string)(unsafe.Pointer(&bytes))
	d.op(p, json.Number(s))
	return cursor, nil
}

func (d *numberDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
			for floatTable[buf[cursor]] {
				cursor++
			}
			num := buf[start:cursor]
			return num, cursor, nil
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return nil, 0, err
			}
			cursor += 4
			return nil, cursor, nil
		case '"':
			return d.stringDecoder.decodeByte(buf, cursor)
		default:
			return nil, 0, errors.ErrUnexpectedEnd("json.Number", cursor)
		}
	}
}

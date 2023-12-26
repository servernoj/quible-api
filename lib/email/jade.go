// Code generated by "jade.go"; DO NOT EDIT.

package email

import (
	"bytes"
	"io"
	"strconv"
)

var (
	escaped   = []byte{'<', '>', '"', '\'', '&'}
	replacing = []string{"&lt;", "&gt;", "&#34;", "&#39;", "&amp;"}
)

func WriteEscString(st string, buffer *WriterAsBuffer) {
	for i := 0; i < len(st); i++ {
		if n := bytes.IndexByte(escaped, st[i]); n >= 0 {
			buffer.WriteString(replacing[n])
		} else {
			buffer.WriteByte(st[i])
		}
	}
}

type WriterAsBuffer struct {
	io.Writer
}

func (w *WriterAsBuffer) WriteString(s string) (n int, err error) {
	n, err = w.Write([]byte(s))
	return
}

func (w *WriterAsBuffer) WriteByte(b byte) (err error) {
	_, err = w.Write([]byte{b})
	return
}

type stringer interface {
	String() string
}

func WriteAll(a interface{}, escape bool, buffer *WriterAsBuffer) {
	switch v := a.(type) {
	case string:
		if escape {
			WriteEscString(v, buffer)
		} else {
			buffer.WriteString(v)
		}
	case int:
		WriteInt(int64(v), buffer)
	case int8:
		WriteInt(int64(v), buffer)
	case int16:
		WriteInt(int64(v), buffer)
	case int32:
		WriteInt(int64(v), buffer)
	case int64:
		WriteInt(v, buffer)
	case uint:
		WriteUint(uint64(v), buffer)
	case uint8:
		WriteUint(uint64(v), buffer)
	case uint16:
		WriteUint(uint64(v), buffer)
	case uint32:
		WriteUint(uint64(v), buffer)
	case uint64:
		WriteUint(v, buffer)
	case float32:
		buffer.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 64))
	case float64:
		buffer.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		WriteBool(v, buffer)
	case stringer:
		if escape {
			WriteEscString(v.String(), buffer)
		} else {
			buffer.WriteString(v.String())
		}
	default:
		buffer.WriteString("\n<<< unprinted type, fmt.Stringer implementation needed >>>\n")
	}
}

func ternary(condition bool, iftrue, iffalse interface{}) interface{} {
	if condition {
		return iftrue
	} else {
		return iffalse
	}
}

// Used part of go source:
// https://github.com/golang/go/blob/master/src/strconv/itoa.go
func WriteUint(u uint64, buffer *WriterAsBuffer) {
	var a [64 + 1]byte
	i := len(a)

	if ^uintptr(0)>>32 == 0 {
		for u > uint64(^uintptr(0)) {
			q := u / 1e9
			us := uintptr(u - q*1e9)
			for j := 9; j > 0; j-- {
				i--
				qs := us / 10
				a[i] = byte(us - qs*10 + '0')
				us = qs
			}
			u = q
		}
	}

	us := uintptr(u)
	for us >= 10 {
		i--
		q := us / 10
		a[i] = byte(us - q*10 + '0')
		us = q
	}

	i--
	a[i] = byte(us + '0')
	buffer.Write(a[i:])
}
func WriteInt(i int64, buffer *WriterAsBuffer) {
	if i < 0 {
		buffer.WriteByte('-')
		i = -i
	}
	WriteUint(uint64(i), buffer)
}
func WriteBool(b bool, buffer *WriterAsBuffer) {
	if b {
		buffer.WriteString("true")
		return
	}
	buffer.WriteString("false")
}

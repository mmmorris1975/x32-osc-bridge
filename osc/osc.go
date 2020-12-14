package osc

import (
	"bytes"
	"errors"
	"io"
)

var ErrInvalidMessageType = errors.New("invalid message type")

func WriteString(s string) []byte {
	msg := append([]byte(s), 0x0)

	if len(msg)%4 > 0 {
		padding := 4 - len(msg)%4
		msg = append(msg, make([]byte, padding)...)
	}

	return msg
}

func ReadString(buf *bytes.Buffer) (string, error) {
	s, err := buf.ReadBytes(0x0)
	if err != nil && errors.Is(err, io.EOF) {
		return "", err
	}

	if len(s)%4 > 0 {
		padding := 4 - len(s)%4
		buf.Next(padding)
	}

	return string(bytes.TrimSuffix(s, []byte{0x0})), nil
}

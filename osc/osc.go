package osc

import (
	"bytes"
	"errors"
	"io"
)

// ErrInvalidMessageType is the message returned if the OSC data doesn't match the message type when doing UnmarshalBinary
var ErrInvalidMessageType = errors.New("invalid message type")

// WriteString returns a properly padded OSC string for the provided input stream
func WriteString(s string) []byte {
	msg := append([]byte(s), 0x0)

	if len(msg)%4 > 0 {
		padding := 4 - len(msg)%4
		msg = append(msg, make([]byte, padding)...)
	}

	return msg
}

// ReadString reads a string from the provided OSC data, stripping the null padding bytes.  The buffer position will
// be incremented to just after the padding of this string field so the next read can read the next data type in the
// message.
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

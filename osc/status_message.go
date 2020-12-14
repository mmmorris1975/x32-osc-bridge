package osc

import (
	"bytes"
	"net"
	"strings"
)

const statusMsg = "/status"

type Status struct {
	State string
	IP    net.IP
	Name  string
}

func (t *Status) MarshalBinary() ([]byte, error) {
	msg := WriteString(statusMsg)

	if t.IP != nil {
		msg = append(msg, WriteString(",sss")...)
		msg = append(msg, WriteString(t.State)...)
		msg = append(msg, WriteString(t.IP.String())...)
		msg = append(msg, WriteString(t.Name)...)
	}

	return msg, nil
}

func (t *Status) UnmarshalBinary(data []byte) error {
	path := append(WriteString(statusMsg), ',')
	if !bytes.HasPrefix(data, path) {
		return ErrInvalidMessageType
	}

	parts := bytes.SplitN(data, []byte(","), 2)
	buf := bytes.NewBuffer(append([]byte{','}, parts[1]...))

	fieldMap, err := ReadString(buf)
	if err != nil {
		return err
	}
	fieldMap = strings.TrimPrefix(fieldMap, ",")

	// these are all strings, no need to type check
	attrs := make([]string, len(fieldMap))
	for i := range fieldMap {
		var v string
		v, err = ReadString(buf)
		if err != nil {
			return err
		}
		attrs[i] = v
	}

	t.State = attrs[0]
	t.IP = net.ParseIP(attrs[1])
	t.Name = attrs[2]

	return nil
}

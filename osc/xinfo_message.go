package osc

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const xinfoMsg = `/xinfo`

type XInfo struct {
	IP      net.IP
	Name    string
	Model   string
	Version string
}

func (t *XInfo) String() string {
	return fmt.Sprintf("&{IP: %s Name: %s Model: %s Version: %s}", t.IP, t.Name, t.Model, t.Version)
}

func (t *XInfo) MarshalBinary() ([]byte, error) {
	msg := WriteString(xinfoMsg)

	if t.IP != nil {
		msg = append(msg, WriteString(",ssss")...)
		msg = append(msg, WriteString(t.IP.String())...)
		msg = append(msg, WriteString(t.Name)...)
		msg = append(msg, WriteString(t.Model)...)
		msg = append(msg, WriteString(t.Version)...)
	}

	return msg, nil
}

func (t *XInfo) UnmarshalBinary(data []byte) error {
	path := append(WriteString(xinfoMsg), ',')
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

	t.IP = net.ParseIP(attrs[0])
	t.Name = attrs[1]
	t.Model = attrs[2]
	t.Version = attrs[3]

	return nil
}

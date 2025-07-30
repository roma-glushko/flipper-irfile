package flipperirlib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type SignalType string

type Protocol string

const (
	ProtocolNEC Protocol = "NEC"
	ProtocolRC5 Protocol = "RC5"
)

const (
	SignalTypeParsed SignalType = "parsed"
	SignalTypeRaw    SignalType = "raw"
)

type SignalLib struct {
	Filetype string
	Version  string
	Signals  []Signal
}

type Signal struct {
	Name string
	Type SignalType

	// Parsed Data
	Protocol string
	Address  uint32
	Command  uint32

	// Raw Data
	Frequency int
	DutyCycle float64
	Data      []int
}

func Unmarshal(s []byte) (*SignalLib, error) {
	lib := SignalLib{}
	signals := make([]Signal, 0, 10)

	var curr Signal

	lines := bytes.Split(s, []byte("\n"))

	for lineno, line := range lines {
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if lib.Filetype == "" && bytes.HasPrefix(line, []byte("Filetype:")) {
			lib.Filetype = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Filetype:"))))
			continue
		}

		if lib.Version == "" && bytes.HasPrefix(line, []byte("Version:")) {
			lib.Version = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Version:"))))
			continue
		}

		if bytes.Equal(line, []byte("#")) {
			if curr.Name != "" {
				signals = append(signals, curr)
				curr = Signal{}
			}

			continue
		}

		parts := bytes.SplitN(line, []byte(":"), 2)

		if len(parts) != 2 {
			continue
		}

		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))

		switch key {
		case "name":
			curr.Name = value
		case "type":
			curr.Type = SignalType(value)
		case "protocol":
			curr.Protocol = value
		case "address":
			addr, err := leHexToUint32(value)

			if err != nil {
				return nil, fmt.Errorf("invalid address at line %d: %v", lineno, err)
			}

			curr.Address = addr
		case "command":
			cmd, err := leHexToUint32(value)

			if err != nil {
				return nil, fmt.Errorf("invalid command at line %d: %v", lineno, err)
			}

			curr.Command = cmd
		case "frequency":
			freq, err := strconv.Atoi(value)

			if err != nil {
				return nil, fmt.Errorf("invalid frequency at line %d: %v", lineno+1, err)
			}

			curr.Frequency = freq
		case "duty_cycle":
			duty, err := strconv.ParseFloat(value, 64)

			if err != nil {
				return nil, fmt.Errorf("invalid duty_cycle at line %d: %v", lineno+1, err)
			}

			curr.DutyCycle = duty
		case "data":
			fields := strings.Fields(value)
			ints := make([]int, len(fields))

			for i, f := range fields {
				n, err := strconv.Atoi(f)
				if err != nil {
					return nil, fmt.Errorf("invalid data int at line %d: %v", lineno+1, err)
				}
				ints[i] = n
			}

			curr.Data = ints
		}
	}

	if curr.Name != "" {
		signals = append(signals, curr)
	}

	lib.Signals = signals

	return &lib, nil

}

func Marshal(l *SignalLib) ([]byte, error) {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "Filetype: %s\n", l.Filetype)
	_, _ = fmt.Fprintf(&buf, "Version: %s\n", l.Version)

	for _, s := range l.Signals {

		buf.WriteString("#\n")

		_, _ = fmt.Fprintf(&buf, "name: %s\n", s.Name)
		_, _ = fmt.Fprintf(&buf, "type: %s\n", s.Type)

		switch s.Type {
		case SignalTypeParsed:
			_, _ = fmt.Fprintf(&buf, "protocol: %s\n", s.Protocol)
			_, _ = fmt.Fprintf(&buf, "address: %s\n", encodeLEUint32Hex(s.Address))
			_, _ = fmt.Fprintf(&buf, "command: %s\n", encodeLEUint32Hex(s.Command))
		case SignalTypeRaw:
			_, _ = fmt.Fprintf(&buf, "frequency: %d\n", s.Frequency)
			_, _ = fmt.Fprintf(&buf, "duty_cycle: %.6f\n", s.DutyCycle)
			buf.WriteString("data:")

			for _, v := range s.Data {
				_, _ = fmt.Fprintf(&buf, " %d", v)
			}

			buf.WriteString("\n")
		}
	}

	return buf.Bytes(), nil
}

func leHexToUint32(s string) (uint32, error) {
	b := strings.Fields(s)

	if len(b) != 4 {
		return 0, fmt.Errorf("expected 4 bytes, got %d", len(b))
	}

	var result uint32

	for i := 0; i < 4; i++ {
		b, err := strconv.ParseUint(b[i], 16, 8)
		if err != nil {
			return 0, err
		}

		result |= uint32(b) << (8 * i)
	}

	return result, nil
}

func encodeLEUint32Hex(val uint32) string {
	return fmt.Sprintf("%02X %02X %02X %02X", byte(val), byte(val>>8), byte(val>>16), byte(val>>24))
}

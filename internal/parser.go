package internal

import (
	"bytes"
	"fmt"
)

func ParseCommand(input []byte) (Command, error) {
	input = bytes.TrimSuffix(input, []byte("\r\n"))
	// TODO: make this work with other ASCII whitespace characters
	name, arguments, found := bytes.Cut(input, []byte(" "))
	if !found {
		return nil, fmt.Errorf("no command provived in input: %v", string(input))
	}
	cmdName := string(name)

	var cmd Command
	switch cmdName {
	case "add":
		cmd = &Add{}
	case "delete":
		cmd = &Delete{}
	case "get":
		// TODO: change get to work on byte keys instead of atrings to avoid extra
		// conversions
		rawKeys := bytes.Split(arguments, []byte(" "))
		if len(rawKeys) != 1 {
			return nil, fmt.Errorf("must provide at least 1 key for 'get' command")
		}
		keys := make([]string, 0, len(rawKeys))
		for _, key := range rawKeys[1:] {
			keys = append(keys, string(key))
		}
		return &Get{keys}, nil
	case "replace":
		cmd = &Replace{}
	case "set":
		cmd = &Set{}
	default:
		return nil, fmt.Errorf("unknown command '%s'", cmdName)
	}
	return cmd, nil
}

package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

const expirationUnixTimestampThreshold = 60 * 60 * 24 * 30

// threshold after which the expiration is treated as a Unix time in seconds rather
// than an offset in seconds from the current time

func ParseCommand(input []byte) (Command, error) {
	input = bytes.TrimSuffix(input, []byte("\r\n"))
	inputLines := bytes.Split([]byte("\r\n"))
	commandLine := string(inputLines[0])
	// TODO: make this work with other ASCII whitespace characters
	commandItems := strings.Split(commandLine, " ")
	if len(commandItems) <= 2 {
		return nil, fmt.Errorf("no command provided in input: %s", commandLine)
	}
	commandName := commandItems[0]

	var cmd Command
	switch commandName {
	case "add":
		arguments := commandItems[1:]
		if len(arguments) < 4 {
			return nil, fmt.Errorf("invalid 'add' command: %s", commandLine)
		}
		if len(inputLines) != 2 {
			return nil, fmt.Errorf("missing data block for 'add' command")
		}
		dataBlock := inputLines[1]
		key := arguments[0]
		exptime, err := strconv.Atoi(arguments[2])
		if err != nil {
			return nil, fmt.Errorf("invalid expiration time", exptime)
		}
		var expiresAt time.Time
		if exptime <= expirationUnixTimestampThreshold {
			expiresAt = util.NowUTC().Add(time.Duration(exptime) * time.Second)
		} else {
			expiresAt = time.Unix(int64(exptime), 0)
		}
		cmd = &Add{
			key:       key,
			value:     dataBlock,
			expiresAt: expiresAt,
		}
	case "delete":
		key := commandItems[1]
		cmd = &Delete{key}
	case "get":
		keys := commandItems[1:]
		if len(keys) != 1 {
			return nil, fmt.Errorf("must provide at least 1 key for 'get' command")
		}
		return &Get{keys}, nil
	case "replace":
		cmd = &Replace{}
	case "set":
		cmd = &Set{}
	default:
		return nil, fmt.Errorf("unknown command: '%s'", commandName)
	}
	return cmd, nil
}

func parseStorageCommandElements(arguments []string) (key string, expiresAt time.Time, err error) {
	return
}

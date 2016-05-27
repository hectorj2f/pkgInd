package transport

import (
	"errors"
	"fmt"
	"strings"
)

type MessageRequest struct {
	Command      string
	Package      string
	Dependencies []string
}
type MessageResponseCode string

const (
	IndexCmd  = "INDEX"
	QueryCmd  = "QUERY"
	RemoveCmd = "REMOVE"

	ERROR = "ERROR"
	FAIL  = "FAIL"
	OK    = "OK"

//  "((QUERY|INDEX|REMOVE)\|+[a-z]\|*)\w+"
)

func ValidateMessage(rawMsg string) error {
	msgSlice := strings.Split(rawMsg, "|")
	if len(msgSlice) > 2 && len(msgSlice) < 4 {
		if msgSlice[0] != IndexCmd && msgSlice[0] != RemoveCmd && msgSlice[0] != QueryCmd {
			return errors.New(fmt.Sprintf("Wrong command for the message '<command>|<package>|<dependencies>\\n' - '%s'", msgSlice[0]))
		}

		if msgSlice[1] == "" {
			return errors.New(fmt.Sprintf("Wrong package for the message '<command>|<package>|<dependencies>\\n' - '<package>' is required %s", msgSlice[1]))
		}

		dependenciesSlice := strings.Split(msgSlice[2], ",")
		if len(dependenciesSlice) > 1 && dependenciesSlice[0] != "" {
			for _, pkg := range dependenciesSlice {
				if pkg == "" {
					return errors.New(fmt.Sprintf("Wrong message format dependencies '<command>|<package>|<dependencies>\\n' dependencies should be a comma separated list '%v'", dependenciesSlice))
				}
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Wrong message format '<command>|<package>|<dependencies>\\n' '%v'", msgSlice))
}

func ExtractMessage(rawMsg string) (*MessageRequest, error) {
	msgSlice := strings.Split(rawMsg, "|")
	msg := &MessageRequest{
		Command:      msgSlice[0],
		Package:      msgSlice[1],
		Dependencies: make([]string, 0),
	}
	if len(msgSlice) > 2 {
		dependenciesSlice := strings.Split(msgSlice[2], ",")
		if len(dependenciesSlice) > 0 && dependenciesSlice[0] != "" {
			msg.Dependencies = append(msg.Dependencies, dependenciesSlice...)
		}
	}

	return msg, nil
}

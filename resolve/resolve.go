package resolve

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func Lookup(typ, key string) (string, error) {
	switch strings.ToLower(typ) {
	case "env":
		return os.Getenv(key), nil
	case "file":
		data, err := ioutil.ReadFile(key)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}
	return "", errors.New("unsupported lookup type: " + typ)
}

var reTags = regexp.MustCompile(`(?i)\{\s*([a-z][a-z0-9]*):([^}]*)\s*\}`)

func Value(input string) (string, error) {
	matches := reTags.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if len(match) == 3 {
			typ := match[1]
			key := match[2]

			replaceWith, err := Lookup(typ, key)
			if err != nil {
				return "", err
			} else if replaceWith == "" {
				return "", errors.New(typ + "-lookup returned empty string")
			}
			input = strings.Replace(input, match[0], replaceWith, -1)
		} else {
			return "", errors.New("invalid match")
		}
	}

	return input, nil
}

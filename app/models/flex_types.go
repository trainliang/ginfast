package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type FlexUint uint

func (v *FlexUint) UnmarshalJSON(data []byte) error {
	parsed, err := parseUintToken(data)
	if err != nil {
		return err
	}
	*v = FlexUint(parsed)
	return nil
}

func (v *FlexUint) UnmarshalText(text []byte) error {
	parsed, err := parseUintString(string(text))
	if err != nil {
		return err
	}
	*v = FlexUint(parsed)
	return nil
}

type FlexInt int

func (v *FlexInt) UnmarshalJSON(data []byte) error {
	parsed, err := parseIntToken(data)
	if err != nil {
		return err
	}
	*v = FlexInt(parsed)
	return nil
}

func (v *FlexInt) UnmarshalText(text []byte) error {
	parsed, err := parseIntString(string(text))
	if err != nil {
		return err
	}
	*v = FlexInt(parsed)
	return nil
}

type FlexInt8 int8

func (v *FlexInt8) UnmarshalJSON(data []byte) error {
	parsed, err := parseInt8Token(data)
	if err != nil {
		return err
	}
	*v = FlexInt8(parsed)
	return nil
}

func (v *FlexInt8) UnmarshalText(text []byte) error {
	parsed, err := parseInt8String(string(text))
	if err != nil {
		return err
	}
	*v = FlexInt8(parsed)
	return nil
}

type FlexString string

func (v *FlexString) UnmarshalJSON(data []byte) error {
	token, err := parseRawToken(data)
	if err != nil {
		return err
	}
	*v = FlexString(token)
	return nil
}

func (v *FlexString) UnmarshalText(text []byte) error {
	*v = FlexString(strings.TrimSpace(string(text)))
	return nil
}

func StringToUint(value string) uint {
	parsed, err := parseUintString(value)
	if err != nil {
		return 0
	}
	return parsed
}

func parseRawToken(data []byte) (string, error) {
	token := strings.TrimSpace(string(data))
	if token == "" || token == "null" {
		return "", nil
	}

	if strings.HasPrefix(token, "\"") {
		var decoded string
		if err := json.Unmarshal(data, &decoded); err != nil {
			return "", err
		}
		return strings.TrimSpace(decoded), nil
	}

	return token, nil
}

func parseUintToken(data []byte) (uint, error) {
	return parseUintStringFromToken(data)
}

func parseIntToken(data []byte) (int, error) {
	token, err := parseRawToken(data)
	if err != nil {
		return 0, err
	}
	return parseIntString(token)
}

func parseInt8Token(data []byte) (int8, error) {
	token, err := parseRawToken(data)
	if err != nil {
		return 0, err
	}
	return parseInt8String(token)
}

func parseUintStringFromToken(data []byte) (uint, error) {
	token, err := parseRawToken(data)
	if err != nil {
		return 0, err
	}
	return parseUintString(token)
}

func parseUintString(value string) (uint, error) {
	token := strings.TrimSpace(value)
	if token == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func parseIntString(value string) (int, error) {
	token := strings.TrimSpace(value)
	if token == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(parsed), nil
}

func parseInt8String(value string) (int8, error) {
	token := strings.TrimSpace(value)
	if token == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(token, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(parsed), nil
}

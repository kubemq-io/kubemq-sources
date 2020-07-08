package types

import (
	"fmt"
	"strconv"
)

type Metadata map[string]string

func NewMetadata() Metadata {
	return map[string]string{}
}
func UnmarshallMetadata(meta string) (Metadata, error) {
	if meta == "" {
		return NewMetadata(), nil
	}
	m := Metadata{}
	err := json.UnmarshalFromString(meta, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (m Metadata) String() string {
	str, err := json.MarshalToString(&m)
	if err != nil {
		return ""
	}
	return str
}
func (m Metadata) Set(key, value string) Metadata {
	m[key] = value
	return m
}
func (m Metadata) Get(key string) string {
	return m[key]
}

func (m Metadata) ParseString(key, defaultValue string) string {

	if val, ok := m[key]; ok && val != "" {
		return val
	} else {
		return defaultValue
	}
}

func (m Metadata) ParseStringMap(key string, stringMap map[string]string) (string, error) {
	if val, ok := stringMap[m[key]]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("no valid key found")
	}
}

func (m Metadata) MustParseString(key string) (string, error) {
	if val, ok := m[key]; ok && val != "" {
		return val, nil
	} else {
		return "", fmt.Errorf("value of key %s cannot be empty", key)
	}
}

func (m Metadata) MustParseInt(key string) (int, error) {

	if val, ok := m[key]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid conversion error for value %s", val)
		}
		return int(parsedVal), nil
	} else {
		return 0, fmt.Errorf("key %s not foud for int coneversion", val)
	}
}

func (m Metadata) ParseInt(key string, defaultValue int) int {
	if val, ok := m[key]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return defaultValue
		} else {
			return int(parsedVal)
		}
	} else {
		return defaultValue
	}
}
func (m Metadata) ParseIntWithRange(key string, defaultValue, min, max int) (int, error) {
	val := m.ParseInt(key, defaultValue)
	if val < min {
		return 0, fmt.Errorf("conversion value cannot be lower than %d", min)
	}
	if val > max {
		return 0, fmt.Errorf("conversion value cannot be higher than %d", min)
	}
	return val, nil
}

func (m Metadata) MustParseIntWithRange(key string, min, max int) (int, error) {
	val, err := m.MustParseInt(key)
	if err != nil {
		return 0, err
	}
	if val < min {
		return 0, fmt.Errorf("conversion value cannot be lower than %d", min)
	}
	if val > max {
		return 0, fmt.Errorf("conversion value cannot be higher than %d", min)
	}
	return val, nil
}

func (m Metadata) ParseBool(key string, defaultValue bool) bool {
	if val, ok := m[key]; ok && val != "" {
		parsedVal, err := strconv.ParseBool(val)
		if err != nil {
			return defaultValue
		} else {
			return parsedVal
		}
	} else {
		return defaultValue
	}
}
func (m Metadata) MustParseBool(key string) (bool, error) {

	if val, ok := m[key]; ok && val != "" {
		parsedVal, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("invalid bool conversion error for value %s", val)
		}
		return parsedVal, nil
	} else {
		return false, fmt.Errorf("key %s not foud for bool coneversion", val)
	}
}
func (m Metadata) MustParseJsonMap(key string) (map[string]string, error) {
	if val, ok := m[key]; ok && val != "" {
		if val == "" {
			return map[string]string{}, nil
		}
		parsedVal := make(map[string]string)
		err := json.Unmarshal([]byte(val), &parsedVal)
		if err != nil {
			return nil, fmt.Errorf("invalid json conversion to map[string]string %s", val)
		}
		return parsedVal, nil
	} else {
		return map[string]string{}, nil
	}
}

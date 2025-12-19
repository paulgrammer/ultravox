package ultravox

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// UltravoxDuration is a wrapper around time.Duration that marshals to seconds
type UltravoxDuration time.Duration

// String returns the duration as a string in standard Go duration format
func (d UltravoxDuration) String() string {
	return time.Duration(d).String()
}

// formatDuration is a helper that formats the duration as a string in seconds
func (d UltravoxDuration) formatDuration() string {
	seconds := time.Duration(d).Seconds()
	// Format with no decimal places if it's a whole number
	if seconds == float64(int64(seconds)) {
		return fmt.Sprintf("%.0fs", seconds)
	}
	return fmt.Sprintf("%gs", seconds)
}

// MarshalJSON converts the duration to a string in seconds like "60s"
func (d UltravoxDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.formatDuration())
}

// parseDuration is a helper that parses a duration from a string
// supporting multiple formats:
// - Duration strings ("30s", "1m30s", etc.)
// - Numeric strings ("30") as seconds
func parseDuration(s string) (UltravoxDuration, error) {
	// Try parsing as duration string first ("30s", "1m", etc.)
	if parsed, err := time.ParseDuration(s); err == nil {
		return UltravoxDuration(parsed), nil
	}

	// Try parsing as numeric string ("30")
	if seconds, err := strconv.ParseFloat(s, 64); err == nil {
		return UltravoxDuration(time.Duration(seconds * float64(time.Second))), nil
	}

	return 0, fmt.Errorf("invalid duration format: %q", s)
}

// UnmarshalJSON converts JSON data to duration supporting multiple formats:
// - Numbers (30) as seconds
// - Numeric strings ("30") as seconds
// - Duration strings ("30s", "1m30s", etc.)
func (d *UltravoxDuration) UnmarshalJSON(data []byte) error {
	var rawValue interface{}
	if err := json.Unmarshal(data, &rawValue); err != nil {
		return err
	}

	switch v := rawValue.(type) {
	case float64:
		// Direct number (30)
		*d = UltravoxDuration(time.Duration(v * float64(time.Second)))
		return nil

	case string:
		parsed, err := parseDuration(v)
		if err != nil {
			return err
		}
		*d = parsed
		return nil

	default:
		return fmt.Errorf("duration must be a number or string, got %T", rawValue)
	}
}

// MarshalYAML converts the duration to a string in seconds like "60s"
func (d UltravoxDuration) MarshalYAML() (interface{}, error) {
	return d.formatDuration(), nil
}

// UnmarshalYAML converts YAML data to duration supporting multiple formats:
// - Numbers (30) as seconds
// - Numeric strings ("30") as seconds
// - Duration strings ("30s", "1m30s", etc.)
func (d *UltravoxDuration) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("duration must be a scalar value, got %v", value.Kind)
	}

	parsed, err := parseDuration(value.Value)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

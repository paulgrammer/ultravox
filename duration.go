package ultravox

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// UltravoxDuration is a wrapper around time.Duration that marshals to seconds
type UltravoxDuration time.Duration

// String returns the duration as a string in standard Go duration format
func (d UltravoxDuration) String() string {
	return time.Duration(d).String()
}

// MarshalJSON converts the duration to a string in seconds like "60s"
func (d UltravoxDuration) MarshalJSON() ([]byte, error) {
	seconds := time.Duration(d).Seconds()
	// Format with no decimal places if it's a whole number
	if seconds == float64(int64(seconds)) {
		return json.Marshal(fmt.Sprintf("%.0fs", seconds))
	}
	return json.Marshal(fmt.Sprintf("%gs", seconds))
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
		// Try parsing as duration string first ("30s", "1m", etc.)
		if parsed, err := time.ParseDuration(v); err == nil {
			*d = UltravoxDuration(parsed)
			return nil
		}

		// Try parsing as numeric string ("30")
		if seconds, err := strconv.ParseFloat(v, 64); err == nil {
			*d = UltravoxDuration(time.Duration(seconds * float64(time.Second)))
			return nil
		}

		return fmt.Errorf("invalid duration format: %q", v)

	default:
		return fmt.Errorf("duration must be a number or string, got %T", rawValue)
	}
}

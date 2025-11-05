package service

import (
	"fmt"
	"time"
)

// GenerateID creates a timestamp-based unique ID in format YYYYMMDD-HHMMSS-nnnnn
// Example: 20251104-143000-00001
func GenerateID() string {
	now := time.Now()
	dateTime := now.Format("20060102-150405")
	micro := now.Nanosecond() / 1000 % 100000
	return fmt.Sprintf("%s-%05d", dateTime, micro)
}

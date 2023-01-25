package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestNext(t *testing.T) {

	tests := []struct {
		now        string
		expression string
		expected   string
		//want       *CronExpression
	}{
		// Simple cases
		{"2012-07-09 14:45", "0 0/15 * * * *", "2012-07-09 14:45"},
		{"2012-07-09 14:46", "0 0/15 * * * *", "2012-07-09 15:00"},
		{"2012-07-09 14:59", "0 0/15 * * * *", "2012-07-09 15:00"},
		{"2012-07-09 14:59:59", "0 0/15 * * * *", "2012-07-09 15:00"},

		// shift hours
		{"2012-07-09 15:45", "0 20-35/15 * * * *", "2012-07-09 16:20"},

		// shift days
		{"2012-07-09 23:46", "0 */15 * * * *", "2012-07-10 00:00"},
		{"2012-07-09 23:45", "0 20-35/15 * * * *", "2012-07-10 00:20"},
		{"2012-07-09 23:35:51", "15/35 20-35/15 * * * *", "2012-07-10 00:20:15"},
		{"2012-07-09 23:35:51", "15/35 20-35/15 1/2 * * *", "2012-07-10 01:20:15"},
		{"2012-07-09 23:35:51", "15/35 20-35/15 10-12 * * *", "2012-07-10 10:20:15"},

		{"2012-07-09 23:35:51", "15/35 20-35/15 1/2 */2 * *", "2012-07-11 01:20:15"},
		{"2012-07-09 23:35:51", "15/35 20-35/15 * 9-20 * *", "2012-07-10 00:20:15"},
		{"2012-07-09 23:35:51", "15/35 20-35/15 * 9-20 Jul *", "2012-07-10 00:20:15"},

		// Wrap around months
		{"2012-07-09 23:35", "0 0 0 9 Apr-Oct ?", "2012-08-09 00:00"},
		{"2012-07-09 23:35", "0 0 0 */5 Apr,Aug,Oct Mon", "2012-08-06 00:00"},
		{"2012-07-09 23:35", "0 0 0 */5 Oct Mon", "2012-10-01 00:00"},

		// Wrap around years
		{"2012-07-09 23:35", "0 0 0 * Feb Mon", "2013-02-04 00:00"},
		{"2012-07-09 23:35", "0 0 0 * Feb Mon/2", "2013-02-01 00:00"},

		// Wrap around minute, hour, day, month, and year
		{"2012-12-31 23:59:45", "0 * * * * *", "2013-01-01 00:00:00"},

		// Leap year
		{"2012-07-09 23:35", "0 0 0 29 Feb ?", "2016-02-29 00:00"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			cr, err := Parse(tt.expression)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			actual := cr.Next(parseTime(tt.now))
			expected := parseTime(tt.expected)
			if actual != expected {
				t.Errorf("Expected: %v Actual: %v", expected, actual)
			}
		})
	}
}

func parseTime(val string) time.Time {
	parsed, err := time.Parse("2006-01-02 15:04", val)
	if err != nil {
		parsed, _ = time.Parse("2006-01-02 15:04:05", val)
	}
	return parsed
}

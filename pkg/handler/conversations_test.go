package handler

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseFlexibleDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDate string
		wantErr  bool
	}{
		// Standard formats (existing)
		{
			name:     "YYYY-MM-DD",
			input:    "2025-07-15",
			wantDate: "2025-07-15",
			wantErr:  false,
		},
		{
			name:     "YYYY/MM/DD",
			input:    "2025/07/15",
			wantDate: "2025-07-15",
			wantErr:  false,
		},

		// New flexible month-year formats
		{
			name:     "Month Year - July 2025",
			input:    "July 2025",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Year Month - 2025 July",
			input:    "2025 July",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Abbreviated Month Year - Jul 2025",
			input:    "Jul 2025",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Year Abbreviated Month - 2025 Jul",
			input:    "2025 Jul",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Case insensitive - july 2025",
			input:    "july 2025",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Case insensitive - JULY 2025",
			input:    "JULY 2025",
			wantDate: "2025-07-01",
			wantErr:  false,
		},

		// Day-Month-Year formats
		{
			name:     "1-July-2025",
			input:    "1-July-2025",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "July-25-2025",
			input:    "July-25-2025",
			wantDate: "2025-07-25",
			wantErr:  false,
		},
		{
			name:     "July 10 2025",
			input:    "July 10 2025",
			wantDate: "2025-07-10",
			wantErr:  false,
		},
		{
			name:     "10 July 2025",
			input:    "10 July 2025",
			wantDate: "2025-07-10",
			wantErr:  false,
		},
		{
			name:     "31-December-2025",
			input:    "31-December-2025",
			wantDate: "2025-12-31",
			wantErr:  false,
		},
		{
			name:     "2025 July 10",
			input:    "2025 July 10",
			wantDate: "2025-07-10",
			wantErr:  false,
		},

		// Various month names
		{
			name:     "January full name",
			input:    "January 2025",
			wantDate: "2025-01-01",
			wantErr:  false,
		},
		{
			name:     "February abbreviated",
			input:    "Feb 2025",
			wantDate: "2025-02-01",
			wantErr:  false,
		},
		{
			name:     "September with Sept abbreviation",
			input:    "Sept 2025",
			wantDate: "2025-09-01",
			wantErr:  false,
		},

		// Relative dates
		{
			name:     "today",
			input:    "today",
			wantDate: time.Now().UTC().Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "yesterday",
			input:    "yesterday",
			wantDate: time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "Today with capital T",
			input:    "Today",
			wantDate: time.Now().UTC().Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "Yesterday with capital Y",
			input:    "Yesterday",
			wantDate: time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "TODAY all caps",
			input:    "TODAY",
			wantDate: time.Now().UTC().Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "YESTERDAY all caps",
			input:    "YESTERDAY",
			wantDate: time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "tomorrow",
			input:    "tomorrow",
			wantDate: time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "5 days ago",
			input:    "5 days ago",
			wantDate: time.Now().UTC().AddDate(0, 0, -5).Format("2006-01-02"),
			wantErr:  false,
		},
		{
			name:     "1 day ago",
			input:    "1 day ago",
			wantDate: time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02"),
			wantErr:  false,
		},

		// Edge cases
		{
			name:     "Whitespace trimming",
			input:    "  July 2025  ",
			wantDate: "2025-07-01",
			wantErr:  false,
		},
		{
			name:     "Invalid month name",
			input:    "Jully 2025",
			wantDate: "",
			wantErr:  true,
		},
		{
			name:     "Invalid date format",
			input:    "2025-13-01",
			wantDate: "",
			wantErr:  true,
		},
		{
			name:     "Invalid day for month",
			input:    "31-February-2025",
			wantDate: "",
			wantErr:  true,
		},
		{
			name:     "Empty string",
			input:    "",
			wantDate: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotDate, err := parseFlexibleDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlexibleDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotDate != tt.wantDate {
				t.Errorf("parseFlexibleDate() gotDate = %v, want %v", gotDate, tt.wantDate)
			}
		})
	}
}

func TestBuildDateFilters(t *testing.T) {
	tests := []struct {
		name    string
		before  string
		after   string
		on      string
		during  string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "On with flexible format July 2025",
			before:  "",
			after:   "",
			on:      "July 2025",
			during:  "",
			want:    map[string]string{"on": "2025-07-01"},
			wantErr: false,
		},
		{
			name:    "Before and After with flexible formats",
			before:  "December 2025",
			after:   "January 2025",
			on:      "",
			during:  "",
			want:    map[string]string{"before": "2025-12-01", "after": "2025-01-01"},
			wantErr: false,
		},
		{
			name:    "During with day format",
			before:  "",
			after:   "",
			on:      "",
			during:  "15-July-2025",
			want:    map[string]string{"during": "2025-07-15"},
			wantErr: false,
		},
		{
			name:    "Error: on with other filters",
			before:  "2025-12-01",
			after:   "",
			on:      "July 2025",
			during:  "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Error: during with before",
			before:  "2025-12-01",
			after:   "",
			on:      "",
			during:  "July 2025",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Error: after date is after before date",
			before:  "January 2025",
			after:   "December 2025",
			on:      "",
			during:  "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid: complex date formats",
			before:  "31-December-2025",
			after:   "1-January-2025",
			on:      "",
			during:  "",
			want:    map[string]string{"before": "2025-12-31", "after": "2025-01-01"},
			wantErr: false,
		},
		{
			name:    "Error: invalid date format",
			before:  "",
			after:   "",
			on:      "Jully 2025",
			during:  "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDateFilters(tt.before, tt.after, tt.on, tt.during)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildDateFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("buildDateFilters() got map length = %v, want %v", len(got), len(tt.want))
					return
				}
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("buildDateFilters() got[%s] = %v, want %v", k, got[k], v)
					}
				}
			}
		})
	}
}

func TestLimitByExpression_Valid(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		minSecs int64 // inclusive
		maxSecs int64 // exclusive
	}{
		{"1 day", "1d", 0, 86400},
		{"2 days", "2d", 86400, 172800},
		{"1 week", "1w", 6 * 86400, 7 * 86400},
		{"2 weeks", "2w", 13 * 86400, 14 * 86400},
		{"1 month", "1m", 28 * 86400, 31 * 86400},
		{"2 months", "2m", 2 * 28 * 86400, 2 * 31 * 86400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slackLimit, oldestStr, latestStr, err := limitByExpression(tt.input, defaultConversationsExpressionLimit)
			if err != nil {
				t.Fatalf("expected no error for %q, got %v", tt.input, err)
			}
			if slackLimit != 100 {
				t.Errorf("expected slackLimit=100 for %q, got %d", tt.input, slackLimit)
			}

			// Parse the "1234567890.000000" format back to an integer
			o, err := strconv.ParseInt(strings.TrimSuffix(oldestStr, ".000000"), 10, 64)
			if err != nil {
				t.Fatalf("invalid oldest timestamp %q: %v", oldestStr, err)
			}
			l, err := strconv.ParseInt(strings.TrimSuffix(latestStr, ".000000"), 10, 64)
			if err != nil {
				t.Fatalf("invalid latest timestamp %q: %v", latestStr, err)
			}

			if l <= o {
				t.Errorf("for %q expected latest(%d) > oldest(%d)", tt.input, l, o)
			}
			diff := l - o
			if diff < tt.minSecs || diff >= tt.maxSecs {
				t.Errorf(
					"for %q expected span in [%d, %d), got %d",
					tt.input, tt.minSecs, tt.maxSecs, diff,
				)
			}
		})
	}
}

func TestLimitByExpression_Invalid(t *testing.T) {
	invalid := []string{
		"d",   // too short
		"0d",  // zero
		"-1d", // negative
		"1x",  // bad suffix
		"1",   // missing suffix
		"01",  // no suffix + zero value
	}

	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			_, _, _, err := limitByExpression(input, defaultConversationsExpressionLimit)
			if err == nil {
				t.Errorf("expected error for %q, got nil", input)
			}
		})
	}
}

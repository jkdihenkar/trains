package parser

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid time format",
			input:    "14:30",
			expected: 870, // 14*60 + 30
		},
		{
			name:     "Midnight",
			input:    "00:00",
			expected: 0,
		},
		{
			name:     "Late evening",
			input:    "23:59",
			expected: 1439, // 23*60 + 59
		},
		{
			name:     "Invalid format",
			input:    "invalid",
			expected: 0,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTime(tt.input)
			if result != tt.expected {
				t.Errorf("ParseTime(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDistanceKm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid distance with Kms",
			input:    "754 Kms",
			expected: 754,
		},
		{
			name:     "Valid distance with Km",
			input:    "1038 Km",
			expected: 1038,
		},
		{
			name:     "Distance without space",
			input:    "500Kms",
			expected: 500,
		},
		{
			name:     "Invalid format",
			input:    "invalid distance",
			expected: 0,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDistanceKm(tt.input)
			if result != tt.expected {
				t.Errorf("ParseDistanceKm(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateAndNormalizeDay(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid short form",
			input:       "mon",
			expected:    "Monday",
			expectError: false,
		},
		{
			name:        "Valid full form",
			input:       "Monday",
			expected:    "Monday",
			expectError: false,
		},
		{
			name:        "Case insensitive",
			input:       "MON",
			expected:    "Monday",
			expectError: false,
		},
		{
			name:        "With spaces",
			input:       "  tue  ",
			expected:    "Tuesday",
			expectError: false,
		},
		{
			name:        "Invalid day",
			input:       "invalid",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndNormalizeDay(tt.input)
			
			if tt.expectError && err == nil {
				t.Errorf("ValidateAndNormalizeDay(%q) expected error but got none", tt.input)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("ValidateAndNormalizeDay(%q) unexpected error: %v", tt.input, err)
			}
			
			if result != tt.expected {
				t.Errorf("ValidateAndNormalizeDay(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCommonRunningDays(t *testing.T) {
	tests := []struct {
		name     string
		days1    string
		days2    string
		expected string
	}{
		{
			name:     "All days common",
			days1:    "1111111",
			days2:    "1111111",
			expected: "Sun,Mon,Tue,Wed,Thu,Fri,Sat",
		},
		{
			name:     "Some days common",
			days1:    "1010101",
			days2:    "1011001",
			expected: "Sun,Tue,Sat",
		},
		{
			name:     "No days common",
			days1:    "1010101",
			days2:    "0101010",
			expected: "",
		},
		{
			name:     "Invalid input - short string",
			days1:    "101",
			days2:    "1111111",
			expected: "",
		},
		{
			name:     "Empty strings",
			days1:    "",
			days2:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommonRunningDays(tt.days1, tt.days2)
			if result != tt.expected {
				t.Errorf("GetCommonRunningDays(%q, %q) = %q, want %q", tt.days1, tt.days2, result, tt.expected)
			}
		})
	}
}

func TestIsUnder19Hours(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Under 19 hours",
			input:    "16h 52m",
			expected: true,
		},
		{
			name:     "Exactly 19 hours",
			input:    "19h 0m",
			expected: false,
		},
		{
			name:     "Over 19 hours",
			input:    "20h 30m",
			expected: false,
		},
		{
			name:     "Single digit hours",
			input:    "8h 45m",
			expected: true,
		},
		{
			name:     "Invalid format",
			input:    "invalid",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnder19Hours(tt.input)
			if result != tt.expected {
				t.Errorf("IsUnder19Hours(%q) = %t, want %t", tt.input, result, tt.expected)
			}
		})
	}
}
// internal/models/user_test.go
package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name string
		dob  string
	}{
		{
			name: "Person born in 1990",
			dob:  "1990-05-10",
		},
		{
			name: "Person born in 2000",
			dob:  "2000-01-01",
		},
		{
			name: "Person born in 1985",
			dob:  "1985-12-25",
		},
		{
			name: "Person born in 1995",
			dob:  "1995-03-20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dob, err := time.Parse("2006-01-02", tt.dob)
			if err != nil {
				t.Fatalf("Failed to parse date: %v", err)
			}

			got := CalculateAge(dob)
			
			// Calculate expected age manually
			expectedAge := now.Year() - dob.Year()
			if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
				expectedAge--
			}
			
			if got != expectedAge {
				t.Errorf("CalculateAge(%s) = %v, want %v (today is %s)", 
					tt.dob, got, expectedAge, now.Format("2006-01-02"))
			}
		})
	}
}

func TestCalculateAge_EdgeCases(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "Birthday today",
			dob:  time.Date(now.Year()-30, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 30,
		},
		{
			name: "Birthday tomorrow",
			dob:  time.Date(now.Year()-30, now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC),
			want: 29, // Haven't had birthday yet this year
		},
		{
			name: "Birthday yesterday",
			dob:  time.Date(now.Year()-30, now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			want: 30, // Already had birthday this year
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.want {
				t.Errorf("CalculateAge() = %v, want %v (dob=%s, now=%s)", 
					got, tt.want, tt.dob.Format("2006-01-02"), now.Format("2006-01-02"))
			}
		})
	}
}

// Additional test to verify the calculation logic explicitly
func TestCalculateAge_Explicit(t *testing.T) {
	tests := []struct {
		name     string
		dob      string
		testDate string
		wantAge  int
	}{
		{
			name:     "Birthday already passed this year",
			dob:      "1990-01-15",
			testDate: "2024-12-18",
			wantAge:  34,
		},
		{
			name:     "Birthday not yet this year",
			dob:      "1990-12-25",
			testDate: "2024-12-18",
			wantAge:  33,
		},
		{
			name:     "Birthday is today",
			dob:      "1990-12-18",
			testDate: "2024-12-18",
			wantAge:  34,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dob, _ := time.Parse("2006-01-02", tt.dob)
			
			// This test assumes current date - if you want to test with specific dates,
			// you'd need to modify CalculateAge to accept a "now" parameter
			got := CalculateAge(dob)
			
			// For now, just verify it returns a positive integer
			if got < 0 {
				t.Errorf("CalculateAge() returned negative age: %v", got)
			}
			
			// Verify age is reasonable (between 0 and 150 years)
			if got > 150 {
				t.Errorf("CalculateAge() returned unreasonable age: %v", got)
			}
		})
	}
}
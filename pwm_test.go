package pwm

import (
	"testing"
)

var pwmA = Pwm{
	0: {"A": 1, "C": 0, "G": 0, "T": 0},
}
var pwmEq2 = Pwm{
	0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
	1: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
}
var pwmC2 = Pwm{
	0: {"A": 0.1, "C": 0.7, "G": 0.1, "T": 0.1},
	1: {"A": 0.1, "C": 0.7, "G": 0.1, "T": 0.1},
}
var inPwmA = Pwm{
	0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
}

// Test for panics based on: https://stackoverflow.com/questions/31595791/how-to-test-panics
// See kmers_test.go in rbns for how to test returned values if there is no panic
func TestValidatePWM(t *testing.T) {
	tables := []struct {
		name 	string // name of the test
		pwm Pwm // pwm to test
		wantPanic bool // whether the test should trigger a panic
	}{
		{"valid pwmA",
			pwmA,
			false,
		},
		{"valid pwmEq2",
			pwmEq2,
			false,
		},
		{"valid pwmC2",
			pwmC2,
			false,
		},
		// Invalid PWMs, expect a Panic
		{"first position is 1, not 0",
			Pwm{
				1: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
			},
			true,
		},
		{"two positions are 0 and 2",
			Pwm{
				0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
				2: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
			},
			true,
		},
		{"C missing",
			Pwm{
				0: {"A": 0.25, "G": 0.5, "T": 0.25},
			},
			true,
		},
		{"position 1 probs add up to 0.75",
			Pwm{
				0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
				1: {"A": 0., "C": 0.25, "G": 0.25, "T": 0.25},
			},
			true,
		},
	}

	for _, tt := range tables {
		// Test if there's a panic
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Validate() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()
			// Test for panic
			tt.pwm.Validate()
		})
	}
}

func TestValidateNtSeq(t *testing.T) {
	tables := []struct {
		name 	string // name of the test
		seq ntSeq // sequence to test
		wantPanic bool // whether the test should trigger a panic
	}{
		{"valid seq",
			"ACGAAACTTAA",
			false,
		},
		{"lower case nt",
			"aCGAAACTTAA",
			true,
		},
		{"Us instead of Ts",
			"ACGAAACUUAA",
			true,
		},
		{"Ns",
			"ACGAANCTTAA",
			true,
		},
	}

	for _, tt := range tables {
		// Test if there's a panic
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Validate() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()
			tt.seq.Validate()
		})
	}
}
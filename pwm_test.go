package pwm

import (
	"testing"
)

var pwmA = Pwm{
	0: {"A": 1, "C": 0, "G": 0, "T": 0},
}
var pwmEq = Pwm{
	0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
}
var inPwmA = Pwm{
	0: {"A": 0.25, "C": 0.25, "G": 0.25, "T": 0.25},
}

// Test for panics based on: https://stackoverflow.com/questions/31595791/how-to-test-panics
func TestValidatePWM(t *testing.T) {
	tables := []struct {
		name 	string // name of the test
		pwm Pwm // name of the test
		wantPanic bool // whether the test should trigger a panic (i.e., k < 1)
	}{
		// a valid PWM
		{"valid pwmA",
			pwmA,
			false,
		},
		// Invalid PWMs
		// - positions not 0, 1, 2, ....
		{"first position is 1",
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
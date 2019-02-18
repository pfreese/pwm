package pwm

import (
	"reflect"
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

func TestValidatePWM(t *testing.T) {
	tables := []struct {
		name 		string 	// name of the test
		pwm 		Pwm 	// pwm to test
		wantPanic 	bool 	// whether the test should trigger a panic
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
		name 		string 	// name of the test
		seq 		ntSeq 	// sequence to test
		wantPanic	bool 	// whether the test should trigger a panic
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

func TestAddPseudocount(t *testing.T) {
	tables := []struct {
		name      string // name of the test
		pwm       Pwm    // pwm to test
		pseudo	float64    // pseudcount to add to each base
		expPwm 	Pwm   // expected PWM after pseudocount
		wantPanic 	bool   // whether the test should trigger a panic
	}{
		{
			"add 1 pseudocount to each nt",
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
			},
			1,
			Pwm{
				0: {"A": 2./5, "C": 1./5, "G": 1./5, "T": 1./5},
			},
			false,
		},
		{
			"add 1 pseudocount to all positions",
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
				1: {"A": 0.5, "C": 0.5, "G": 0, "T": 0},
			},
			1,
			Pwm{
				0: {"A": 2./5, "C": 1./5, "G": 1./5, "T": 1./5},
				1: {"A": 1.5/5, "C": 1.5/5, "G": 1./5, "T": 1./5},
			},
			false,
		},
		{
			"add 0 pseudocount leaves original PWM",
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
			},
			0,
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
			},
			false,
		},
		{
			"add 0.25 pseudocount",
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
			},
			0.25,
			Pwm{
				0: {"A": 1.25/2, "C": 0.25/2, "G": 0.25/2, "T": 0.25/2},
			},
			false,
		},
		{
			"panic if pseudocount <0",
			Pwm{
				0: {"A": 1, "C": 0, "G": 0, "T": 0},
			},
			-0.25,
			Pwm{},
			true,
		},
	}
	for _, tt := range tables {
		// Test if there's a panic
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("addPseudocount() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()
			(&tt.pwm).addPseudocount(tt.pseudo)
			if !reflect.DeepEqual(tt.pwm, tt.expPwm) {
				t.Errorf("Failed to add pseudocount: expected %v, got %v", tt.expPwm, tt.pwm)
			}
		})
	}
}
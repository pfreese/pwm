package pwm

import (
	"fmt"
	"math"
)

type Nt string
type PosProb map[Nt]float64
type Pwm map[int]PosProb

var nts = [4]Nt{"A", "C", "G", "T"}

func (c Nt) String() string {
	return fmt.Sprintf("%v", string(c))
}

func (p Pwm) Validate() string {
	for i := 0; i < len(p); i++ {
		// panic if position is not in pwm
		if _, ok := p[i]; !ok {
			panic(fmt.Sprintf("position %d not in pwm - must be indexed consecutively from 0", i))
		}
		// panic if "A", "C", "G", and "T" not in pwm, or if their
		// probabilities don't sum to "1"
		var prob float64
		for _, nt := range nts {
			if ntProb, ok := p[i][nt]; ok {
				prob += ntProb
			} else {
				panic(fmt.Sprintf("nt %s not in pwm", nt))
			}

		}
		if math.Abs(prob - 1) > 0.0001 {
			panic(fmt.Sprintf("pos. %d prob sums to %f", i, prob))
		}
	}
	return "ok"
}

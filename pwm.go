package pwm

import (
	"fmt"
	"math"
)

const minCount = 0.001

type ntSeq string
type Nt string
type PosProb map[Nt]float64
type Pwm map[int]PosProb

var nts = [4]Nt{"A", "C", "G", "T"}

func (c Nt) String() string {
	return fmt.Sprintf("%v", string(c))
}

func (s ntSeq) String() string {
	return fmt.Sprintf("%v", string(s))
}

// Ensure that a PWM (1) has positions indexed consecutively starting from 0,
// (2) each position has "A", "C", "G", and "T" entries, and the 4 probabilities sum to 1
func (p Pwm) Validate() {
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
}

// Ensure that a ntSeq consists solely of "A", "C", "G", and "T"
func (s ntSeq) Validate() {
	for i, c := range s {
		isNt := false
		for _, nt := range nts {
			if Nt(string(c)) == nt {
				isNt = true
			}
		}
		if !isNt {
			panic(fmt.Sprintf("position %d (=%v) not a valid nt ('A'/'C'/'G'/'T'", i, string(c)))
		}
	}
}

func (p *Pwm) addPseudocount(pseudo float64) {
	p.Validate()
	if pseudo < 0 {
		panic(fmt.Sprintf("pseudo %f is less than 0", pseudo))
	}
	var countsWithPsuedo = 1 + 4*pseudo
	//var pPseudo Pwm
	pPseudo := make(Pwm)
	for i, probs := range *p {
		var pseudoProbs = PosProb{
			"A": (probs["A"] + pseudo) / countsWithPsuedo,
			"C": (probs["C"] + pseudo) / countsWithPsuedo,
			"G": (probs["G"] + pseudo) / countsWithPsuedo,
			"T": (probs["T"] + pseudo) / countsWithPsuedo,
		}
		pPseudo[i] = pseudoProbs
	}
	*p = pPseudo
}

// Add a pseudocount if any entry is less than minCount
func (p *Pwm) addPseudoIfNecessary() {
	p.Validate()
	var needsPseudo = false
	for _, probs := range *p {
		for _, nt := range nts {
			if probs[nt] < minCount {
				needsPseudo = true
				break
			}
		}
	}
	if needsPseudo {
		p.addPseudocount(minCount)
	}
}

func (p *Pwm) scoreSeq(s ntSeq) (logProb float64) {
	if len(s) != len(*p) || len(s) == 0 {
		return math.Inf(-1)
	}
	for i, c := range s {
		prob := (*p)[i][Nt(string(c))]
		if prob == 0 {
			return math.Inf(-1)
		}
		logProb += math.Log10(prob)
	}
	return logProb
}

package spellcheck

import "errors"

var ErrCheckTimeout = errors.New("check request takes too long, checker not responding")

// Returns error if spellchecker service is not responding 
type SpellChecker interface {
        Check(text string) (CheckResponse, error)
}

type CheckResponse []MisspellInfo

type MisspellInfo map[string]interface{}

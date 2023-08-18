package spellcheck

// Returns error if spellchecker service is not responding 
type SpellChecker interface {
        Check(text string) (CheckResponse, error)
}

type CheckResponse []MisspellInfo

type MisspellInfo map[string]interface{}

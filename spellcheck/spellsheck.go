package spellcheck

type SpellChecker interface {
        Check(text string) (bool, CheckResponse, error)
}

type CheckResponse []MisspellInfo

type MisspellInfo map[string]interface{}

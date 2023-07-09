package models

type Composition struct {
	ID         string
	ComposerID string
	Name       string
	Version    string
	Mirrors    []CompositionMirror
}

func (c *Composition) Describe() string {
	return c.Name
}

type CompositionMirror struct {
	ID            string
	CompositionID string
	Link          string
	Version       string
}

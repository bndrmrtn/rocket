package tokenizer

// Models / Schemas

type ModelAnnotation struct {
	Annotation string
	Values     []string
}

type ModelConfig struct {
	Type        string
	Attributes  []string
	Annotations []ModelAnnotation
}

type Model map[string]ModelConfig

type Schema map[string]ModelConfig

// Hashing

type Hashing struct {
	Name      string
	Provider  string
	Arguments []string
}

// Query

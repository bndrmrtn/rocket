package tokenizer

type Generated struct {
	Models  map[string]Model `json:"models"`
	Hashing []Hashing        `json:"hashing"`
	Queries []Query          `json:"queries"`
}

func NewGenerated() *Generated {
	return &Generated{
		Models:  map[string]Model{},
		Hashing: []Hashing{},
		Queries: []Query{},
	}
}

// Models / Schemas

type ModelAnnotation struct {
	Annotation string   `json:"annotation"`
	Arguments  []string `json:"arguments,omitempty"`
}

type ModelRelation struct {
	Model string `json:"model"`
	Field string `json:"field"`
	Type  string `json:"type"`
}

type ModelConfig struct {
	Type        string            `json:"type"`
	Attributes  []string          `json:"attributes,omitempty"`
	Annotations []ModelAnnotation `json:"annotations,omitempty"`
	Relation    *ModelRelation    `json:"relation,omitempty"`
	Ignore      bool              `json:"ignore"`
}

type Model map[string]ModelConfig

// Query

type QueryArg struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Query struct {
	Name      string     `json:"name"`
	Arguments []QueryArg `json:"arguments,omitempty"`

	Tokens []string `json:"tokens,omitempty"`
}

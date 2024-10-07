package tokenizer

import "errors"

type Generated struct {
	Enums     map[string]map[string]string `json:"enums"`
	Models    map[string]Model             `json:"models"`
	ModelKeys map[string][]string          `json:"model_keys,omitempty"`
	Hashing   []Hashing                    `json:"hashing"`
	Queries   []Query                      `json:"queries"`
}

func NewGenerated() *Generated {
	return &Generated{
		Models:    map[string]Model{},
		Hashing:   []Hashing{},
		Queries:   []Query{},
		Enums:     make(map[string]map[string]string),
		ModelKeys: make(map[string][]string),
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
	Type        string           `json:"type"`
	Attributes  []string         `json:"attributes,omitempty"`
	Annotations ModelAnnotations `json:"annotations,omitempty"`
	Relation    *ModelRelation   `json:"relation,omitempty"`
	Ignore      bool             `json:"ignore"`
}

type Model map[string]ModelConfig

type ModelAnnotations []ModelAnnotation

func (a ModelAnnotations) Get(name string) (string, error) {
	for _, val := range a {
		if val.Annotation == name {
			return val.Arguments[0], nil
		}
	}

	return "", errors.New("annotation does not exists")
}

// Query

type QueryArg struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Query struct {
	Name      string     `json:"name"`
	Arguments []QueryArg `json:"arguments,omitempty"`

	Tokens []string `json:"tokens,omitempty"`

	BT BuildToken `json:"build_token,omitempty"`
}

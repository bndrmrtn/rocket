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
	Relation    ModelRelation     `json:"relation,omitempty"`
	Ignore      bool              `json:"ignore"`
}

type Model map[string]ModelConfig

// Query

type QueryArg struct {
	Name string
	Type string
}

type Query struct {
	Name      string
	Arguments []QueryArg

	Builder QueryBuilder
}

// QueryBuilder contains the actual steps for a query
type QueryBuilder struct {
	Method QueryBuilderMethod  // get, set, update, delete
	Get    *QueryBuilderGetter // nil if method is not get
}

type QueryBuilderMethod string

const (
	QueryMethodGet    QueryBuilderMethod = "get"
	QueryMethodSet    QueryBuilderMethod = "set"
	QueryMethodUpdate QueryBuilderMethod = "update"
	QueryMethodDelete QueryBuilderMethod = "delete"
)

type QuerySelectField struct {
	Model  string
	Fields map[string]string // field : type
}

type QueryMethod struct {
	Name  string
	Value string
}

type QueryBuilderGetter struct {
	MultiResult    bool
	Fields         []QuerySelectField
	Limit          int
	AppliedMethods []QueryMethod
}

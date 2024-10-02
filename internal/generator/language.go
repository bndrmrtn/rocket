package generator

import (
	"strings"
)

type Language struct {
	Name      string
	Extension string
}

func GetLanguage(name string) (*Language, error) {
	name = strings.ToLower(name)

	switch name {
	case "go", "golang":
		return &Language{
			Name:      "Go",
			Extension: "go",
		}, nil
	case "javascript", "js":
		return &Language{
			Name:      "JavaScript",
			Extension: "js",
		}, nil
	}

	return nil, unsupported("language not supported: " + name)
}

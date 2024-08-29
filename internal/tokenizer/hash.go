package tokenizer

// Hashing

type Hashing struct {
	Name      string   `json:"name"`
	Provider  string   `json:"provider"`
	Arguments []string `json:"arguments,omitempty"`
}

func getHashAlgoList() []string {
	return []string{"bcrypt", "sha256", "sha512", "md5", "sha1"}
}

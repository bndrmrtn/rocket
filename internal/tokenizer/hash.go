package tokenizer

// Hashing

type Hashing struct {
	Name      string
	Provider  string
	Arguments []string
}

func getHashAlgoList() []string {
	return []string{"bcrypt", "sha256", "sha512", "md5", "sha1"}
}

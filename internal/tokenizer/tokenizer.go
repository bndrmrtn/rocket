package tokenizer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Tokenizer struct {
	keywords    map[string]TokenType
	fileName    string
	src         string
	currentLine int
	currentPos  int
	tokens      []Token
	buf         string
}

func New(fileName string) (*Tokenizer, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return &Tokenizer{
		keywords:    keywords(),
		fileName:    fileName,
		src:         string(content) + "\n",
		currentLine: 1,
		currentPos:  0,
		tokens:      []Token{},
		buf:         "",
	}, nil
}

func (l *Tokenizer) shift() string {
	if l.buf == "" {
		l.buf = string(l.src[l.currentPos])
	}
	l.currentPos++
	current := string(l.src[l.currentPos])
	l.buf += current
	return current
}

func (l *Tokenizer) addToken(value string, t TokenType) {
	l.tokens = append(l.tokens, Token{
		Value:     value,
		TokenType: t,
		Line:      l.currentLine,
		TokenPos:  l.currentPos,
		FileName:  l.fileName,
	})
}

func (l *Tokenizer) currentPosOK() bool {
	return l.currentPos+1 < len(l.src)
}

func (l *Tokenizer) nextPosOK() bool {
	return l.currentPos+2 < len(l.src)
}

func (l *Tokenizer) Tokenize() error {
	for l.currentPosOK() {
		currentChar := rune(l.src[l.currentPos])

		if unicode.IsSpace(currentChar) && currentChar != '\n' {
			l.shift()
			continue
		}

		if currentChar == '\n' {
			l.shift()
			l.currentLine++
			continue
		}

		if currentChar == '/' {
			// comments
			if l.currentPosOK() && l.src[l.currentPos+1] == '/' {
				for l.src[l.currentPos] != '\n' {
					l.shift()
				}
				l.currentLine++
				l.shift()
			}
			continue
		}

		if unicode.IsLetter(currentChar) || currentChar == '_' {
			has, token, err := l.isTokenStarted()

			if err != nil {
				return err
			}

			l.shift()

			if !has {
				continue
			}

			var buffer string
			for l.currentPosOK() && l.src[l.currentPos] != '{' {
				char := l.src[l.currentPos]
				if char == '\n' {
					l.currentLine++
				}
				buffer += string(char)
				l.shift()
			}

			l.addToken(strings.TrimSpace(buffer), *token)
			l.shift()

			buffer = ""
			var braceCount = 0
			for l.currentPosOK() {
				char := l.src[l.currentPos]
				if char == '\n' {
					l.currentLine++
				}

				text := string(char)
				buffer += text
				l.shift()

				if text == "{" {
					braceCount++
				}

				if text == "}" {
					if braceCount > 0 {
						braceCount--
					} else {
						break
					}
				}
			}

			l.addToken(strings.TrimSpace(buffer), TypeValue)
			l.buf = ""

			continue
		}

		return NewErrorWithPosition("Invalid character", Token{
			Line:     l.currentLine,
			FileName: l.fileName,
			Value:    l.buf,
		})
	}

	return nil
}

func (l *Tokenizer) GetTokens() []Token {
	return l.tokens
}

func (l *Tokenizer) isTokenStarted() (bool, *TokenType, error) {
	for name, tokenType := range typeTokens() {
		if strings.HasSuffix(l.buf, name) {
			if l.nextPosOK() && l.src[l.currentPos+1] != ' ' {
				return false, nil, NewErrorWithPosition("Invalid character", Token{
					Value:    name,
					FileName: l.fileName,
					TokenPos: l.currentPos,
					Line:     l.currentLine,
				})
			}
			return true, &tokenType, nil
		}
	}

	char := rune(l.src[l.currentPos+1])
	if l.nextPosOK() && unicode.IsSpace(char) {
		return false, nil, NewErrorWithPosition("Invalid character", Token{
			Value:    l.buf,
			FileName: l.fileName,
			TokenPos: l.currentPos,
			Line:     l.currentLine,
		})
	}
	return false, nil, nil
}

func (l *Tokenizer) GetPosition() string {
	return fmt.Sprintf("File: %v:%v", l.fileName, l.currentLine)
}

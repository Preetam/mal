package reader

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Preetam/mal/dl/types"
)

type Reader struct {
	Tokens []string
}

func ReadString(str string) (types.MalType, error) {
	reader := &Reader{
		Tokens: tokenize(str),
	}
	return reader.ReadForm()
}

func (r *Reader) ReadForm() (types.MalType, error) {
	token := r.Peek()
	if token == "" {
		return nil, errors.New("EOF")
	}
	switch token[0] {
	case '(':
		// Read list
		r.Next()
		val, err := r.readList()
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return r.readAtom()
	}
}

func (r *Reader) readList() (types.MalType, error) {
	list := types.MalList{}
	for {
		token := r.Peek()
		switch token {
		case ")":
			r.Next()
			return list, nil
		case "":
			return nil, errors.New("EOF")
		}
		val, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		list = append(list, val)
	}
	return list, nil
}

func (r *Reader) readAtom() (types.MalType, error) {
	token := r.Next()
	if n, err := strconv.ParseInt(token, 10, 64); err == nil {
		return types.MalInt(n), nil
	}
	return types.MalSymbol(token), nil
}

func (r *Reader) Next() string {
	if len(r.Tokens) == 0 {
		return ""
	}
	token := r.Tokens[0]
	r.Tokens = r.Tokens[1:]
	return token
}

func (r *Reader) Peek() string {
	if len(r.Tokens) == 0 {
		return ""
	}
	return r.Tokens[0]
}

func tokenize(str string) []string {
	tokens := []string{}
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		if (match[1] == "") ||
			// comment
			(match[1][0] == ';') {
			continue
		}
		tokens = append(tokens, match[1])
	}
	return tokens
}

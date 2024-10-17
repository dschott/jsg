package main

import (
	"net/url"
	"path"
	"strings"
	"unicode"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Mapper struct {
	Types       map[string]Type
	Initialisms map[string]struct{}
}

func (m *Mapper) Map(schema *jsonschema.Schema) (*Type, error) {
	if m.Types == nil {
		m.Types = make(map[string]Type)
	}

	var (
		typ Type
		err error
	)

	typ.Pkg = "main"

	typ.Name, err = m.TypeName(schema)
	if err != nil {
		return nil, err
	}

	return &typ, nil
}

func (m *Mapper) TypeName(schema *jsonschema.Schema) (string, error) {
	name := schema.Title
	if name == "" {
		u, err := url.Parse(schema.Location)
		if err != nil {
			return "", err
		}
		base := path.Base(u.Path)
		ext := path.Ext(base)
		name = strings.TrimSuffix(base, ext)
	}
	return m.ToIdentifier(name), nil
}

func (m *Mapper) AddInitialism(initialism string) {
	if m.Initialisms == nil {
		m.Initialisms = make(map[string]struct{})
	}
	m.Initialisms[initialism] = struct{}{}
}

func (m *Mapper) ToIdentifier(s string) string {
	s = strings.TrimSpace(s)

	var (
		output        = make([]rune, 0, len(s)+1)
		lastWord      = make([]rune, 0, 10)
		lastIsUpper   bool
		lastIsDigit   bool
		lastIsLetter  bool
		lastIsSpecial bool
	)

	for i, r := range s {
		isDigit := unicode.IsDigit(r)
		isLetter := !isDigit && unicode.IsLetter(r)
		isUpper := isLetter && unicode.IsUpper(r)
		isLower := isLetter && !isUpper
		isSpecial := !isDigit && !isLetter
		newToken := i == 0 || (isSpecial && !lastIsSpecial) || (isDigit && !lastIsDigit) || (isUpper && !lastIsUpper) || (isLower && !lastIsLetter)

		switch {
		case i == 0 && isDigit:
			output = append(output, '_')
		case newToken && isLower:
			r = unicode.ToUpper(r)
		}

		if len(m.Initialisms) > 0 {
			if newToken {
				m.applyInitialism(output, lastWord)
				lastWord = lastWord[:0]
				if isLetter {
					lastWord = append(lastWord, r)
				}
			} else if len(lastWord) > 0 {
				lastWord = append(lastWord, unicode.ToUpper(r))
			}
		}

		if !isSpecial {
			output = append(output, r)
		}
		lastIsDigit = isDigit
		lastIsUpper = isUpper
		lastIsLetter = isLetter
		lastIsSpecial = isSpecial
	}

	m.applyInitialism(output, lastWord)

	return string(output)
}

func (m *Mapper) applyInitialism(target, word []rune) {
	if _, ok := m.Initialisms[string(word)]; ok {
		copy(target[len(target)-len(word):], word)
	}
}

// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import (
	"errors"
	"strings"
	"unicode"
)

var ErrNilTokenizer = errors.New("nil tokenizer")

type Document interface {
	Text() string
}

type indexedDocument struct {
	d Document
	k map[string]int
	n int
}

type TextDocument struct {
	t string
}

func NewTextDocument(text string) *TextDocument {
	return &TextDocument{
		t: text,
	}
}

func (t *TextDocument) Text() string {
	return t.t
}

func (id *indexedDocument) index(t Tokenizer) error {
	if t == nil {
		return ErrNilTokenizer
	}
	k, err := t(id.d.Text())
	if err != nil {
		return err
	}
	id.k = k
	for _, v := range k {
		id.n += v
	}
	return nil
}

func Tokenize(text string) (map[string]int, error) {
	lower := strings.ToLower(text)
	words := strings.FieldsFunc(lower, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	qk := make(map[string]int)
	for _, w := range words {
		if len(w) >= 2 {
			qk[w]++
		}
	}
	return qk, nil
}

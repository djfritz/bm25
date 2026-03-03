// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

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

// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import "testing"

func TestCollectionNilTokenizer(t *testing.T) {
	d := NewTextDocument("this is a document")
	c := new(Collection)
	err := c.AddDocument(d)
	if err == nil {
		t.Fatal("nil error on nil tokenizer")
	}
}

func TestCollection(t *testing.T) {
	c := new(Collection)
	c.SetTokenizer(Tokenize)
	d := NewTextDocument("this is a document")
	d2 := NewTextDocument("this is my other other document")
	err := c.AddDocument(d)
	if err != nil {
		t.Fatal(err)
	}
	err = c.AddDocument(d2)
	if err != nil {
		t.Fatal(err)
	}

	q := "other document"
	scores, err := c.Score(q, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if scores[0].S != 1.0529239429556823 {
		t.Fatal("invalid score", scores[0].S)
	}
	if scores[1].S != 0.21449594916935835 {
		t.Fatal("invalid score", scores[1].S)
	}
}

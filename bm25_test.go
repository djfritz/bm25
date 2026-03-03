// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import (
	"fmt"
	"testing"
)

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
	scores, err := c.Score(q, 2)
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

func TestCollectionParallel(t *testing.T) {
	c := new(Collection)
	c.SetTokenizer(Tokenize)
	c.SetParallel(31) // should leave a remainder

	for i := 0; i < 1000; i++ {
		d := NewTextDocument(fmt.Sprintf("document number %v with some unique words like alpha%v beta%v", i, i, i))
		err := c.AddDocument(d)
		if err != nil {
			t.Fatal(err)
		}
	}

	scores, err := c.Score("document number", 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 1000 {
		t.Fatal("invalid response", len(scores))
	}
}

func TestCollectionParallelResults(t *testing.T) {
	c := new(Collection)
	c.SetTokenizer(Tokenize)
	c.SetParallel(7)

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
	scores, err := c.Score(q, 2)
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

func TestCollectionSetParallelInvalid(t *testing.T) {
	c := new(Collection)
	err := c.SetParallel(-1)
	if err == nil {
		t.Fatal("nil error on invalid parallel")
	}
}

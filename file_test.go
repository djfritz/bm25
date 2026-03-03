// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import (
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	c := new(Collection)
	c.SetTokenizer(Tokenize)

	d := NewTextDocument("this is a document")
	d2 := NewTextDocument("this is my other other document")
	if err := c.AddDocument(d); err != nil {
		t.Fatal(err)
	}
	if err := c.AddDocument(d2); err != nil {
		t.Fatal(err)
	}

	tmp, err := os.CreateTemp("", "test-*.bm25")
	if err != nil {
		t.Fatal(err)
	}
	f := tmp.Name()
	tmp.Close()
	defer os.Remove(f)
	if err := c.Save(f); err != nil {
		t.Fatal(err)
	}

	c2, err := Load(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(c2.d) != 2 {
		t.Fatal("expected 2 documents, got", len(c2.d))
	}

	if c2.avgdl != c.avgdl {
		t.Fatalf("avgdl mismatch: got %v, want %v", c2.avgdl, c.avgdl)
	}

	if c2.d[0].d.Text() != "this is a document" {
		t.Fatal("document 0 text mismatch")
	}
	if c2.d[1].d.Text() != "this is my other other document" {
		t.Fatal("document 1 text mismatch")
	}

	if c2.d[0].n != c.d[0].n {
		t.Fatalf("document 0 n mismatch: got %v, want %v", c2.d[0].n, c.d[0].n)
	}
	if c2.d[1].k["other"] != 2 {
		t.Fatal("document 1 token count for 'other' mismatch")
	}

	c2.SetTokenizer(Tokenize)
	scores, err := c2.Score("other document", 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if scores[0].S != 1.0529239429556823 {
		t.Fatal("invalid score after load", scores[0].S)
	}
	if scores[1].S != 0.21449594916935835 {
		t.Fatal("invalid score after load", scores[1].S)
	}
}

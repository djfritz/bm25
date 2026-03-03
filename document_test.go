// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import "testing"

func TestDocumentText(t *testing.T) {
	d := NewTextDocument("this is a document")
	if d.Text() != "this is a document" {
		t.Fatal("invalid text")
	}
}

func TestIndexedDocument(t *testing.T) {
	d := NewTextDocument("this is a document, and the document is this")
	id := &indexedDocument{
		d: d,
	}
	err := id.index(Tokenize)
	if err != nil {
		t.Fatal(err)
	}
	if len(id.k) != 5 {
		t.Fatal("invalid token count", len(id.k))
	}
	if id.k["this"] != 2 {
		t.Fatal("invalid token count for 'this'")
	}
	if id.k["document"] != 2 {
		t.Fatal("invalid token count for 'document'")
	}
}

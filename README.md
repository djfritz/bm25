# bm25

A simple [Okapi BM25](https://en.wikipedia.org/wiki/Okapi_BM25) ranking function implemented in Go. It supports document retrieval, parallel scoring, custom tokenizers, and collection serialization.

The implementation was derived from *An Introduction to Information Retrieval*, Manning et al., page 233.

## Install

```
go get github.com/djfritz/bm25
```

## Usage

```go
package main

import (
        "fmt"
        "log"

        "github.com/djfritz/bm25"
)

func main() {
        // Create a collection and set a tokenizer
        c := new(bm25.Collection)
        c.SetTokenizer(bm25.Tokenize)

        // Add documents
        c.AddDocument(bm25.NewTextDocument("the quick brown fox jumps over the lazy dog"))
        c.AddDocument(bm25.NewTextDocument("the lazy dog sat on the porch"))
        c.AddDocument(bm25.NewTextDocument("the fox is quick and brown"))

        // Score documents against a query
        results, err := c.Score("quick fox", 3)
        if err != nil {
                log.Fatal(err)
        }

        for _, r := range results {
                fmt.Printf("score: %.4f  text: %s\n", r.S, r.D.Text())
        }
}
```

## Custom Document Types

Any type that implements the `Document` interface can be added to a collection:

```go
type Document interface {
        Text() string
}
```

## Parallel Scoring

By default, scoring uses `runtime.NumCPU()` goroutines. Override with:

```go
c.SetParallel(4) // use 4 goroutines
c.SetParallel(0) // reset to NumCPU (default)
```

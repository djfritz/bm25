// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

// Package bm25 implements a modified Okapi BM25 ranking function. It supports
// document retrieval, and and differs from BM25 by normalizing ranks on [0,1].
//
// The implementation was derived from "An Introduction to Information
// Retrieval, Manning et al., page 233".
package bm25

import (
	"container/heap"
	"errors"
	"maps"
	"math"
	"runtime"
	"slices"
	"sync"
)

const (
	k1 = 1.5  // free parameter
	b  = 0.75 // free parameter
)

var (
	ErrInvalidThreshold = errors.New("threshold must be between 0 and 1")
	ErrInvalidN         = errors.New("n must be > 0")
	ErrInvalidQ         = errors.New("query must be non-zero length")
	ErrInvalidParallel  = errors.New("number of threads must be >= 0")
)

type Collection struct {
	mtx       sync.Mutex
	d         []*indexedDocument
	avgdl     float64 // average document length in the collection
	tokenizer Tokenizer
	p         int
}

type Tokenizer func(string) (map[string]int, error)

type ScoredDocument struct {
	S float64
	D Document
}

func (c *Collection) SetTokenizer(t Tokenizer) {
	c.tokenizer = t
}

func (c *Collection) SetParallel(p int) error {
	if p < 0 {
		return ErrInvalidParallel
	}
	c.p = p
	return nil
}

func (c *Collection) AddDocument(d Document) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	id := &indexedDocument{
		d: d,
	}
	err := id.index(c.tokenizer)
	if err != nil {
		return err
	}
	c.d = append(c.d, id)
	c.calculateAvgdl()
	return nil
}

func (c *Collection) calculateAvgdl() {
	var total int
	for _, v := range c.d {
		total += v.n
	}
	c.avgdl = float64(total) / float64(len(c.d))
}

type scoredHeap []*ScoredDocument

func (s scoredHeap) Len() int           { return len(s) }
func (s scoredHeap) Less(i, j int) bool { return s[i].S < s[j].S }
func (s scoredHeap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *scoredHeap) Push(x any) {
	*s = append(*s, x.(*ScoredDocument))
}

func (s *scoredHeap) Pop() any {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func (c *Collection) Score(q string, n int, t float64) ([]*ScoredDocument, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if t < 0 || t > 1 {
		return nil, ErrInvalidThreshold
	}
	if n < 1 {
		return nil, ErrInvalidN
	}
	if len(q) == 0 {
		return nil, ErrInvalidQ
	}

	// tokenize the query
	qk, err := c.tokenizeQuery(q)
	if err != nil {
		return nil, err
	}
	idfs := c.idf(qk)

	h := new(scoredHeap)
	heap.Init(h)

	p := c.p
	if p == 0 {
		p = runtime.NumCPU()
	}

	// split documents into p partitions, or the length of documents if it's less than p
	part := p
	if len(c.d) < p {
		part = len(c.d)
	}
	step := len(c.d) / part
	remainder := len(c.d) % part

	coalesce := make(chan *ScoredDocument, len(c.d))
	var wg sync.WaitGroup
	var last int
	for i := range part {
		wg.Add(1)
		start := last
		end := start + step
		if i < remainder {
			end++
		}
		if end > len(c.d) {
			end = len(c.d)
		}
		go scorePartition(coalesce, c.d[start:end], c.avgdl, idfs, qk, &wg)
		last = end
	}

	wg.Wait()
	close(coalesce)

	for s := range coalesce {
		heap.Push(h, s)
		if h.Len() > n {
			heap.Pop(h)
		}
	}

	var r []*ScoredDocument
	for h.Len() > 0 {
		r = append(r, heap.Pop(h).(*ScoredDocument))
	}
	slices.Reverse(r)
	return r, nil
}

func (c *Collection) tokenizeQuery(q string) ([]string, error) {
	if c.tokenizer == nil {
		return nil, ErrNilTokenizer
	}
	m, err := c.tokenizer(q)
	if err != nil {
		return nil, err
	}
	return slices.Sorted(maps.Keys(m)), nil
}

func (c *Collection) idf(qk []string) []float64 {
	var idfs []float64
	for _, v := range qk {
		var nq int
		for _, w := range c.d {
			if _, ok := w.k[v]; ok {
				nq++
			}
		}
		idf := math.Log(((float64(len(c.d)) - float64(nq) + 0.5) / (float64(nq) + 0.5)) + 1)
		idfs = append(idfs, idf)
	}
	return idfs
}

func score(id *indexedDocument, avgdl float64, idfs []float64, qk []string) float64 {
	var s float64
	for i, v := range qk {
		idf := idfs[i]
		fdq := id.k[v]
		if fdq == 0 {
			continue
		}
		x := idf * ((float64(fdq) * (k1 + 1)) / (float64(fdq) + k1*(1-b+b*(float64(id.n)/avgdl))))
		s += x
	}
	return s
}

func scorePartition(coalesce chan *ScoredDocument, docs []*indexedDocument, avgdl float64, idfs []float64, qk []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range docs {
		s := score(v, avgdl, idfs, qk)
		coalesce <- &ScoredDocument{
			S: s,
			D: v.d,
		}
	}
}

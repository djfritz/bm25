// Copyright 2026 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package bm25

import (
	"bytes"
	"encoding/gob"
	"os"
)

func init() {
	gob.Register(&TextDocument{})
}

func (t *TextDocument) GobEncode() ([]byte, error) {
	return []byte(t.t), nil
}

func (t *TextDocument) GobDecode(data []byte) error {
	t.t = string(data)
	return nil
}

func (id *indexedDocument) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&id.d); err != nil {
		return nil, err
	}
	if err := enc.Encode(id.k); err != nil {
		return nil, err
	}
	if err := enc.Encode(id.n); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (id *indexedDocument) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&id.d); err != nil {
		return err
	}
	if err := dec.Decode(&id.k); err != nil {
		return err
	}
	if err := dec.Decode(&id.n); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Save(f string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(c.d); err != nil {
		return err
	}
	if err := enc.Encode(c.avgdl); err != nil {
		return err
	}
	return nil
}

func Load(f string) (*Collection, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := &Collection{}
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&c.d); err != nil {
		return nil, err
	}
	if err := dec.Decode(&c.avgdl); err != nil {
		return nil, err
	}
	return c, nil
}

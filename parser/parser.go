// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

//go:generate pigeon -o grammar.peg.go ./grammar.peg
//go:generate goimports -w ./grammar.peg.go

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Filesystem interface {
	Open(filename string) (io.ReadCloser, error)
	// Abs makes "path" absolute, when relative to the directory "dir".
	Abs(dir, path string) (string, error)
}

type Parser struct {
	Filesystem Filesystem // For handling includes. Can be set to nil to fall back to os package.
	Files      map[string]*Thrift
}

func New() *Parser {
	return &Parser{
		Files: map[string]*Thrift{},
	}
}

func (p *Parser) Parse(r io.Reader, opts ...Option) (*Thrift, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	name := "<reader>"
	if named, ok := r.(namedReader); ok {
		name = named.Name()
	}
	i, err := Parse(name, b, opts...)
	if err != nil {
		return nil, err
	}
	t := i.(*Thrift)
	t.Filename = name
	return t, nil
}

// ParseFile parses a thrift file *and all the included thrift files*.
// Produces in output a map[full path name]=>parsed Thrift struct which includes all the
// thrift files that have been parsed so far.
// The parser memoizes the result of the parsing, so each thrift file is actually parsed only once.
func (p *Parser) ParseFile(filename string, opts ...Option) (map[string]*Thrift, string, error) {
	absPath, err := p.abs("", filename)
	if err != nil {
		return nil, "", err
	}

	pathsToParse := []string{absPath}
	for len(pathsToParse) > 0 {
		path := pathsToParse[len(pathsToParse)-1]
		pathsToParse = pathsToParse[:len(pathsToParse)-1]
		if _, ok := p.Files[path]; ok {
			break
		}
		rd, err := p.open(path)
		if err != nil {
			return nil, "", err
		}
		thrift, err := p.Parse(rd, opts...)
		rd.Close()
		if err != nil {
			return nil, "", err
		}
		p.Files[path] = thrift

		basePath := filepath.Dir(path)
		for incName, incPath := range thrift.Includes {
			includePath, err := p.abs(basePath, incPath)
			if err != nil {
				return nil, "", err
			}
			thrift.Includes[incName] = includePath
			if _, ok := p.Files[includePath]; !ok {
				pathsToParse = append(pathsToParse, includePath)
			}
		}
	}

	return p.Files, absPath, nil
}

func (p *Parser) open(path string) (io.ReadCloser, error) {
	if p.Filesystem == nil {
		return os.Open(path)
	}
	return p.Filesystem.Open(path)
}

func (p *Parser) abs(dir, path string) (string, error) {
	if p.Filesystem == nil {
		absPath, err := filepath.Abs(filepath.Join(dir, path))
		if err != nil {
			return "", err
		}
		return filepath.Clean(absPath), nil
	}
	return p.Filesystem.Abs(dir, path)
}

type namedReader interface {
	Name() string
}

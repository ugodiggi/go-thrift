// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.
//
// types.go contains the type definition for the types that the parser parses out
// of the thrift files. It should be read side-by-side with grammar.peg.
// Look at the Thrift struct for the top-level output of parsing a thrift file.

package parser

import "fmt"

// Pos represents a token position in a file.
type Pos struct {
	Line int
	Col  int
}

// TemplateInstance contains the reference to the template that this type is an instance of,
// for types that are instances of a template.
type TemplateInstance struct {
	TemplateName string  `json:",omitempty"`
	TypeArgs     []*Type `json:",omitempty"`
}

// Type is the thrift type of a piece of data.
// When parsing, the thrift type of a piece of data (e.g. the type of thrift field, of a method
// argument or response type, of a constant, ...), will be parsed as a Type object; the definition
// of a thrift type will be parsed as a specific struct (Typedef, Enum, SEnum, Field, Struct, ...)
type Type struct {
	Pos         Pos
	Name        string        `json:",omitempty"`
	KeyType     *Type         `json:",omitempty"` // If map
	ValueType   *Type         `json:",omitempty"` // If map, list, or set
	Annotations []*Annotation `json:",omitempty"`

	TemplateInstance *TemplateInstance `json:",omitempty"` // If template instance
}

// Typedef is the definition of a thrift typedef.
type Typedef struct {
	*Type

	Pos         Pos
	Alias       string
	Annotations []*Annotation `json:",omitempty"`
}

type EnumValue struct {
	Pos         Pos
	Comment     string
	Name        string
	Value       int
	Annotations []*Annotation `json:",omitempty"`
}

// Enum is the definition of a thrift enum.
type Enum struct {
	Pos         Pos
	Comment     string
	Name        string
	Values      map[string]*EnumValue
	Annotations []*Annotation `json:",omitempty"`
}

type SEnumValue struct {
	Pos         Pos
	Comment     string
	Value       string
	Annotations []*Annotation `json:",omitempty"`
}

// SEnum is the definition of a thrift senum - analogous to an enum, but its value is the string value.
type SEnum struct {
	Pos         Pos
	Comment     string
	Name        string
	Values      map[string]*SEnumValue
	Annotations []*Annotation `json:",omitempty"`
}

type Constant struct {
	Pos     Pos
	Comment string
	Name    string
	Type    *Type
	Value   interface{}
}

// Field is the definition of a struct's field.
type Field struct {
	Pos         Pos
	Comment     string
	ID          int
	Name        string
	Optional    bool
	Type        *Type
	Default     interface{}   `json:",omitempty"`
	Annotations []*Annotation `json:",omitempty"`
}

// Struct is the definition of a thrift struct.
type Struct struct {
	Pos         Pos
	Comment     string
	Name        string
	Fields      []*Field
	Annotations []*Annotation `json:",omitempty"`
}

// TemplateDef is the definition of a templated thrift struct.
type TemplateDef struct {
	Struct

	TypeArgNames []string `json:",omitempty"`
}

// Method is the definition of a thrift method.
type Method struct {
	Pos         Pos
	Comment     string
	Name        string
	Oneway      bool
	ReturnType  *Type
	Arguments   []*Field
	Exceptions  []*Field      `json:",omitempty"`
	Annotations []*Annotation `json:",omitempty"`
}

// Service is the definition of a thrift service.
type Service struct {
	Pos         Pos
	Comment     string
	Name        string
	Extends     string `json:",omitempty"`
	Methods     map[string]*Method
	Annotations []*Annotation `json:",omitempty"`
}

// Thrift is the output of parsing a whole thrift file.
type Thrift struct {
	Filename     string
	Includes     map[string]string       `json:",omitempty"` // name -> unique identifier (absolute path generally)
	Imports      map[string]*Thrift      `json:",omitempty"` // name -> imported file
	Typedefs     map[string]*Typedef     `json:",omitempty"`
	Namespaces   map[string]string       `json:",omitempty"`
	Constants    map[string]*Constant    `json:",omitempty"`
	Enums        map[string]*Enum        `json:",omitempty"`
	SEnums       map[string]*SEnum       `json:",omitempty"`
	Structs      map[string]*Struct      `json:",omitempty"`
	Exceptions   map[string]*Struct      `json:",omitempty"`
	Unions       map[string]*Struct      `json:",omitempty"`
	TemplateDefs map[string]*TemplateDef `json:",omitempty"`
	Services     map[string]*Service     `json:",omitempty"`
}

type Identifier string

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

type Annotation struct {
	Pos   Pos
	Name  string
	Value string
}

func (t *Type) String() string {
	switch t.Name {
	case "map":
		return fmt.Sprintf("map<%s,%s>", t.KeyType.String(), t.ValueType.String())
	case "list":
		return fmt.Sprintf("list<%s>", t.ValueType.String())
	case "set":
		return fmt.Sprintf("set<%s>", t.ValueType.String())
	}
	return t.Name
}

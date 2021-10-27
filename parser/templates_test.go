package parser

import "testing"

func init() {
	inTests = true
}

func TestParseTemplates(t *testing.T) {
	thrift, err := parse(`// Some comment
template SomeTemplate<T0, T1> {
  1: optional T0 field0;
  2: optional T1 field1;
}

template Template2<Type0> {
  1: optional Type0 field;
}

struct SomeStruct {
	1: optional SomeTemplate<string, i32> templateField;
	2: optional string abc
	3: optional Template2<double> template2;
}`)

	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	expectedStructs := []*Struct{
		&Struct{
			Name: "SomeStruct",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "templateField",
					Optional: true,
					Type: &Type{
						Name: "SomeTemplate__string__i32",
						TemplateInstance: &TemplateInstance{
							TemplateName: "SomeTemplate",
							TypeArgs: []*Type{
								&Type{Name: "string"},
								&Type{Name: "i32"},
							},
						},
					},
				},
				{
					ID:       2,
					Name:     "abc",
					Optional: true,
					Type: &Type{
						Name: "string",
					},
				},
				{
					ID:       3,
					Name:     "template2",
					Optional: true,
					Type: &Type{
						Name: "Template2__double",
						TemplateInstance: &TemplateInstance{
							TemplateName: "Template2",
							TypeArgs: []*Type{
								&Type{Name: "double"},
							},
						},
					},
				},
			},
		},
		&Struct{
			Name:    "SomeTemplate__string__i32",
			Comment: "Some comment",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "field0",
					Optional: true,
					Type: &Type{
						Name: "string",
					},
				},
				{
					ID:       2,
					Name:     "field1",
					Optional: true,
					Type: &Type{
						Name: "i32",
					},
				},
			},
		},
		&Struct{
			Name: "Template2__double",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "field",
					Optional: true,
					Type: &Type{
						Name: "double",
					},
				},
			},
		},
	}
	assertStructsEqual(t, expectedStructs, thrift.Structs)
}

func TestParseTemplate1(t *testing.T) {
	thrift, err := parse(`// Some comment
template Template1<T0> {
  1: optional T0 field0;
}

struct SomeStruct {
	1: optional Template1<string> structField;
}`)

	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	expectedStructs := []*Struct{
		&Struct{
			Name: "SomeStruct",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "structField",
					Optional: true,
					Type: &Type{
						Name: "Template1__string",
						TemplateInstance: &TemplateInstance{
							TemplateName: "Template1",
							TypeArgs: []*Type{
								&Type{Name: "string"},
							},
						},
					},
				},
			},
		},
		&Struct{
			Name:    "Template1__string",
			Comment: "Some comment",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "field0",
					Optional: true,
					Type: &Type{
						Name: "string",
					},
				},
			},
		},
	}
	assertStructsEqual(t, expectedStructs, thrift.Structs)
}

func TestParseTemplate3(t *testing.T) {
	thrift, err := parse(`// Some comment
template Template3<T0, T1, T2> {
  1: optional T0 field0;
  2: optional T1 field1;
  3: optional T2 field2;
}

struct SomeStruct {
	1: optional Template3<string, i32, double> structField;
}`)

	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	expectedStructs := []*Struct{
		&Struct{
			Name: "SomeStruct",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "structField",
					Optional: true,
					Type: &Type{
						Name: "Template3__string__i32__double",
						TemplateInstance: &TemplateInstance{
							TemplateName: "Template3",
							TypeArgs: []*Type{
								&Type{Name: "string"},
								&Type{Name: "i32"},
								&Type{Name: "double"},
							},
						},
					},
				},
			},
		},
		&Struct{
			Name:    "Template3__string__i32__double",
			Comment: "Some comment",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "field0",
					Optional: true,
					Type: &Type{
						Name: "string",
					},
				},
				{
					ID:       2,
					Name:     "field1",
					Optional: true,
					Type: &Type{
						Name: "i32",
					},
				},
				{
					ID:       3,
					Name:     "field2",
					Optional: true,
					Type: &Type{
						Name: "double",
					},
				},
			},
		},
	}
	assertStructsEqual(t, expectedStructs, thrift.Structs)
}

func TestParseTemplateWithContainers(t *testing.T) {
	thrift, err := parse(`// Some comment
template TemplateWithContainers<T0> {
  1: optional list<T0> field1;
  2: optional map<string, T0> field2;
}

struct SomeStruct {
	1: optional TemplateWithContainers<i64> structField;
}`)

	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	expectedStructs := []*Struct{
		&Struct{
			Name: "SomeStruct",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "structField",
					Optional: true,
					Type: &Type{
						Name: "TemplateWithContainers__i64",
						TemplateInstance: &TemplateInstance{
							TemplateName: "TemplateWithContainers",
							TypeArgs: []*Type{
								&Type{Name: "i64"},
							},
						},
					},
				},
			},
		},
		&Struct{
			Name:    "TemplateWithContainers__i64",
			Comment: "Some comment",
			Fields: []*Field{
				{
					ID:       1,
					Name:     "field1",
					Optional: true,
					Type: &Type{
						Name: "list",
						ValueType: &Type{
							Name: "i64",
						},
					},
				},
				{
					ID:       2,
					Name:     "field2",
					Optional: true,
					Type: &Type{
						Name: "map",
						KeyType: &Type{
							Name: "string",
						},
						ValueType: &Type{
							Name: "i64",
						},
					},
				},
			},
		},
	}
	assertStructsEqual(t, expectedStructs, thrift.Structs)
}

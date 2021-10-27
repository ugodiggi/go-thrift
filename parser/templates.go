package parser

import "fmt"

func tmplInstanceName(ti *TemplateInstance) string {
	name := ti.TemplateName + "__"
	for i, a := range ti.TypeArgs {
		if i > 0 {
			name += "__"
		}
		name += a.Name
	}
	return name
}

func tmplResolveArgs(typeArgNames []string, typeArgs []*Type, fields []*Field) ([]*Field, error) {
	tas := make(map[string]*Type)
	if len(typeArgs) != len(typeArgNames) {
		return nil, fmt.Errorf("Template args mismatch. Expected %d args, got %d", len(typeArgNames), len(typeArgs))
	}
	for i, n := range typeArgNames {
		tas[n] = typeArgs[i]
	}

	result := []*Field{}
	for _, f := range fields {
		ff := &Field{
			Pos:         f.Pos,
			Comment:     f.Comment,
			ID:          f.ID,
			Name:        f.Name,
			Optional:    f.Optional,
			Type:        f.Type,
			Default:     f.Default,
			Annotations: f.Annotations,
		}
		if ta, ok := tas[ff.Type.Name]; ok {
			ff.Type = ta
		}
		result = append(result, ff)
	}
	return result, nil
}

func (p *Parser) RenderTemplates() (*Parser, error) {
	// Find all template instances.
	instances := make(map[string]*TemplateInstance)
	collectType := func(t *Type) {
		if t == nil {
			return
		}
		if ti := t.TemplateInstance; ti != nil {
			name := tmplInstanceName(ti)
			if _, ok := instances[name]; !ok {
				instances[name] = ti
			}
		}
	}
	collect := func(s *Struct) {
		for _, f := range s.Fields {
			collectType(f.Type)
		}
	}
	for _, f := range p.Files {
		for _, s := range f.Structs {
			collect(s)
		}
		for _, s := range f.Exceptions {
			collect(s)
		}
		for _, s := range f.Unions {
			collect(s)
		}
		for _, s := range f.Services {
			for _, m := range s.Methods {
				collectType(m.ReturnType)
				for _, arg := range m.Arguments {
					collectType(arg.Type)
				}
				for _, exc := range m.Exceptions {
					collectType(exc.Type)
				}
			}
		}
	}

	instancesLeft := []*TemplateInstance{}
	for _, i := range instances {
		instancesLeft = append(instancesLeft, i)
	}

	// Traverse all template instances, rendering them.
	// More TemplateInstances may come out now.
	for n := len(instancesLeft); n > 0; n = len(instancesLeft) {
		var ti *TemplateInstance
		ti, instancesLeft = instancesLeft[n-1], instancesLeft[:n-1]
		found := false
		for _, f := range p.Files {
			tDef, ok := f.TemplateDefs[ti.TemplateName]
			if !ok {
				continue
			} else {
				found = true
			}

			name := tmplInstanceName(ti)
			fields, err := tmplResolveArgs(tDef.TypeArgNames, ti.TypeArgs, tDef.Fields)
			if err != nil {
				return nil, err
			}
			f.Structs[name] = &Struct{
				Pos:         tDef.Pos,
				Comment:     tDef.Comment,
				Name:        name,
				Fields:      fields,
				Annotations: tDef.Annotations,
			}
			break
		}
		if !found {
			return nil, fmt.Errorf("Undefined template %s", ti.TemplateName)
		}
	}

	return &Parser{
		Filesystem: p.Filesystem,
		Files:      p.Files,
	}, nil
}

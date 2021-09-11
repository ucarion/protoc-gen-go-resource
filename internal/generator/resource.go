package generator

import (
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type resource struct {
	NameField     *protogen.Field
	ParseFunc     protogen.GoIdent
	FullParseFunc protogen.GoIdent
	ParsedType    protogen.GoIdent
	Type          resourceType
	Patterns      []pattern
}

func newMessageResource(m *protogen.Message) (*resource, error) {
	o := m.Desc.Options().(*descriptorpb.MessageOptions)
	d := proto.GetExtension(o, annotations.E_Resource).(*annotations.ResourceDescriptor)
	if d == nil {
		return nil, nil
	}

	r, err := newBareResource(m.GoIdent.GoImportPath, d)
	if err != nil {
		return nil, err
	}

	fieldName := "name" // the default, unless otherwise specified
	if d.NameField != "" {
		fieldName = d.NameField
	}

	for _, f := range m.Fields {
		if string(f.Desc.Name()) == fieldName {
			r.NameField = f
		}
	}

	if r.NameField == nil {
		return nil, fmt.Errorf("%v specifies %q as name field, but no field with that name exists", m.GoIdent, fieldName)
	}

	return r, nil
}

func newFileResources(f *protogen.File) ([]*resource, error) {
	var out []*resource

	o := f.Desc.Options().(*descriptorpb.FileOptions)
	rs := proto.GetExtension(o, annotations.E_ResourceDefinition).([]*annotations.ResourceDescriptor)
	for _, r := range rs {
		res, err := newBareResource(f.GoImportPath, r)
		if err != nil {
			return nil, err
		}

		out = append(out, res)
	}

	return out, nil
}

func newBareResource(importPath protogen.GoImportPath, d *annotations.ResourceDescriptor) (*resource, error) {
	t, err := newResourceType(d.Type)
	if err != nil {
		return nil, err
	}

	var patterns []pattern
	for _, s := range d.Pattern {
		p, err := newPattern(s)
		if err != nil {
			return nil, err
		}

		patterns = append(patterns, p)
	}

	parseFunc := protogen.GoIdent{
		GoName:       "Parse" + t.TypeName + "Name",
		GoImportPath: importPath,
	}

	parseFullFunc := protogen.GoIdent{
		GoName:       "ParseFull" + t.TypeName + "Name",
		GoImportPath: importPath,
	}

	parsedType := protogen.GoIdent{
		GoName:       "Parsed" + t.TypeName + "Name",
		GoImportPath: importPath,
	}

	return &resource{
		ParseFunc:     parseFunc,
		ParsedType:    parsedType,
		FullParseFunc: parseFullFunc,
		Type:          t,
		Patterns:      patterns,
	}, nil
}

type resourceType struct {
	ServiceName string
	TypeName    string
}

func newResourceType(s string) (resourceType, error) {
	i := strings.IndexByte(s, '/')
	if i == -1 {
		return resourceType{}, fmt.Errorf("invalid resource type: %q", s)
	}

	return resourceType{ServiceName: s[:i], TypeName: s[i+1:]}, nil
}

package generator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type pattern []segment

func newPattern(s string) (pattern, error) {
	var out pattern
	var vars int
	for _, part := range strings.Split(s, "/") {
		if strings.HasPrefix(part, "{") {
			if !strings.HasSuffix(part, "}") {
				return nil, fmt.Errorf("invalid segment: %q", part)
			}

			out = append(out, segment{
				Name: strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}"),
				VarIndex: vars,
			})

			vars++
		} else {
			out = append(out, segment{Name: part, VarIndex: -1})
		}
	}

	return out, nil
}

type segment struct {
	Name     string
	VarIndex int
}

func (s *segment) Var() bool {
	return s.VarIndex != -1
}

func (s *segment) FieldName() string {
	return strcase.ToCamel(s.Name + "_ID")
}

// Code generated by protoc-gen-go-resource. DO NOT EDIT.
package external

import (
	"fmt"
	"strings"
)

type ParsedExternalName struct {
	ExternalID string
}

func ParseExternalName(s string) (ParsedExternalName, error) {
	p := strings.Split(s, "/")
	if len(p) != 2 {
		return ParsedExternalName{}, fmt.Errorf("parse %q: bad number of segments, want: 2, got: %d", s, len(p))
	}

	var out ParsedExternalName

	if p[0] != "external" {
		return ParsedExternalName{}, fmt.Errorf("parse %q: bad segment 0, want: %q, got: %q", s, "external", p[0])
	}

	out.ExternalID = p[1]

	return out, nil
}

func ParseFullExternalName(s string) (ParsedExternalName, error) {
	if !strings.HasPrefix(s, "//example.com/") {
		return ParsedExternalName{}, fmt.Errorf("parse %q: invalid prefix, want: %q", s, "//example.com/")
	}

	return ParseExternalName(strings.TrimPrefix(s, "//example.com/"))
}

func (n ParsedExternalName) Name() string {
	var out string

	out += "external"

	out += "/"

	out += n.ExternalID

	return out
}

func (n ParsedExternalName) FullName() string {
	return "//example.com/" + n.Name()
}

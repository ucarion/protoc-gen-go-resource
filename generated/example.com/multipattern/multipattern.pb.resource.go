// Code generated by protoc-gen-go-resource. DO NOT EDIT.
package multipattern

import (
	"fmt"
	"strings"
)

type ParsedBookName_0 struct {
	PublisherID string

	BookID string
}

func ParseBookName_0(s string) (ParsedBookName_0, error) {
	p := strings.Split(s, "/")
	if len(p) != 4 {
		return ParsedBookName_0{}, fmt.Errorf("parse %q: bad number of segments, want: 4, got: %d", s, len(p))
	}

	var out ParsedBookName_0

	if p[0] != "publishers" {
		return ParsedBookName_0{}, fmt.Errorf("parse %q: bad segment 0, want: %q, got: %q", s, "publishers", p[0])
	}

	out.PublisherID = p[1]

	if p[2] != "books" {
		return ParsedBookName_0{}, fmt.Errorf("parse %q: bad segment 2, want: %q, got: %q", s, "books", p[2])
	}

	out.BookID = p[3]

	return out, nil
}

func ParseBookNameFull_0(s string) (ParsedBookName_0, error) {
	if !strings.HasPrefix(s, "//example.com/") {
		return ParsedBookName_0{}, fmt.Errorf("parse %q: invalid prefix, want: %q", s, "//example.com/")
	}

	return ParseBookName_0(strings.TrimPrefix(s, "//example.com/"))
}

func (n ParsedBookName_0) Name() string {
	var out string

	out += "publishers"

	out += "/"

	out += n.PublisherID

	out += "/"

	out += "books"

	out += "/"

	out += n.BookID

	return out
}

func (n ParsedBookName_0) FullName() string {
	return "//example.com/" + n.Name()
}

func (n ParsedBookName_0) mustEmbedParsedBookName() {}

type ParsedBookName_1 struct {
	AuthorID string

	BookID string
}

func ParseBookName_1(s string) (ParsedBookName_1, error) {
	p := strings.Split(s, "/")
	if len(p) != 4 {
		return ParsedBookName_1{}, fmt.Errorf("parse %q: bad number of segments, want: 4, got: %d", s, len(p))
	}

	var out ParsedBookName_1

	if p[0] != "authors" {
		return ParsedBookName_1{}, fmt.Errorf("parse %q: bad segment 0, want: %q, got: %q", s, "authors", p[0])
	}

	out.AuthorID = p[1]

	if p[2] != "books" {
		return ParsedBookName_1{}, fmt.Errorf("parse %q: bad segment 2, want: %q, got: %q", s, "books", p[2])
	}

	out.BookID = p[3]

	return out, nil
}

func ParseBookNameFull_1(s string) (ParsedBookName_1, error) {
	if !strings.HasPrefix(s, "//example.com/") {
		return ParsedBookName_1{}, fmt.Errorf("parse %q: invalid prefix, want: %q", s, "//example.com/")
	}

	return ParseBookName_1(strings.TrimPrefix(s, "//example.com/"))
}

func (n ParsedBookName_1) Name() string {
	var out string

	out += "authors"

	out += "/"

	out += n.AuthorID

	out += "/"

	out += "books"

	out += "/"

	out += n.BookID

	return out
}

func (n ParsedBookName_1) FullName() string {
	return "//example.com/" + n.Name()
}

func (n ParsedBookName_1) mustEmbedParsedBookName() {}

type ParsedBookName interface {
	Name() string
	FullName() string
	mustEmbedParsedBookName()
}

func ParseBookName(s string) (ParsedBookName, error) {
	var errs []string
	var res ParsedBookName
	var err error

	res, err = ParseBookName_0(s)
	if err == nil {
		return res, nil
	}

	errs = append(errs, fmt.Sprintf("pattern 0: %v", err))

	res, err = ParseBookName_1(s)
	if err == nil {
		return res, nil
	}

	errs = append(errs, fmt.Sprintf("pattern 1: %v", err))

	return nil, fmt.Errorf("no pattern matches input: %v", strings.Join(errs, "; "))
}

func ParseFullBookName(s string) (ParsedBookName, error) {
	if !strings.HasPrefix(s, "//example.com/") {
		return nil, fmt.Errorf("parse %q: invalid prefix, want: %q", s, "//example.com/")
	}

	return ParseBookName(strings.TrimPrefix(s, "//example.com/"))
}

func (x *Book) ParseName() (ParsedBookName, error) {
	return ParseBookName(x.Name)
}

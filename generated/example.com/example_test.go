package example_test

import (
	"testing"

	"example.com/external"
	"example.com/multipattern"
	"example.com/namefield"
	"example.com/reference"
	"example.com/simple"
	"github.com/google/go-cmp/cmp"
)

func TestSimple_ParseFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out simple.ParsedThingName
		Err string
	}{
		{In: "things/foo", Out: simple.ParsedThingName{ThingID: "foo"}},
		{In: "thing/foo", Err: `parse "thing/foo": bad segment 0, want: "things", got: "thing"`},
		{In: "things/foo/bar", Err: `parse "things/foo/bar": bad number of segments, want: 2, got: 3`},
		{In: "", Err: `parse "": bad number of segments, want: 2, got: 1`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := simple.ParseThingName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestSimple_ParseFullFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out simple.ParsedThingName
		Err string
	}{
		{In: "//example.com/things/foo", Out: simple.ParsedThingName{ThingID: "foo"}},
		{In: "example.com/things/foo", Err: `parse "example.com/things/foo": invalid prefix, want: "//example.com/"`},
		{In: "//foo.example.com/things/foo", Err: `parse "//foo.example.com/things/foo": invalid prefix, want: "//example.com/"`},
		{In: "//example.com/thing/foo", Err: `parse "thing/foo": bad segment 0, want: "things", got: "thing"`},
		{In: "//example.com/things/foo/bar", Err: `parse "things/foo/bar": bad number of segments, want: 2, got: 3`},
		{In: "", Err: `parse "": invalid prefix, want: "//example.com/"`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := simple.ParseFullThingName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestSimple_NameFunc(t *testing.T) {
	parsed := simple.ParsedThingName{ThingID: "foo"}
	if d := cmp.Diff("things/foo", parsed.Name()); d != "" {
		t.Fatalf("bad Name:\n%s", d)
	}
}

func TestSimple_FullNameFunc(t *testing.T) {
	parsed := simple.ParsedThingName{ThingID: "foo"}
	if d := cmp.Diff("//example.com/things/foo", parsed.FullName()); d != "" {
		t.Fatalf("bad FullName:\n%s", d)
	}
}

func TestSimple_ParseMultipartFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out simple.ParsedProjectThingName
		Err string
	}{
		{In: "projects/foo/things/bar", Out: simple.ParsedProjectThingName{ProjectID: "foo", ThingID: "bar"}},
		{In: "project/foo/things/bar", Err: `parse "project/foo/things/bar": bad segment 0, want: "projects", got: "project"`},
		{In: "projects/foo/thing/bar", Err: `parse "projects/foo/thing/bar": bad segment 2, want: "things", got: "thing"`},
		{In: "projects/foo", Err: `parse "projects/foo": bad number of segments, want: 4, got: 2`},
		{In: "", Err: `parse "": bad number of segments, want: 4, got: 1`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := simple.ParseProjectThingName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestSimple_ParseMultipartFullFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out simple.ParsedProjectThingName
		Err string
	}{
		{In: "//example.com/projects/foo/things/bar", Out: simple.ParsedProjectThingName{ProjectID: "foo", ThingID: "bar"}},
		{In: "example.com/projects/foo/things/bar", Err: `parse "example.com/projects/foo/things/bar": invalid prefix, want: "//example.com/"`},
		{In: "//foo.example.com/projects/foo/things/bar", Err: `parse "//foo.example.com/projects/foo/things/bar": invalid prefix, want: "//example.com/"`},
		{In: "//example.com/project/foo/things/bar", Err: `parse "project/foo/things/bar": bad segment 0, want: "projects", got: "project"`},
		{In: "//example.com/projects/foo/thing/bar", Err: `parse "projects/foo/thing/bar": bad segment 2, want: "things", got: "thing"`},
		{In: "//example.com/projects/foo", Err: `parse "projects/foo": bad number of segments, want: 4, got: 2`},
		{In: "", Err: `parse "": invalid prefix, want: "//example.com/"`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := simple.ParseFullProjectThingName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestSimple_MultipartNameFunc(t *testing.T) {
	parsed := simple.ParsedProjectThingName{ProjectID: "foo", ThingID: "bar"}
	if d := cmp.Diff("projects/foo/things/bar", parsed.Name()); d != "" {
		t.Fatalf("bad Name:\n%s", d)
	}
}

func TestSimple_MultipartFullNameFunc(t *testing.T) {
	parsed := simple.ParsedProjectThingName{ProjectID: "foo", ThingID: "bar"}
	if d := cmp.Diff("//example.com/projects/foo/things/bar", parsed.FullName()); d != "" {
		t.Fatalf("bad Name:\n%s", d)
	}
}

func TestSimple_ParseMethod(t *testing.T) {
	thing := simple.Thing{Name: "things/foo"}
	got, err := thing.ParseName()
	if err != nil {
		t.Fatalf("bad parse err: %v", err)
	}

	if d := cmp.Diff(simple.ParsedThingName{ThingID: "foo"}, got); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}
}

func TestExternal_ParseFunc(t *testing.T) {
	got, err := external.ParseExternalName("external/foo")
	if err != nil {
		t.Fatalf("bad parse err: %v", err)
	}

	if d := cmp.Diff(external.ParsedExternalName{ExternalID: "foo"}, got); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}
}

func TestMultiPattern_ParseFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out multipattern.ParsedBookName
		Err string
	}{
		{In: "publishers/foo/books/bar", Out: multipattern.ParsedBookName_0{PublisherID: "foo", BookID: "bar"}},
		{In: "authors/foo/books/bar", Out: multipattern.ParsedBookName_1{AuthorID: "foo", BookID: "bar"}},
		{In: "publisher/foo/books/bar", Err: `no pattern matches input: pattern 0: parse "publisher/foo/books/bar": bad segment 0, want: "publishers", got: "publisher"; pattern 1: parse "publisher/foo/books/bar": bad segment 0, want: "authors", got: "publisher"`},
		{In: "", Err: `no pattern matches input: pattern 0: parse "": bad number of segments, want: 4, got: 1; pattern 1: parse "": bad number of segments, want: 4, got: 1`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := multipattern.ParseBookName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestMultiPattern_ParseFullFunc(t *testing.T) {
	testCases := []struct {
		In  string
		Out multipattern.ParsedBookName
		Err string
	}{
		{In: "//example.com/publishers/foo/books/bar", Out: multipattern.ParsedBookName_0{PublisherID: "foo", BookID: "bar"}},
		{In: "//example.com/authors/foo/books/bar", Out: multipattern.ParsedBookName_1{AuthorID: "foo", BookID: "bar"}},
		{In: "//example.com/publisher/foo/books/bar", Err: `no pattern matches input: pattern 0: parse "publisher/foo/books/bar": bad segment 0, want: "publishers", got: "publisher"; pattern 1: parse "publisher/foo/books/bar": bad segment 0, want: "authors", got: "publisher"`},
		{In: "example.com/publishers/foo/books/bar", Err: `parse "example.com/publishers/foo/books/bar": invalid prefix, want: "//example.com/"`},
		{In: "", Err: `parse "": invalid prefix, want: "//example.com/"`},
	}

	for _, tt := range testCases {
		t.Run(tt.In, func(t *testing.T) {
			out, err := multipattern.ParseFullBookName(tt.In)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if d := cmp.Diff(tt.Err, errMsg); d != "" {
				t.Fatalf("bad err:\n%s", d)
			}

			if d := cmp.Diff(tt.Out, out); d != "" {
				t.Fatalf("bad out:\n%s", d)
			}
		})
	}
}

func TestMultiPattern_NameFunc(t *testing.T) {
	var parsed multipattern.ParsedBookName

	parsed = multipattern.ParsedBookName_0{PublisherID: "foo", BookID: "bar"}
	if d := cmp.Diff("publishers/foo/books/bar", parsed.Name()); d != "" {
		t.Fatalf("bad Name:\n%s", d)
	}

	parsed = multipattern.ParsedBookName_1{AuthorID: "foo", BookID: "bar"}
	if d := cmp.Diff("authors/foo/books/bar", parsed.Name()); d != "" {
		t.Fatalf("bad Name:\n%s", d)
	}
}

func TestMultiPattern_ParseMethod(t *testing.T) {
	thing := multipattern.Book{Name: "authors/foo/books/bar"}
	got, err := thing.ParseName()
	if err != nil {
		t.Fatalf("bad parse err: %v", err)
	}

	if d := cmp.Diff(multipattern.ParsedBookName_1{AuthorID: "foo", BookID: "bar"}, got); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}
}

func TestNameField_ParseFunc(t *testing.T) {
	person := namefield.Person{PersonName: "persons/foo"}
	got, err := person.ParsePersonName()
	if err != nil {
		t.Fatalf("bad parse err: %v", err)
	}

	if d := cmp.Diff(namefield.ParsedPersonName{PersonID: "foo"}, got); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}
}

func TestReference_ParseFieldMethod(t *testing.T) {
	foo := reference.Foo{Name: "foos/x", Bar: "bars/y"}
	got1, err := foo.ParseBar()
	if err != nil {
		t.Fatalf("bad Parse err: %v", err)
	}

	if d := cmp.Diff(reference.ParsedBarName{BarID: "y"}, got1); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}

	bar := reference.Bar{Name: "bars/z", Foo: "foos/w"}
	got2, err := bar.ParseFoo()
	if err != nil {
		t.Fatalf("bad Parse err: %v", err)
	}

	if d := cmp.Diff(reference.ParsedFooName{FooID: "w"}, got2); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}

	xpkg := reference.CrossPackage{ThingName: "things/foo"}
	got3, err := xpkg.ParseThingName()
	if err != nil {
		t.Fatalf("bad Parse err: %v", err)
	}

	if d := cmp.Diff(simple.ParsedThingName{ThingID: "foo"}, got3); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}

	xpkgExt := reference.CrossPackageExternal{ExternalName: "external/foo"}
	got4, err := xpkgExt.ParseExternalName()
	if err != nil {
		t.Fatalf("bad Parse err: %v", err)
	}

	if d := cmp.Diff(external.ParsedExternalName{ExternalID: "foo"}, got4); d != "" {
		t.Fatalf("bad Parse result:\n%s", d)
	}
}

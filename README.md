# protoc-gen-go-resource

`protoc-gen-go-resource` is a `protoc` plugin that generates convenience
functions for Protobuf messages annotated
with [`google.api.resource`](https://github.com/googleapis/api-common-protos/blob/master/google/api/resource.proto)
.

For example, this example message, taken
from [Google's AIP 123 ("Resource Names")](https://google.aip.dev/123):

```protobuf
syntax = "proto3";

import "google/api/resource.proto";

// A representation of a Pub/Sub topic.
message Topic {
  option (google.api.resource) = {
    type: "pubsub.googleapis.com/Topic"
    pattern: "projects/{project}/topics/{topic}"
  };

  // The resource name of the topic.
  string name = 1;

  // Other fields...
}
```

Would have, in addition to the types usually generated for `Topic`, the
following new functionality:

```go
package pubsub

type ParsedTopicName struct {
	ProjectID string
	TopicID   string
}

func ParseTopicName(s string) (ParsedTopicName, error) {
	// Parses names like "projects/foo/topics/bar" ...
}

func ParseFullTopicName(s string) (ParsedTopicName, error) {
	// Parses names like "//pubsub.googleapis.com/projects/foo/topics/bar" ...
}

func (n ParsedTopicName) Name() string {
	// Returns names like "projects/foo/topics/bar" ...
}

func (n ParsedTopicName) FullName() string {
	// Returns names like "//pubsub.googleapis.com/projects/foo/topics/bar" ...
}
```

You'll also get some new methods on the `Topic` type that the standard Go protoc
generator produces:

```go
package pubsub

func (t *Topic) ParseName() (ParsedTopicName, error) {
	return ParseTopicName(t.Name)
}

func (t *Topic) ParseFullName() (ParsedTopicName, error) {
	return ParseFullTopicName(t.Name)
}
```

This package also has support for
[`google.api.resource_reference`](https://google.aip.dev/122#fields-representing-another-resource),
multi-pattern resources, "externally-defined" resources, and more. However,
there are some differences between this package's support for resources and
Google's usage of those patterns. See ["Differences from AIP /
googleapis"](#differences-from-aip--googleapis) for more on this.

## Supported Functionality

This section details the features supported in this package. Note that all the
features below can be used in concert, but for ease of understanding they are
described independently of one another.

### Basic `google.api.resource` usage

As described at the beginning of this README, you can annotate a resource like
so:

```protobuf
message Thing {
  option (google.api.resource) = {
    type: "example.com/Thing"
    pattern: "things/{thing}"
  };
  
  string name = 1;
}
```

This will generate a struct for representing a parsed `Thing.name` called
`ParsedThingName`. You parse strings like `things/foo` into a `ParsedThingName`
using the generated `ParseThingName` function, and you can generate strings like
`things/foo` from a `ParsedThingName` using the `Name` method on
`ParsedThingName`.

As a convenience method, you'll also get a new `ParseName` method added to your
protoc-generated `Thing` type. This method is essentially just a convenience
wrapper around calling `ParseThingName(thing.Name)`.

### Full resource names

As detailed in [AIP-122](https://google.aip.dev/122#full-resource-names),
resources have "full resource names" in addition to their ordinary "names". The
full resource name of a type is formed by combining the service name of the
resource with its (non-full) name. 

For instance, to take the `Thing` example above, the service name is
`example.com` (the part before the slash in `example.com/Thing`). So this would
be a `Thing` name:

```text
things/foo
```

And this would be a `Thing` *full* name:

```text
//example.com/things/foo
```

In the generated code, you can parse full names using the generated
`ParseFullXXXName` (e.g. `ParseFullThingName` in this example). You can generate
full names from a parsed name using the `FullName` method on the `ParsedXXXName`
type (e.g. `ParsedThingName`).

### Multi-pattern resources

Resources can have multiple patterns. For example, if a `Thing` can have names
like:

```text
things/foo
projects/bar/things/baz
```

Then you can express that by passing `pattern` multiple times to the
`google.api.resource` annotation like so:

```protobuf
message Thing {
  option (google.api.resource) = {
    type: "example.com/Thing"
    pattern: "things/{thing}"
    pattern: "projects/{project}/things/{thing}"
  };
  
  string name = 1;
}
```

When generating code from a multi-pattern resource, this package will produce an
interface instead of a struct for the `ParsedXXXName` (e.g. `ParsedThingName`)
type. You can still call `ParseXXXName` / `ParseFullXXXName` (e.g.
`ParseThingName` / `ParseFullThingName`) and the `Name` / `FullName` methods on
the result, but you will need to use a type switch to handle the different
possible return results.

For each possible pattern, you'll have a different concrete struct that
represents that pattern. In the example above, you'll have the following
structs:

```go
type ParsedThingName_0 struct {
	ThingID string
}

type ParsedThingName_1 struct {
	ProjectID string
	ThingID   string
}
```

Each of those structs will implement the `ParsedThingName` interface, so you can
write code like this to handle different possibilities:

```go
n, err := examplepb.ParseThingName("projects/bar/things/baz")
if err != nil {
	return err
}

switch n := n.(type) {
case examplepb.ParsedThingName_0:
	fmt.Println("you gave me a project-less thing with id", n.ThingID)
case examplepb.ParsedThingName_1:
	fmt.Println("you gave me thing with id", n.ThingID, "in project", n.ProjectID)
default:
	// this code will only run if you generated code from a resource with 2+ patterns, 
	// but haven't yet updated your code to handle ParsedThingName_2 or beyond.
	//
	// ...
}
```

Because the generated interface looks like this:

```go
type ParsedThingName interface {
	Name() string
	FullName() string
	
	// this will be hidden from godoc, but under the hood there's a third method
	mustEmbedParsedThingName()
}
```

You can always call `Name` or `FullName` on the result of parsing a
multi-pattern name. However, you cannot implement the interface yourself,
because there is an unexported method on the interface. This is the same code
generation technique used by the Go code generator for generating `oneof`
implementations.

### Alternative name fields

By default, whenever you add a `google.api.resource` declaration, it's implied
that the `pattern` values you supply refer to a protobuf field called `name` on
the message. If you want to use a different field, you can specify a
`name_field` on your `google.api.resource` call:

```protobuf
message Thing {
  option (google.api.resource) = {
    type: "example.com/Thing"
    pattern: "things/{thing}"
    name_field: "my_cool_name"
  };
  
  string my_cool_name = 1;
}
```

Specifying a non-default `name_field` does not affect the names of the types
generated by this package. However, it does affect the name of the generated
method added to `Thing`. Instead of generating a method called `ParseName`, this
package will instead generate a method called `ParseMyCoolName`.

### `google.api.resource_definition`

If you want to define a resource that doesn't correspond to a protobuf message,
you can add the `google.api.resource_definition` annotation to a proto file,
like so:

```protobuf
option (google.api.resource_definition) = {
  type: "example.com/External"
  pattern: "external/{external}"
};
```

This will generate the usual parser functions and types for a resource named
`External`, namely:

```go
type ParsedExternalName struct {
	ExternalID string
}

func ParseExternalName(s string) (ParsedExternalName, error) {
	// ...
}

func ParseFullExternalName(s string) (ParsedExternalName, error) {
	// ...
}

func (n ParsedExternalName) Name() string {
	// ...
}

func (n ParsedExternalName) FullName() string {
	// ...
}
```

However, since there is no Golang type associated with the resource, no
convenience method for parsing a name can be generated.

You can specify `google.api.resource_definition` multiple times in a single
`.proto` file. A `google.api.resource_definition` can have multiple patterns.
The `name_field` of a `google.api.resource_definition` has no impact on the
generated code.

### `google.api.resource_reference`

You can indicate that a field contains a name of a resource using the
`google.api.resource_reference` field annotation. When you do this, this package
will generate a convenience method for parsing that field.

For example:

```protobuf
message Person {
  string favorite_thing = 1 [(google.api.resource_reference) = {
    type: "example.com/Thing"
  }];
}

message Thing {
  option (google.api.resource) = {
    type: "example.com/Thing"
    pattern: "things/{thing}"
  };

  string name = 1;
}
```

Will generate, in addition to the usual set of generated code for `Thing`, a new method on the `Person`:

```go
func (p *Person) ParseFavoriteThing() (ParsedThingName, error) {
	return ParseThingName(p.FavoriteThing)
}
```

This package determines which type and parsing function to use thanks to the
`type` you provide to `google.api.resource_reference`, and matching it against a
`google.api.resource` or `google.api.resource_definition` declaration.

## Differences from AIP / googleapis

The `google.api.resource`, `google.api.resource_definition`, and
`google.api.resource_reference` annotations are loosely specified in three
places:

* Google's [API Improvement Proposals](https://google.aip.dev/), especially
  [AIP-122](https://google.aip.dev/122), [AIP-123](https://google.aip.dev/123),
  [AIP-180](https://google.aip.dev/180#changing-resource-names),
  and [AIP-4231](https://google.aip.dev/client-libraries/4231), describe how the
  annotations are meant to be used.
* Google's [api-common-protos](https://github.com/googleapis/api-common-protos)
  contains declarations of these annotations that can be used by `protoc` and
  imported by `.proto` files.
* Google's [googleapis](https://github.com/googleapis/googleapis) contains the
  interface definitions of many of Google's public APIs. These provide real-world
  examples of the annotations in use.
  
Based on the information from these sources, there are known differences between
how this package interprets the annotations versus how Google appears to
interpret them internally.

### No support for tilde separators

Google supports use of tilde as separators within a path segment. For instance,
consider [`AdGroupAd` in
googleapis](https://github.com/googleapis/googleapis/blob/b64d30b2b245cb805d1aa328e50258fdff6e628e/google/ads/googleads/v6/resources/ad_group_ad.proto#L44):

```protobuf
message AdGroupAd {
  option (google.api.resource) = {
    type: "googleads.googleapis.com/AdGroupAd"
    pattern: "customers/{customer_id}/adGroupAds/{ad_group_id}~{ad_id}"
  };

  // ...
}
```

This package does not support such patterns, but such support may be added at a
later date. This functionality is not yet implemented simply because of the
complexity required to generate code for parsing such patterns, not for any deep
technical reason.

### Name field must exist

Google's use of the resource annotation does not appear to require that that
name field, `name` by default, does not need to exist. This package considers
the name field not existing to be an error.

`AdGroupAd` in the previous section is an example of a message that specifies
`name` as the name field, but doesn't have a `name` field.

### No support for `child_type`

[AIP-122](https://google.aip.dev/122#fields-representing-a-resources-parent)
documents that, for resources which may have multiple different "parent" types,
`google.api.resource_reference` may use `child_type` instead of `type` to
indicate that the referred-to type is any one of the types which may be a parent
of the given `child_type`.

This package does not support `child_type`, because such support would require
that this package be aware of the parent-child relationships of all resources.
This information is not explicitly provided in the annotations this package
understands. Instead, parent-child relationships are implicit in the shape of
the patterns of resource names.

Determining parent-child relationships from annotations would require doing some
inference based on the patterns of names for each type. Although this is doable,
for the result to be usable in practice, some additional tooling would be
required to make it clear what parent/child relationships the code generator is
inferring, so that developers can debug what's going wrong if the generated code
isn't as expected.

As a result, `child_type` support may be added in the future, but will not be
included without additional tooling for users to be able to gain insight into
the code generator's inferences.

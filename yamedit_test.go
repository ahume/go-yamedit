package yamedit_test

import (
	"testing"

	"github.com/ahume/go-yamedit"
)

func TestGet(t *testing.T) {
	yaml := `foo:
  bar: you
  anno:
    a.b.c/123: right
  brain:
  - 1
  - two
  - three:
    foo: bar`

	tests := []struct {
		path string
		want string
	}{
		{path: "/foo/bar", want: "you"},
		{path: "/foo/brain/0", want: "1"},
		{path: "/foo/brain/2/foo", want: "bar"},
		{path: "/foo/anno/a.b.c~1123", want: "right"},
	}

	for _, tc := range tests {
		r, err := yamedit.Get([]byte(yaml), tc.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if r != tc.want {
			t.Errorf("did not correct value. want: %s, got: %s", tc.want, r)
		}
	}
}

func TestEdit(t *testing.T) {
	yaml := `foo:
  bar: you
  anno:
    a.b.c/123: right
  data:
    config.json: |
      line1
      line2
  brain:
  - zero
  - one
  - two:
    foo: bar`

	tests := []struct {
		path     string
		newValue string
		want     string
	}{
		{path: "/foo/bar", newValue: "fool", want: `foo:
  bar: fool
  anno:
    a.b.c/123: right
  data:
    config.json: |
      line1
      line2
  brain:
  - zero
  - one
  - two:
    foo: bar`},
		{path: "/foo/brain/1", newValue: "fool", want: `foo:
  bar: you
  anno:
    a.b.c/123: right
  data:
    config.json: |
      line1
      line2
  brain:
  - zero
  - fool
  - two:
    foo: bar`},
		{path: "/foo/brain/2/foo", newValue: "fool", want: `foo:
  bar: you
  anno:
    a.b.c/123: right
  data:
    config.json: |
      line1
      line2
  brain:
  - zero
  - one
  - two:
    foo: fool`},
		{path: "/foo/anno/a.b.c~1123", newValue: "left", want: `foo:
  bar: you
  anno:
    a.b.c/123: left
  data:
    config.json: |
      line1
      line2
  brain:
  - zero
  - one
  - two:
    foo: bar`},
		{path: "/foo/data/config.json", newValue: `>
      {"ok":"then"}`, want: `foo:
  bar: you
  anno:
    a.b.c/123: right
  data:
    config.json: >
      {"ok":"then"}
  brain:
  - zero
  - one
  - two:
    foo: bar`},
	}

	for _, tc := range tests {
		r, err := yamedit.Edit([]byte(yaml), tc.path, tc.newValue)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if string(r) != tc.want {
			t.Errorf("did not correct value. \nWANT:\n-------- \n%s \nGOT:\n-------- \n%s", tc.want, r)
		}
	}
}

# yamedit

![CI](https://github.com/BrandwatchLtd/kuber-webhook/workflows/CI/badge.svg?branch=master&event=push)

Allows yaml files that are typically edited by humans to be updated without losing any formatting or comments.

It is not a fully-fledged YAML parser/writer, and only deals with updating values for existing fields in a YAML file. It cannot add new fields or new sub-trees to the data structure.

## Example


```golang
yaml := []byte(`---
apiVersion: v1
kind: Namespace
metadata:
  name: platform
  labels:
    team: pricing`)

path := "/metadata/labels/team"


updated, err := yamedit.Edit(yaml, path, "checkout")
if err != nil {
  log.Fatal(err)
}

fmt.Println(updated)
// apiVersion: v1
// kind: Namespace
// metadata:
//  name: platform
//  labels:
//    team: checkout
```
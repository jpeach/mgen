# mgen

When you want to know the version of a module you depend on at
runtime, `mgen` will generate a trivial package (see example
[here](./internal/vers/vers.go)) from the `go.mod` file in your project.

Integrate `mgen` into your build with a make rule like this:
```
internal/vers/vers.go: go.mod
	mkdir -p $$(dirname $@)
	go run . --package=vers go.mod > $@
```

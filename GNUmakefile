
internal/vers/vers.go: go.mod
	mkdir -p $$(dirname $@)
	go run . --package=vers go.mod > $@

clean:
	rm -f mgen

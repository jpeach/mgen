package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	"golang.org/x/mod/modfile"
)

type version struct {
	Path string
	Vers string
}

type context struct {
	Package      string
	Requirements []version
}

func usage() int {
	fmt.Printf("Generate Go code to store module versions\n\n")
	fmt.Printf("Usage:\n      mgen [OPTIONS] go.mod\n")
	fmt.Printf("\n")
	fmt.Printf("Options\n")
	pflag.PrintDefaults()
	return 64 // EX_USAGE
}

func parse(path string) (*modfile.File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.ParseLax(path, data, nil)
}

func expand(c *context) ([]byte, error) {
	t, err := template.New("mgen").Parse(`
package {{ .Package }}

   
var _versions = map[string]string {
{{ range .Requirements }} "{{ .Path }}":  "{{ .Vers }}",
{{ end }}
}

func VersionOf(mod string) string {
    return _versions[mod]
}

func Modules() []string {
	return []string{
{{ range .Requirements }} "{{ .Path }}",
{{ end }}
	}
}

`)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	var packageName *string = pflag.String("package", "deps", "Name of the generated package")
	var showHelp *bool = pflag.Bool("help", false, "Show help")

	pflag.Parse()

	if *showHelp {
		usage()
		os.Exit(0)
	}

	if len(pflag.Args()) != 1 {
		usage()
		os.Exit(64) // EX_USAGE
	}

	mfile, err := parse(pflag.Args()[0])
	if err != nil {
		fmt.Printf("%s", err)
	}

	c := context{
		Package: *packageName,
	}

	for _, mod := range mfile.Require {
		c.Requirements = append(c.Requirements,
			version{
				Path: mod.Mod.Path,
				Vers: mod.Mod.Version,
			})
	}

	data, err := expand(&c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	data, err = format.Source(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Printf(string(data))
}

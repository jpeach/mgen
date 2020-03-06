package vers

var _versions = map[string]string{
	"github.com/spf13/pflag": "v1.0.5",
	"golang.org/x/mod":       "v0.2.0",
}

func VersionOf(mod string) string {
	return _versions[mod]
}

func Modules() []string {
	return []string{
		"github.com/spf13/pflag",
		"golang.org/x/mod",
	}
}

package yaag

var Debug = false

type Config struct {
	On bool

	BaseUrls map[string]string

	DocTitle string
	DocPath  string
}

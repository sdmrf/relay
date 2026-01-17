package paths

type Layout string

const (
	SystemLayout   Layout = "system"
	PortableLayout Layout = "portable"
)

type Paths struct {
	InstallDir string
	DataDir    string
	BinDir     string
	ConfigDir  string
	CacheDir   string
}

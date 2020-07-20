package shared

type Arg struct{}

type App struct {
	Name        string
	Icon        []byte
	Description string
	Versions    []Version
}

type Version struct {
	Name        string
	Changelog   string
	ArchiveName string
}

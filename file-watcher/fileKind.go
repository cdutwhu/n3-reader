package filewatcher

type (
	EmFileKind string
)

const (
	UnknownKind EmFileKind = "unknown file kind"
	Resource    EmFileKind = "resource"
	Query       EmFileKind = "query"
	Command     EmFileKind = "command"
)

func (fk EmFileKind) String() string {
	return string(fk)
}

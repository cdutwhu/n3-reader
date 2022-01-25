package filewatcher

import (
	"os"
	"path/filepath"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

type (
	EmFileType string
)

const (
	UnknownType EmFileType = "unknown file type"
	Text        EmFileType = "text"
	Image       EmFileType = "image"
	Audio       EmFileType = "audio"
	Video       EmFileType = "video"
	Executable  EmFileType = "executable?"
	Binary      EmFileType = "binary"
	Deleted     EmFileType = "deleted"
)

func (ft EmFileType) String() string {
	return string(ft)
}

var (
	mContType = map[string]EmFileType{
		"text/plain; charset=utf-8": Text,
		"application/pdf":           Text,
		"application/octet-stream":  Binary,
	}

	mBinType = map[string]EmFileType{
		".rmvb": Video,
		".exe":  Executable,
		"":      Executable,
	}
)

func FileType(file string) EmFileType {

	if !fd.FileExists(file) {
		return Deleted
	}

	// Open File
	f, err := os.Open(file)
	lk.FailOnErr("%v", err)
	defer f.Close()

	// Get the content
	contentType, err := GetFileContentType(f)
	lk.FailOnErr("%v", err)

	if t, ok := mContType[contentType]; ok {
		if t == Binary {
			ext := filepath.Ext(file)
			if t, ok := mBinType[ext]; ok {
				return t
			}
			lk.Log("New Binary Type@ %v", ext)
		}
		return t
	}

	lk.Warn("New Type@ [%v], please add it to 'fileKindType.go'", contentType)
	return UnknownType
}
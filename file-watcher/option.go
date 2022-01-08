package filewatcher

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	gio "github.com/digisan/gotk/io"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

const (
	dfltWatched  = "./sample-watched"
	dfltInterval = "10s"
)

var selfId = SelfMD5()

type Option func(*Watcher) error

func (w *Watcher) setOption(options ...Option) error {
	for _, opt := range options {
		if err := opt(w); err != nil {
			return err
		}
	}
	return nil
}

func OptID(id string) Option {
	return func(w *Watcher) error {
		if id != "" {
			w.id = id
			return nil
		}
		w.id = selfId
		return nil
	}
}

func OptName(name string) Option {
	return func(w *Watcher) error {
		if name != "" {
			w.name = name
			return nil
		}
		name, err := os.Hostname()
		if err != nil {
			return err
		}
		w.name = fmt.Sprintf("%s-reader-%s", name, selfId[:4])
		return nil
	}
}

// const (
// 	UnknownKind EmFileKind = "unknown file kind"
// 	Resource    EmFileKind = "resource"
// 	Query       EmFileKind = "query"
// 	Command     EmFileKind = "command"
// )
func OptKind(kind string) Option {
	return func(w *Watcher) error {
		if kind != "" {
			w.fileKind = EmFileKind(kind)
			return nil
		}
		w.fileKind = Resource
		return nil
	}
}

// "|" as separator, like "json|xml|csv"
func OptFormat(format string) Option {
	return func(w *Watcher) error {
		if format == "" {
			return errors.New("input format cannot be empty")
		}

		var fmts []string
		fmtsPile := strings.Split(format, "|")
		fmtsComma := strings.Split(format, ",")
		switch {
		case len(fmtsPile) > 1: // "|" to separate if could be 2+ formats
			fmts = fmtsPile
		case len(fmtsComma) > 1: // "," to separate if could be 2+ formats
			fmts = fmtsComma
		default:
			fmts = append(fmts, format) // only one format
		}
		for i := 0; i < len(fmts); i++ {
			fmts[i] = strings.ToLower(fmts[i])
			fmts[i] = strings.Trim(fmts[i], ".") // remove any excess . chars
		}

		w.format = fmts
		return nil
	}
}

func OptWatcher(folder string, suffix string, interval string, recursive bool, inclHidden bool, ignore string, autodel bool) Option {
	return func(w *Watcher) error {

		w.watcher = watcher.New()

		// dot file handling
		w.watcher.IgnoreHiddenFiles(!inclHidden)
		w.inclHidden = inclHidden

		// If no files/folders were specified, watch the current directory.
		if folder == "" {
			var osErr error
			folder, osErr = os.Getwd()
			folder = filepath.Join(folder, dfltWatched)
			if osErr != nil {
				return errors.Wrap(osErr, "no watch folder specified, and cannot determine current working directory")
			}
		}

		// must create folder if it does not exist. otherwise, panic
		gio.MustCreateDir(folder)
		w.folder = folder

		// Get any of the paths to ignore.
		ignoredPaths := strings.Split(ignore, ",")
		for _, path := range ignoredPaths {
			trimmed := strings.TrimSpace(path)
			if trimmed == "" {
				continue
			}
			if err := w.watcher.Ignore(trimmed); err != nil {
				return errors.Wrap(err, "unable to add ignore folder "+trimmed)
			}
		}
		w.ignore = ignore

		// Only files that match the regular expression for file suffix during file listings will be watched.
		if suffix != "" {
			trimSuffix := strings.Trim(suffix, ".")
			r := regexp.MustCompile("([^\\s]+(\\.(?i)(" + trimSuffix + "))$)")
			w.watcher.AddFilterHook(watcher.RegexFilterHook(r, false))
		}
		w.fileExt = suffix

		// Add the watch folder specified.
		if recursive {
			if err := w.watcher.AddRecursive(folder); err != nil {
				return errors.Wrap(err, "unable to add watch folder "+folder+" recursively")
			}
		} else {
			if err := w.watcher.Add(folder); err != nil {
				return errors.Wrap(err, "unable to add watch folder "+folder)
			}
		}
		w.recursive = recursive

		// Parse the interval string into a time.Duration.
		if interval == "" {
			interval = dfltInterval
		}
		parsedInterval, err := time.ParseDuration(interval)
		if err != nil {
			return errors.Wrap(err, "unable to parse watcher interval as duration")
		}
		w.interval = parsedInterval

		// Set Auto-Delete after Creating
		w.autodel = autodel

		return nil
	}
}

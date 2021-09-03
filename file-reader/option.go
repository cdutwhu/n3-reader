package filereader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	goio "github.com/digisan/gotk/io"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

const (
	dfltWatched  = "./watched"
	dfltInterval = "10s"
)

var selfId = SelfMD5()

type Option func(*Reader) error

func (fr *Reader) setOption(options ...Option) error {
	for _, opt := range options {
		if err := opt(fr); err != nil {
			return err
		}
	}
	return nil
}

func OptID(id string) Option {
	return func(fr *Reader) error {
		if id != "" {
			fr.id = id
			return nil
		}
		fr.id = selfId
		return nil
	}
}

func OptName(name string) Option {
	return func(fr *Reader) error {
		if name != "" {
			fr.name = name
			return nil
		}
		name, err := os.Hostname()
		if err != nil {
			return err
		}
		fr.name = fmt.Sprintf("%s-reader-%s", name, selfId[:4])
		return nil
	}
}

func OptFormat(format string) Option {
	return func(fr *Reader) error {
		if format == "" {
			return errors.New("input format cannot be empty")
		}

		format = strings.ToLower(format)
		format = strings.Trim(format, ".") // remove any excess . chars
		switch format {
		case "csv", "json":
			fr.format = format
			return nil
		default:
			return fmt.Errorf("input format [%s] not supported (must be one of csv|json)", format)
		}
	}
}

func OptWatcher(folder string, fileSuffix string, interval string, recursive bool, inclHidden bool, ignore string) Option {
	return func(fr *Reader) error {

		fr.watcher = watcher.New()

		// dot file handling
		fr.watcher.IgnoreHiddenFiles(!inclHidden)
		fr.inclHidden = inclHidden

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
		goio.MustCreateDir(folder)
		fr.watchFolder = folder

		// Get any of the paths to ignore.
		ignoredPaths := strings.Split(ignore, ",")
		for _, path := range ignoredPaths {
			trimmed := strings.TrimSpace(path)
			if trimmed == "" {
				continue
			}
			err := fr.watcher.Ignore(trimmed)
			if err != nil {
				return errors.Wrap(err, "unable to add ignore folder "+trimmed)
			}
		}
		fr.ignore = ignore

		// Only files that match the regular expression for file suffix during file listings
		// will be watched.
		if fileSuffix != "" {
			trimSuffix := strings.Trim(fileSuffix, ".")
			r := regexp.MustCompile("([^\\s]+(\\.(?i)(" + trimSuffix + "))$)")
			fr.watcher.AddFilterHook(watcher.RegexFilterHook(r, false))
		}
		fr.watchFileExt = fileSuffix

		// Add the watch folder specified.
		if recursive {
			if err := fr.watcher.AddRecursive(folder); err != nil {
				return errors.Wrap(err, "unable to add watch folder "+folder+" recursively")
			}
		} else {
			if err := fr.watcher.Add(folder); err != nil {
				return errors.Wrap(err, "unable to add watch folder "+folder)
			}
		}
		fr.recursive = recursive

		// Parse the interval string into a time.Duration.
		if interval == "" {
			interval = dfltInterval
		}
		parsedInterval, err := time.ParseDuration(interval)
		if err != nil {
			return errors.Wrap(err, "unable to parse watcher interval as duration")
		}
		fr.interval = parsedInterval

		return nil
	}
}

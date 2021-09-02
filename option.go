package n3reader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	goio "github.com/digisan/gotk/io"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

type Option func(*FileReader) error

func (fr *FileReader) setOption(options ...Option) error {
	for _, opt := range options {
		if err := opt(fr); err != nil {
			return err
		}
	}
	return nil
}

func ID(id string) Option {
	return func(fr *FileReader) error {
		if id != "" {
			fr.id = id
			return nil
		}
		fr.id = uuid.New().String()
		return nil
	}
}

func Name(name string) Option {
	return func(fr *FileReader) error {
		if name != "" {
			fr.name = name
			return nil
		}
		fr.name = "file-reader"
		return nil
	}
}

func Format(format string) Option {
	return func(fr *FileReader) error {
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

func Watcher(folder string, fileSuffix string, interval string, recursive bool, inclHidden bool, ignore string) Option {
	return func(fr *FileReader) error {

		fr.watcher = watcher.New()

		// dot file handling
		fr.watcher.IgnoreHiddenFiles(!inclHidden)
		fr.inclHidden = inclHidden

		// If no files/folders were specified, watch the current directory.
		if folder == "" {
			var osErr error
			folder, osErr = os.Getwd()
			folder = filepath.Join(folder, "watched")
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
			interval = "10m"
		}
		parsedInterval, err := time.ParseDuration(interval)
		if err != nil {
			return errors.Wrap(err, "unable to parse watcher interval as duration")
		}
		fr.interval = parsedInterval

		return nil
	}
}

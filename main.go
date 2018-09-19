package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	var inputFile string
	var xxxTags string
	var inputPath string
	var areas []textArea
	flag.StringVar(&inputFile, "input", "", "path to input file")
	flag.StringVar(&xxxTags, "XXX_skip", "", "skip tags to inject on XXX fields")
	flag.StringVar(&inputPath, "path", "", "path to parse all *.pb.go files")

	flag.Parse()

	var xxxSkipSlice []string
	if len(xxxTags) > 0 {
		xxxSkipSlice = strings.Split(xxxTags, ",")
	}

	if len(inputFile) > 0 {
		areas, err = parseFile(inputFile, xxxSkipSlice)
		if err != nil {
			err = errors.Wrapf(err, "failed to parse %s", inputFile)
			return
		}
		err = writeFile(inputFile, areas)
		if err != nil {
			err = errors.Wrapf(err, "failed to write to %s", inputFile)
			return
		}
		return
	}
	if len(inputPath) > 0 {
		var files []os.FileInfo
		var path string
		files, err = ioutil.ReadDir(inputPath)
		if err != nil {
			err = errors.Wrapf(err, "directory not found: %s", inputPath)
			return
		}
		for _, f := range files {
			fileName := f.Name()
			if len(fileName) < 6 {
				continue
			}
			if strings.ToLower(fileName[len(fileName)-6:]) != ".pb.go" {
				continue
			}

			path = fmt.Sprintf("%s/%s", inputPath, f.Name())
			areas, err = parseFile(path, xxxSkipSlice)
			if err != nil {
				err = errors.Wrapf(err, "failed to parse %s", path)
				return
			}
			if len(areas) < 1 {
				continue
			}
			err = writeFile(path, areas)
			if err != nil {
				err = errors.Wrapf(err, "failed to write to %s", path)
				return
			}
		}
		return
	}
	err = fmt.Errorf("no valid arguments passed")
	return
}

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	download "github.com/joeybloggs/go-download"
	"github.com/vbauerster/mpb"
)

func main() {
	if len(os.Args) != 2 || strings.HasPrefix(os.Args[1], "-") {
		println("Usage:\n\tsuperdownload [url]")
		os.Exit(111)
	}
	progress := mpb.New().SetWidth(80)
	downloadurl := os.Args[1]
	defer progress.Stop()
	// add progress bars
	options := &download.Options{
		Proxy: func(name string, size int64, r io.Reader) io.Reader {
			bar := progress.AddBar(size).
				PrependName(name, 0, 0).
				PrependCounters("%3s / %3s", mpb.UnitBytes, 18, mpb.DwidthSync|mpb.DextraSpace).
				AppendPercentage(5, 0)

			return bar.ProxyReader(r)
		},
	}
	// do it
	f, err := download.Open(downloadurl, options)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// read the whole file at once
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Base(downloadurl), b, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(111)

	}
}

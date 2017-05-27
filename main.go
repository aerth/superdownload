package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/proxy"

	download "github.com/joeybloggs/go-download"
	"github.com/vbauerster/mpb"
)

var (
	proxypath = flag.String("socks", "", "format: socks5://127.0.0.1:1080")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		println("Usage:\n\tsuperdownload [flags] [url]")
		println("Flags:")
		flag.PrintDefaults()
		os.Exit(111)
	}

	proxyurl, err := url.Parse(*proxypath)
	if err != nil {
		println("error parsing SOCKS:", err.Error())
		os.Exit(111)
	}

	progress := mpb.New().SetWidth(80)
	downloadurl := os.Args[1]
	defer progress.Stop()
	// add progress bars
	options := &download.Options{
		Client: func() http.Client {
			client := http.DefaultClient
			if *proxypath != "" {
				dialer, err := proxy.FromURL(proxyurl, proxy.Direct)
				if err != nil {
					panic(err)
				}
				transport := &http.Transport{Dial: dialer.Dial}
				client.Transport = transport
			}
			return *client
		},
		Proxy: download.ProxyFn(func(name string, download int, size int64, r io.Reader) io.Reader {
			bar := progress.AddBar(size).
				PrependName(name, 0, 0).
				PrependCounters("%3s / %3s", mpb.UnitBytes, 18, mpb.DwidthSync|mpb.DextraSpace).
				AppendPercentage(5, 0)

			return bar.ProxyReader(r)
		}),
	}
	// do it

	fi, err := os.Create(filepath.Base(downloadurl))
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	f, err := download.Open(downloadurl, options)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	defer fi.Close()
	_, err = io.Copy(fi, f)
	if err != nil {
		log.Fatal(err)
	}
	println("download completed successfully")

}

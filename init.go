package main

import "flag"

var flags struct {
	Serve         bool
	NoGzip        bool
	GenerateSpecs bool
	PrettyPrint   bool
	RootJson      string
	ListenOn      string
}

func init() {
	flag.BoolVar(&flags.NoGzip, "no-gzip", false, "Do not send Accept-Encoding to upstream")
	flag.BoolVar(&flags.GenerateSpecs, "generate-specs", false, "generate api specs and exit")
	flag.BoolVar(&flags.Serve, "serve", false, "Start the proxy")
	flag.BoolVar(&flags.PrettyPrint, "pretty-print", false, "Pritty print the resulted json works with -generate-specs")
	flag.StringVar(&flags.RootJson, "root-json", "./definitions/api.json", "where are your definitions")
	flag.StringVar(&flags.ListenOn, "listen", "127.0.0.1:80", "Listen on ")

	flag.Parse()
}

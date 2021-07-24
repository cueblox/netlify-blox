package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/apex/gateway"
	"github.com/cueblox/blox/content"
)

func main() {
	port := flag.Int("port", -1, "specify a port to use http rather than AWS Lambda")
	flag.Parse()
	listener := gateway.ListenAndServe
	portStr := ""

	userConfig, err := ioutil.ReadFile("blox.cue")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := content.NewService(string(userConfig), true)
	if err != nil {
		log.Fatal(err)
	}

	hf, err := repo.GQLHandlerFunc()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", hf)

	h, err := repo.GQLPlaygroundHandler()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/ui", h)

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		staticDir, err := repo.Cfg.GetString("static_dir")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(".", staticDir)))))

	}

	log.Fatal(listener(portStr, nil))
}

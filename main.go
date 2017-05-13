package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"

	"github.com/orijtech/otils"
)

func main() {
	var http1 bool
	var http1Port int
	flag.BoolVar(&http1, "http1", false, "run the website in http1 mode")
	flag.IntVar(&http1Port, "http1-port", 80, "the port to run the website on in HTTP1 mode")
	flag.Parse()

	http1Addr := fmt.Sprintf(":%d", http1Port)
	mux := http.DefaultServeMux
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	if http1 {
		log.Printf("Running in http1 mode on address: %q\n", http1Addr)
		if err := http.ListenAndServe(http1Addr, mux); err != nil {
			log.Fatal(err)
		}
	}

	// Auto redirect all non-https traffic
	go func() {
		log.Printf("Running in http2 mode. Redirecting all non-https traffic to %q", http1Addr)
		redirectURL := otils.FirstNonEmptyString(os.Getenv("COMMAI_REDIRECT_NON_HTTPS_URL"), "https://comma.ai/")
		err := http.ListenAndServe(http1Addr, otils.RedirectAllTrafficTo(redirectURL))
		if err != nil {
			log.Printf("failed to bind the redirector: %v", err)
		}
	}()

	domainsCSV := otils.FirstNonEmptyString(os.Getenv("COMMAI_DOMAINS_CSV"), "comma.ai,www.comma.ai")
	domains := strings.Split(domainsCSV, ",")

	log.Fatal(http.Serve(autocert.NewListener(domains...), mux))
}

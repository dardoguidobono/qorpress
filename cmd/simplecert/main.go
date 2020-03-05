//
//  simplecert
//
//  Created by Philipp Mieden
//  Contact: dreadl0ck@protonmail.ch
//  Copyright © 2018 bestbytes. All rights reserved.
//

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"
)

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello from simplecert"))
}

func main() {

	// do the cert magic
	cfg := simplecert.Default
	cfg.Domains = []string{"x0rzkov.com", "www.x0rzkov.com"}
	cfg.CacheDir = "letsencrypt"
	cfg.SSLEmail = "x0rzkov@protonmail.com"
	// cfg.Local = true
	certReloader, err := simplecert.Init(cfg, func() {
		// this function will be called upon receiving the syscall.SIGINT or syscall.SIGABRT signal
		// and can be used to stop your backend gracefully
		fmt.Println("cleaning up...")
	})
	if err != nil {
		log.Fatal("simplecert init failed: ", err)
	}

	// redirect HTTP to HTTPS
	log.Println("starting HTTP Listener on Port 80")
	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))

	// init strict tlsConfig with certReloader
	tlsconf := tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)

	// now set GetCertificate to the reloaders GetCertificateFunc to enable hot reload
	tlsconf.GetCertificate = certReloader.GetCertificateFunc()

	// init server
	s := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsconf,
		Handler:   Handler{},
	}

	log.Println("now visit: https://" + cfg.Domains[0])

	// lets go
	log.Fatal(s.ListenAndServeTLS("", ""))
}

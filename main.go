package main

import (
	"flag"
	"fmt"
	"gvweb/simplemux"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func accessLogger(l chan string) {
	for {
		log.Print(<-l)
	}
}

type statusWrapper struct {
	status int
	http.ResponseWriter
}

func (s *statusWrapper) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func httpLog(l chan string, sw *statusWrapper, r *http.Request) {

	var remote string
	if len(r.Header["X-Forwarded-For"]) > 0 {
		remote = r.Header["X-Forwarded-For"][0]
	} else {
		remote = r.RemoteAddr
	}
	l <- fmt.Sprintf("%s %s %s %d", remote, r.Method, r.URL, sw.status)
}

func httpWrapper(l chan string, handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWrapper{-1, w}
		w.Header().Set("Server", "Go HTTP handler")
		handler.ServeHTTP(&sw, r)
		httpLog(l, &sw, r)
	})
}

func serveHTTP(port string) {
	accessLogChan := make(chan string, 64)
	go accessLogger(accessLogChan)

	reHandler := simplemux.NewRegexpHandler()
	reHandler.AddRoute("^/$", "GET", homePage)
	reHandler.AddRoute("^/generate$", "POST", generateHandler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/"+gDataDir, http.FileServer(http.Dir(".")))
	http.Handle("/", reHandler)

	var scheme string
	if *gUseTLS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	log.Printf("gvweb(%s) listening at %s port %s\n", gVersion, scheme, port)

	var err error
	srv := http.Server{
		Addr:         ":" + port,
		Handler:      httpWrapper(accessLogChan, http.DefaultServeMux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if *gUseTLS {
		err = srv.ListenAndServeTLS(*gTLSCert, *gTLSKey)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

var gPort = flag.Int("port", 12345, "port number to listen on")
var gUseTLS = flag.Bool("usetls", false, "Use TLS(HTTPS) intead of plain HTTP")
var gTLSCert = flag.String("tlscert", "tls.cert", "Path to TLS certificate file")
var gTLSKey = flag.String("tlskey", "tls.key", "Path to TLS key file")
var gCleanupInterval = flag.Int("purge", 24*60*60, "Remove saved graphs that are older than this amount in seconds. 0 to keep them forever.")
var gVersion = "DEVELOPMENT"

const gDataDir = "data/"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s v.%s\n", os.Args[0], gVersion)
	}

	flag.Parse()

	runtime.GOMAXPROCS(3)

	if *gCleanupInterval > 0 {
		interval := time.Duration(*gCleanupInterval) * time.Second
		initPurge(gDataDir, interval)
		log.Printf("Purging data older than %v\n", interval)
	}

	serveHTTP(fmt.Sprintf("%d", *gPort))
}

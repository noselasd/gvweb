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

func httpLog(l chan string, w http.ResponseWriter, r *http.Request) {

	var remote string
	if len(r.Header["X-Forwarded-For"]) > 0 {
		remote = r.Header["X-Forwarded-For"][0]
	} else {
		remote = r.RemoteAddr
	}
	l <- fmt.Sprintf("%s %s %s", remote, r.Method, r.URL)
}

func httpWrapper(l chan string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Go HTTP handler")
		handler.ServeHTTP(w, r)
		httpLog(l, w, r)
	})
}

func serveHTTP(port string) {
	accessLogChan := make(chan string, 64)
	go accessLogger(accessLogChan)

	reHandler := simplemux.NewRegexpHandler()
	reHandler.AddRoute("^/$", "GET", homePage)
	reHandler.AddRoute("^/generate$", "POST", generateHandler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/"+g_DataDir, http.FileServer(http.Dir(".")))
	http.Handle("/", reHandler)

	var scheme string
	if *g_UseTLS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	log.Printf("gvweb(%s) listening at %s port %s\n", g_Version, scheme, port)

	var err error

	if *g_UseTLS {
		err = http.ListenAndServeTLS(":"+port, *g_TLSCert, *g_TLSKey, httpWrapper(accessLogChan, http.DefaultServeMux))
	} else {
		err = http.ListenAndServe(":"+port, httpWrapper(accessLogChan, http.DefaultServeMux))
	}
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

var g_Port = flag.Int("port", 12345, "port number to listen on")
var g_UseTLS = flag.Bool("usetls", false, "Use TLS(HTTPS) intead of plain HTTP")
var g_TLSCert = flag.String("tlscert", "tls.cert", "Path to TLS certificate file")
var g_TLSKey = flag.String("tlskey", "tls.key", "Path to TLS key file")
var g_CleanupInterval = flag.Int("purge", 24*60*60, "Remove saved graphs that are older than this amount in seconds. 0 to keep them forever.")
var g_Version = "DEVELOPMENT"

const g_DataDir = "data/"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s v.%s\n", os.Args[0], g_Version)
	}

	flag.Parse()

	runtime.GOMAXPROCS(3)

	if *g_CleanupInterval > 0 {
		interval := time.Duration(*g_CleanupInterval) * time.Second
		initPurge(g_DataDir, interval)
		log.Printf("Purging data older than %v\n", interval)
	}

	serveHTTP(fmt.Sprintf("%d", *g_Port))
}

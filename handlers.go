package main

import (
	"fmt"
	"log"
	"net/http"
)

var g_supportedTools = [...]string{"dot", "neato", "twopi", "circo", "fdp", "sfdp", "patchwork"}

func homePage(w http.ResponseWriter, r *http.Request, matches []string) {
	http.ServeFile(w, r, "static/index.html")
}

func isSupportedTool(tool string) bool {
	for _, val := range g_supportedTools {
		if tool == val {
			return true
		}
	}
	return false
}

func generateHandler(w http.ResponseWriter, r *http.Request, matches []string) {
	graph := r.FormValue("graphtext")
	imgType := r.FormValue("imagetype")
	tool := r.FormValue("tool")

	if !isSupportedTool(tool) {
		http.Error(w, fmt.Sprintf("Tool '%s' is not supported", tool), http.StatusBadRequest)
		return
	}
	if len(imgType) == 0 {
		http.Error(w, fmt.Sprintf("imagetype is not specified"), http.StatusBadRequest)
		return
	}

	if len(graph) == 0 {
		http.Error(w, "Empty input", http.StatusBadRequest)
		return
	}

	result := runGraphviz(tool, graph, imgType)
	if result.err != nil {
		log.Print(result.err)
		http.Error(w, result.err.Error(), http.StatusNotAcceptable)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(result.fileName))
		return
	}
}

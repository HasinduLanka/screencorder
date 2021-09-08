package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

var API_POSTs map[string]API_POST_Recieved = map[string]API_POST_Recieved{
	"api/recchunk": RecChunkRecieved,
	"api/final":    FinalRecieved,
	"api/end":      EndRecieved,
}

var API_GETs map[string]API_GET_Recieved = map[string]API_GET_Recieved{
	"api/start":     StartRec,
	"api/handshake": Handshake,
}

func ServeFull(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "404 not found.", http.StatusNotFound)
	// 	return
	// }

	urlpath := strings.ReplaceAll(strings.TrimPrefix(r.URL.Path, "/"), "..", "")

	switch r.Method {
	case "GET":

		for apiPath, api := range API_GETs {
			if strings.HasPrefix(urlpath, apiPath) {
				resp := api(strings.TrimPrefix(urlpath, apiPath+"/"))
				if resp.headers != nil {
					for k, v := range resp.headers {
						w.Header().Set(k, v)
					}
				}
				if len(resp.body) > 0 {
					w.Write(resp.body)
				}
				return
			}
		}

		// If no API
		fmt.Println("Serving " + urlpath)

		http.ServeFile(w, r, path.Join(rootDir, urlpath))

	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		println("POST " + urlpath + " body length " + strconv.Itoa(len(body)))

		for apiPath, api := range API_POSTs {
			if strings.HasPrefix(urlpath, apiPath) {
				resp := api(strings.TrimPrefix(urlpath, apiPath+"/"), body)

				if resp.headers != nil {
					for k, v := range resp.headers {
						w.Header().Set(k, v)
					}
				}
				if len(resp.body) > 0 {
					w.Write(resp.body)
				}

			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func ServeMirrorAsRoot(w http.ResponseWriter, r *http.Request) {

	urlpath := strings.ReplaceAll(strings.TrimPrefix(r.URL.Path, "/"), "..", "")

	if len(urlpath) == 0 || urlpath == "index.html" || urlpath == "index" {
		http.Redirect(w, r, "/mirror/viewer.html", http.StatusSeeOther)
	} else if strings.HasPrefix(urlpath, "mapi/") || strings.HasPrefix(urlpath, "mirror/") {
		ServeFull(w, r)
	}

}

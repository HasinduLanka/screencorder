package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var wsroot string = "workspace/"
var rootDir string = ""

var SpeakerInputName string = ""
var AudioEnabled bool = false
var SSLEnabled bool = false

var API_POSTs map[string]API_POST_Recieved = map[string]API_POST_Recieved{
	"api/recchunk":    RecChunkRecieved,
	"api/mirrorchunk": MirrorChunkRecieved,
	"api/mirecchunk":  MirrorAndRecChunkRecieved,
	"api/final":       FinalRecieved,
	"api/end":         EndRecieved}

var API_GETs map[string]API_GET_Recieved = map[string]API_GET_Recieved{
	"api/start":     StartRec,
	"api/mstart":    StartRecOnlyMirror,
	"api/handshake": Handshake,
	"api/reqview":   View,
	"mapi/reqview":  View}

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

func main() {
	println("                   .                  ")
	println("--------------------------------------")
	println("             Screencorder             ")
	println("--------------------------------------")
	println("                                      ")
	println("        Free and Open source          ")
	println("         project written by           ")
	println("             Bitblazers               ")
	println("                                      ")
	println("       Github.com/HasinduLanka        ")
	println("       Github.com/Bitblazers-lk       ")
	println("--------------------------------------")
	println("                   .                  ")

	HomeDir, HomeDirErr := os.UserHomeDir()
	if HomeDirErr == nil {
		HomeVideo := path.Join(HomeDir, "Videos/screencorder") + "/"
		if os.MkdirAll(HomeVideo, os.ModePerm) == nil {
			wsroot = HomeVideo
		}
	}
	MakeDir(wsroot)

	if FileExists("index.html") {
		cwd, cwderr := filepath.Abs("")
		if cwderr == nil {
			rootDir = cwd + "/"
		} else {
			rootDir = "/"
		}
	} else {
		ex, err := os.Executable()
		if err == nil {
			resolvedPath, resolvedErr := filepath.EvalSymlinks(ex)
			if resolvedErr == nil {
				rootDir = filepath.Dir(resolvedPath) + "/"
			} else {
				rootDir = ex + "/"
			}
		}
	}

	println("Operating on " + rootDir)

	CheckError(InitExec())

	HiOut, HiErr := ExcecCmd("echo 'System calls working'")
	println(HiOut)
	if HiErr != nil {
		println(HiErr.Error())
	}

	// EndTask := make(chan bool)
	// go ExcecCmdTask("echo BashWorks1 > bashworks ; sleep 2 ; echo BashWorks2 >>  bashworks ; sleep 2 ; echo BashWorks3 >>  bashworks ; sleep 3 ; echo BashWorks4 >>  bashworks", EndTask)

	// time.Sleep(time.Second * 1)
	// EndTask <- false

	if AudioEnabled {
		AudioEnabled = false
		DetectSoundInput()
	}

	MirrorMux := http.NewServeMux()
	FullMux := http.NewServeMux()

	// MirrorMux.HandleFunc("/mirror/", ServeFull)
	// MirrorMux.HandleFunc("/api/", ServeFull)
	MirrorMux.HandleFunc("/", ServeMirrorAsRoot)
	FullMux.HandleFunc("/", ServeFull)

	myip, errIp := GetOutboundIP()

	NetworkEnabled := errIp == nil

	if NetworkEnabled {
		CheckSSL()
		println("My local IP address is " + myip)

	} else {
		SSLEnabled = false
	}

	if SSLEnabled {
		println("Starting Screencorder http://localhost:49542 ")
	} else {
		println("Starting Screencorder https://localhost:49542 ")
	}

	if SSLEnabled {
		println("Connect to the same LAN and visit \n https://" + myip + ":49542   for host interface, \n  http://" + myip + ":49543   for mirror")
	} else {
		println("Visit \n  http://localhost:49542   for host interface, \n  http://localhost:49543   for mirror")
	}

	if SSLEnabled {
		go OpenProgram("xdg-open", "https://localhost:49542")
	} else {
		go OpenProgram("xdg-open", "http://localhost:49542")
	}

	go func() {

		if err := http.ListenAndServe(":49543", MirrorMux); err != nil {
			log.Fatal(err)
		}

	}()

	if SSLEnabled {
		if err := http.ListenAndServeTLS(":49542", wsroot+"server.crt", wsroot+"server.key", FullMux); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(":49542", FullMux); err != nil {
			log.Fatal(err)
		}
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		println("Cannot get my IP address. May be offline")
		// PrintError(err)
		return "", errors.New("error")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func CheckSSL() bool {
	SSLEnabled = false

	if _, errc := os.Stat(wsroot + "server.crt"); errc == nil {
		// server.crt exists
		if _, errk := os.Stat(wsroot + "server.key"); errk == nil {
			// server.key exists
			SSLEnabled = true
		}
	}

	if SSLEnabled {
		println("SSL certicate found. I will be HTTPS now")
	} else {
		println("No SSL certicates found. I will be HTTP. If you want to connect from other hosts, try creating a SSL certificate as following")
		println("go to " + wsroot + " folder \n\t openssl genrsa -out server.key 2048")
		println("\t openssl ecparam -genkey -name secp384r1 -out server.key")
		println("\t openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650\n")
	}
	return SSLEnabled
}

func DetectSoundInput() {

	DSo, DSErr := ExcecProgramToString("pacmd", "list-sinks")
	PrintError(DSErr)
	// println(DSo)

	re, _ := regexp.Compile("name: <(.*)>")
	matches := re.FindAllStringSubmatch(DSo, 32)

	if len(matches) <= 0 {
		println("No sound device found")
		AudioEnabled = false
		return
	}

	SoundInputs := map[string]string{}

	for i, match := range matches {
		if len(match) != 2 {
			println("Sound device issue on " + strings.Join(match, " ; "))
		} else {
			SpeakerInputName = match[1]
			println("Sound device detected " + SpeakerInputName)
			SoundInputs[strconv.Itoa(i)] = SpeakerInputName
		}
	}

	if len(SoundInputs) == 0 {
		AudioEnabled = false
		return
	} else if len(SoundInputs) > 1 {
		ch := PromptOptions("Multiple sound devices detected. Which one to use?", SoundInputs)
		SpeakerInputName = SoundInputs[ch]
	}
	AudioEnabled = true
	println("Selected sound device " + SpeakerInputName)

}

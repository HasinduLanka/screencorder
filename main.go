package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var wsroot string = "workspace/"

var API_POSTs map[string]API_POST_Recieved = map[string]API_POST_Recieved{"chunk": ChunkRecieved, "final": FinalRecieved}
var API_GETs map[string]API_GET_Recieved = map[string]API_GET_Recieved{"start": StartRec, "handshake": Handshake}

func hello(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "404 not found.", http.StatusNotFound)
	// 	return
	// }

	urlpath := strings.ReplaceAll(strings.TrimPrefix(r.URL.Path, "/"), "..", "")

	switch r.Method {
	case "GET":

		for apiPath, api := range API_GETs {
			if strings.HasPrefix(urlpath, apiPath) {
				api(strings.TrimPrefix(urlpath, apiPath+"/"))
				return
			}
		}

		// If no API
		fmt.Println("Serving " + urlpath)
		http.ServeFile(w, r, urlpath)

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
				api(strings.TrimPrefix(urlpath, apiPath+"/"), body)
			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
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

	DetectSoundInput()

	http.HandleFunc("/", hello)

	println("Starting Screencorder http://localhost:49542 ")
	go ExcecProgram("xdg-open", "http://localhost:49542")
	if err := http.ListenAndServe(":49542", nil); err != nil {
		log.Fatal(err)
	}
}

var SpeakerInputName string

func DetectSoundInput() {

	DSo, DSErr := ExcecProgramToString("pacmd", "list-sinks")
	PrintError(DSErr)
	// println(DSo)

	re, _ := regexp.Compile("name: <(.*)>")
	matches := re.FindAllStringSubmatch(DSo, -1)

	if len(matches) <= 0 {
		println("No sound device found")
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
		return
	} else if len(SoundInputs) > 1 {
		ch := PromptOptions("Multiple sound devices detected. Which one to use?", SoundInputs)
		SpeakerInputName = SoundInputs[ch]
	}

	println("Selected sound device " + SpeakerInputName)

}

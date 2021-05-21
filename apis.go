package main

import (
	"regexp"
	"strings"
	"time"
)

// var ifile int = 1

type API_POST_Recieved func(string, []byte) []byte
type API_GET_Recieved func(string) []byte

var AudioTasks map[string]chan bool = map[string]chan bool{}
var FinalizingDone map[string]bool = map[string]bool{}

func ChunkRecieved(path string, chunk []byte) []byte {
	WriteFile(wsroot+path+".webm", chunk)
	return []byte("Host recieved " + path)
}

func FinalRecieved(path string, body []byte) []byte {
	// End Audio recording
	EndTask, found := AudioTasks[path]
	if found {
		EndTask <- false
		delete(AudioTasks, path)
	}

	FinalizingDone[path+".webm"] = false
	println("Finalizing " + path)

	WriteFile(wsroot+path+".fflist", body)

	HiOut, HiErr := ExcecProgram("ffmpeg", "-f", "concat", "-safe", "0", "-i", path+".fflist", "-c", "copy", path+"-video.webm")
	MuxOut, MuxErr := ExcecCmd("ffmpeg -i " + path + "-video.webm -i " + path + ".mp3 -map 0:v -map 1:a -c:v copy -shortest " + path + ".webm")

	println(HiOut)
	PrintError(HiErr)

	println(MuxOut)
	PrintError(MuxErr)

	ExcecCmd("rm -f ./" + path + "-*.webm")
	ExcecCmd("rm -f ./" + path + ".fflist")
	ExcecCmd("rm -f ./" + path + ".mp3")

	delete(FinalizingDone, path+".webm")
	println("Finalized " + path)

	return []byte("Final Recieved " + path)
}

func EndRecieved(path string, body []byte) []byte {
	println("Ending " + path)

	WriteFile(wsroot+path+".end.fflist", body)

	Sbody := string(body)
	re := regexp.MustCompile("file '(.*)'")
	matches := re.FindAllStringSubmatch(Sbody, -1)

	FinalFiles := make([]string, len(matches))

	for i, match := range matches {
		if len(match) == 2 {
			fl := match[1]
			println("Checking " + fl)
			FinalFiles[i] = wsroot + fl

			_, contains := FinalizingDone[fl]
			for contains {
				println("Waiting for " + fl)
				time.Sleep(1 * time.Second)
				_, contains = FinalizingDone[fl]
			}

		} else {
			println("Parse error " + wsroot + strings.Join(match, ", "))
		}
	}

	HiOut, HiErr := ExcecProgram("ffmpeg", "-f", "concat", "-safe", "0", "-i", path+".end.fflist", "-c", "copy", path+".rec.webm")

	println(HiOut)
	PrintError(HiErr)

	for _, match := range FinalFiles {
		println("Delete file " + match)
		DeleteFiles(match)
	}

	// ExcecCmd("rm -f ./" + path + "-*.webm")
	DeleteFiles(wsroot + path + ".end.fflist")
	println("Ended " + path)

	return []byte("End Recieved " + path)
}

func Handshake(path string) []byte {
	ExcecProgram("echo", "Host Ready")
	return []byte("Host Ready")
}

func StartRec(path string) []byte {
	// ExcecProgram("echo", "start recording")
	EndTask := make(chan bool)
	if len(SpeakerInputName) == 0 {
		return []byte("No Audio")
	}
	go ExcecCmdTask("parec -d alsa_output.pci-0000_00_1f.3.analog-stereo.monitor | lame -r -V0 - "+path+".mp3", EndTask)
	AudioTasks[path] = EndTask
	return []byte("Started")
}

package main

import (
	"regexp"
	"strings"
	"time"
)

// var ifile int = 1

type API_POST_Recieved func(string, []byte) Response
type API_GET_Recieved func(string) Response

type Response struct {
	body    []byte
	headers map[string]string
}

func BodyResponse(body []byte) Response {
	return Response{body, map[string]string{}}
}

var AudioTasks map[string]chan bool = map[string]chan bool{}
var FinalizingDone map[string]bool = map[string]bool{}

var ViewerChunk []byte
var ViewerChunkPath string

func MirrorAndRecChunkRecieved(path string, chunk []byte) Response {

	go MirrorAndRecChunkRecievedHandler(path, chunk)

	return BodyResponse([]byte("Host recieved " + path))
}

func MirrorAndRecChunkRecievedHandler(path string, chunk []byte) {
	WriteFile(wsroot+path+".webm", chunk)
	ExcecCmd("ffmpeg -i " + path + ".webm -cpu-used -8 -deadline realtime -c copy " + path + ".m.webm")

	var err error
	ViewerChunk, err = LoadFile(wsroot + path + ".m.webm")
	PrintError(err)
	ViewerChunkPath = path
}

func RecChunkRecieved(path string, chunk []byte) Response {
	go WriteFile(wsroot+path+".webm", chunk)
	return BodyResponse([]byte("Host recieved " + path))
}

func MirrorChunkRecieved(path string, chunk []byte) Response {

	go func() {
		MirrorAndRecChunkRecievedHandler(path, chunk)

		DeleteFiles(wsroot + path + ".webm")
		DeleteFiles(wsroot + path + ".m.webm")
	}()

	return BodyResponse([]byte("Host recieved " + path))
}

func FinalRecieved(path string, body []byte) Response {
	if AudioEnabled {
		// End Audio recording
		EndTask, found := AudioTasks[path]
		if found {
			EndTask <- false
			delete(AudioTasks, path)
		}
	}

	FinalizingDone[path+".webm"] = false
	println("Finalizing " + path)

	WriteFile(wsroot+path+".fflist", body)

	if AudioEnabled {
		HiOut, HiErr := ExcecProgram("ffmpeg", "-f", "concat", "-safe", "0", "-i", path+".fflist", "-c", "copy", path+"-video.mkv")
		MuxOut, MuxErr := ExcecCmd("ffmpeg -i " + path + "-video.mkv -i " + path + ".mp3 -map 0:v -map 1:a -c:v copy -shortest " + path + ".mkv")

		println(HiOut)
		PrintError(HiErr)

		println(MuxOut)
		PrintError(MuxErr)

		go DeleteFiles(wsroot + path + ".mp3")

	} else {
		HiOut, HiErr := ExcecProgram("ffmpeg", "-f", "concat", "-safe", "0", "-i", path+".fflist", "-c", "copy", path+".mkv")

		println(HiOut)
		PrintError(HiErr)
	}

	go DeleteFiles(wsroot + path + "-video.mkv")

	go ExcecCmd("rm -f ./" + path + "-*.webm")
	go DeleteFiles(wsroot + path + ".fflist")

	delete(FinalizingDone, path+".webm")
	println("Finalized " + path)

	return BodyResponse([]byte("Final Recieved " + path))
}

func EndRecieved(path string, body []byte) Response {
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

	return BodyResponse([]byte("End Recieved " + path))
}

func Handshake(path string) Response {
	ExcecProgram("echo", "Host Ready")
	return BodyResponse([]byte("Host Ready"))
}

func StartRec(path string) Response {
	// ExcecProgram("echo", "start recording")
	EndTask := make(chan bool)
	if !AudioEnabled || (len(SpeakerInputName) == 0) {
		AudioEnabled = false
		return BodyResponse([]byte("No Audio"))
	}
	ExcecCmdTask("parec -d "+SpeakerInputName+".monitor | lame -r -V0 - "+path+".mp3", EndTask)
	AudioTasks[path] = EndTask
	return BodyResponse([]byte("Started with audio"))
}

func StartRecOnlyMirror(path string) Response {
	return BodyResponse([]byte("Started only mirror"))
}

func View(path string) Response {
	if len(ViewerChunkPath) == 0 {
		println("Viewer : waiting ")
		return Response{[]byte{}, map[string]string{"cpath": "wait"}}
	}

	if path == "new" {
		println("Viewer : new : " + ViewerChunkPath)
		return Response{ViewerChunk, map[string]string{"cpath": ViewerChunkPath}}
	} else if path == ViewerChunkPath {
		println("Viewer : same")
		return Response{[]byte{}, map[string]string{"cpath": "same"}}
	} else {
		println("Viewer : " + ViewerChunkPath)
		return Response{ViewerChunk, map[string]string{"cpath": ViewerChunkPath}}
	}

}

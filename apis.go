package main

// var ifile int = 1

type API_POST_Recieved func(string, []byte)
type API_GET_Recieved func(string) []byte

var AudioTasks map[string]chan bool = map[string]chan bool{}

func ChunkRecieved(path string, chunk []byte) {
	WriteFile(wsroot+path+".webm", chunk)
}

func FinalRecieved(path string, body []byte) {
	// End Audio recording
	EndTask, found := AudioTasks[path]
	if found {
		EndTask <- false
		delete(AudioTasks, path)
	}

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

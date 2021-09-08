package main

import (
	"errors"
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

var RecievedChunks map[string]string = make(map[string]string, 64)
var AudioTasks map[string]chan bool = map[string]chan bool{}

// var FinalizingDone map[string]bool = map[string]bool{}

func RecChunkRecieved(path string, chunk []byte) Response {
	go func() {

		paths := strings.Split(path, "/")
		if len(paths) < 3 {
			PrintError(errors.New("RecChunkRecieved: Path is not valid : " + path))
			return
		}

		CEncodeType := paths[0] // r : re-encode , c : copy
		ChunkName := paths[1]   // Record-1631091385274-0000000002
		ChunkType := paths[2]   // webm , mp4

		SavedChunkName := ChunkName + "." + ChunkType

		EncodeType := ""
		OutputType := ""

		if NoReEncode {
			CEncodeType = "c"
		}

		switch CEncodeType {
		case "c":
			EncodeType = " -c copy "
			OutputType = DefaultVideoType
		case "r":
			EncodeType = " "
			OutputType = DefaultVideoType
		case "uf":
			EncodeType = DefaultVideoCodec + " -crf 22 -preset ultrafast "
			OutputType = DefaultVideoType
		case "f":
			EncodeType = DefaultVideoCodec + " -crf 20 -preset veryfast "
			OutputType = DefaultVideoType
		case "m":
			EncodeType = DefaultVideoCodec + " -crf 18 -preset veryfast "
			OutputType = DefaultVideoType
		case "q":
			EncodeType = DefaultVideoCodec + " -crf 16 -preset medium "
			OutputType = DefaultVideoType
		case "hq":
			EncodeType = DefaultVideoCodec + " -crf 14 -preset slow "
			OutputType = DefaultVideoType
		case "uq":
			EncodeType = DefaultVideoCodec + " -crf 10 -preset slow "
			OutputType = DefaultVideoType
		default:
			PrintError(errors.New("RecChunkRecieved: Encode type is not valid : " + CEncodeType))
			EncodeType = " "
			OutputType = DefaultVideoType
		}

		OutputChunkName := ChunkName + "." + OutputType
		WriteFile(wsroot+SavedChunkName, chunk)

		POut, PErr := ExcecCmd("ffmpeg -i " + SavedChunkName + EncodeType + OutputChunkName)
		println("\n\n-------- Chunk encode  -------------\n" + GetErrorString(PErr) + "\n" + POut + "\n---------------------\n\n")

		RecievedChunks[ChunkName] = OutputChunkName

		go DeleteFiles(wsroot + SavedChunkName)

	}()
	return BodyResponse([]byte("Host recieved " + path))
}

func FinalRecieved(para_paths string, body []byte) Response {

	paths := strings.Split(para_paths, "/")
	if len(paths) != 2 {
		PrintError(errors.New("FinalRecieved: Path is not valid : " + para_paths))
		return BodyResponse([]byte("Host : FinalRecieved : Paths are not valid : " + para_paths))
	}

	path := paths[0]
	newpath := paths[1]

	go func() {

		println("Finalizing " + path)

		if AudioEnabled {
			// End Audio recording
			EndTask, found := AudioTasks[path]
			if found {
				EndTask <- false
				delete(AudioTasks, path)
			}

			if newpath == "end" {
				ExcecCmdToString("pkill parec")
			} else {
				go startRecSysAudio(newpath)
			}
		}

		Sbody := string(body)
		matches := strings.Split(Sbody, "\n")

		ChunkList := make([]string, 0, len(matches)+2)

		for _, match := range matches {
			if len(match) > 0 {
				fl, found := RecievedChunks[match]

				RecieveWaitTimeout := 8000
				for !found && RecieveWaitTimeout > 0 {
					time.Sleep(100 * time.Millisecond)
					fl, found = RecievedChunks[match]
					RecieveWaitTimeout--
				}

				if RecieveWaitTimeout <= 0 && !found {
					PrintError(errors.New("FinalRecieved: Chunk not recieved on time : " + match))
					continue
				}

				delete(RecievedChunks, match)

				println("FinalRecieved : Checking chunk " + fl)
				chnk := wsroot + fl

				WaitTimeout := 2000
				for !FileExists(chnk) && WaitTimeout > 0 {
					time.Sleep(500 * time.Millisecond)
					WaitTimeout--
				}

				if WaitTimeout <= 0 {
					PrintError(errors.New("FinalRecieved: Chunk not found : " + fl))
					continue
				}

				ChunkList = append(ChunkList, fl)

			}
		}

		var ConcatChunkList []string

		if AudioEnabled {

			fflist := ""
			ConcatChunkList = make([]string, 0, 4)

			for _, chfile := range ChunkList {
				if strings.HasPrefix(chfile, "Ch-") {
					fflist += "file '" + chfile + "'\n"
				} else {
					ConcatChunkList = append(ConcatChunkList, chfile)
				}
			}

			ChFFLIST := path + "-pr.fflist"
			WriteFile(wsroot+ChFFLIST, []byte(fflist))

			VPath := path + "-video." + DefaultVideoType
			HiOut, HiErr := ExcecProgramToString("ffmpeg", "-f", "concat", "-safe", "0", "-i", ChFFLIST, "-c", "copy", VPath)

			AVPath := path + "-av." + DefaultVideoType
			MuxOut, MuxErr := ExcecCmdToString("ffmpeg -i " + VPath + " -i " + path + ".wav -map 0:v -map 1:a -c:v copy -shortest " + AVPath)

			ConcatChunkList = append(ConcatChunkList, AVPath)
			ChunkList = append(ChunkList, AVPath) // To delete

			go func() {
				println("\n\n" + HiOut + "\n\n")
				PrintError(HiErr)

				println("\n\n" + MuxOut + "\n\n")
				PrintError(MuxErr)

				go DeleteFiles(wsroot + path + ".wav")
				go DeleteFiles(wsroot + VPath)
				go DeleteFiles(wsroot + ChFFLIST)
			}()

		} else {
			ConcatChunkList = ChunkList
		}

		{
			fflist := ""
			for _, chfile := range ConcatChunkList {
				fflist += "file '" + chfile + "'\n"
			}

			WriteFile(wsroot+path+".fflist", []byte(fflist))

			HiOut, HiErr := ExcecProgramToString("ffmpeg", "-f", "concat", "-safe", "0", "-i", path+".fflist", "-c", "copy", path+"."+DefaultVideoType)

			println("\n\n" + HiOut + "\n\n")
			PrintError(HiErr)
		}

		RecievedChunks[path] = path + "." + DefaultVideoType

		go func() {
			for _, chnk := range ChunkList {
				go DeleteFiles(wsroot + chnk)
			}
		}()

		go DeleteFiles(wsroot + path + ".fflist")

		println("Finalized " + path)
	}()

	return BodyResponse([]byte("Final Recieved " + path))
}

func EndRecieved(path string, body []byte) Response {
	go func() {
		println("Ended " + path)
	}()

	return BodyResponse([]byte("End Recieved " + path))
}

func Handshake(path string) Response {
	go ExcecProgram("echo", "Host Ready")
	return BodyResponse([]byte("Host Ready"))
}

func StartRec(path string) Response {

	if !AudioEnabled || (len(SpeakerInputName) == 0) {
		AudioEnabled = false
		return BodyResponse([]byte("Started without system audio"))
	}

	go startRecSysAudio(path)

	return BodyResponse([]byte("Started with system audio"))
}

func startRecSysAudio(filename string) {
	EndTask := make(chan bool)

	ExcecCmdToString("pkill parec")
	ExcecCmdTask("parec -d "+SpeakerInputName+".monitor --file-format=wav "+filename+".wav", EndTask)

	AudioTasks[filename] = EndTask
}

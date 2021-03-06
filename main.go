package main

import (
	"errors"
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

var AudioEnabled bool = false
var DefaultVideoType string = "mkv"
var NoReEncode bool = false

var FFMPEGArgs string = " "
var HTTPPort string = "49542"

var EndFileDir string
var wsroot string = "workspace/"
var rootDir string = ""

var SpeakerInputName string = ""
var SSLEnabled bool = false

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
	println("                                      ")
	println("                                      ")

	HomeDir, HomeDirErr := os.UserHomeDir()
	if HomeDirErr == nil {
		HomeVideo := path.Join(HomeDir, "Videos/screencorder") + "/"
		if os.MkdirAll(HomeVideo, os.ModePerm) == nil {
			wsroot = HomeVideo
		}
	}
	EndFileDir = wsroot

	args := os.Args
	if RunArgs(args) {
		return
	}

	// wsroot might be changed by command line arguments

	if !strings.HasSuffix(wsroot, "/") {
		wsroot = wsroot + "/"
	}

	if !strings.HasSuffix(EndFileDir, "/") {
		EndFileDir = EndFileDir + "/"
	}

	MakeDir(wsroot)
	MakeDir(EndFileDir)

	if CheckDir(EndFileDir) {
		println("Output Directory " + EndFileDir + " is working")
	} else {
		println("Error : Output Directory " + EndFileDir + " is not working properly")
		return
	}

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

	if PrintError(InitExec([]string{"ffmpeg", "echo", "sh"})) {
		return
	}

	HiOut, HiErr := ExcecCmd("echo 'System calls working'")
	println(HiOut)
	if HiErr != nil {
		println(HiErr.Error())
	}

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
		println("Starting Screencorder http://localhost: " + HTTPPort)
	} else {
		println("Starting Screencorder https://localhost:" + HTTPPort)
	}

	if SSLEnabled {
		println("Connect to the same LAN and visit \n https://" + myip + ":" + HTTPPort + "   for host interface")
	} else {
		println("Visit \n  http://localhost:" + HTTPPort + "   for host interface")
	}

	if SSLEnabled {
		go OpenProgram("xdg-open", "https://localhost:"+HTTPPort)
	} else {
		go OpenProgram("xdg-open", "http://localhost:"+HTTPPort)
	}

	if SSLEnabled {
		if err := http.ListenAndServeTLS(":"+HTTPPort, wsroot+"server.crt", wsroot+"server.key", FullMux); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(":"+HTTPPort, FullMux); err != nil {
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

	if PrintError(InitExec([]string{"pacmd", "parec", "lame"})) {
		return
	}

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

func CheckDir(dir string) bool {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	if WriteFile(dir+"test-directory", []byte("Ensure this folder is working")) {
		if DeleteFiles(dir + "test-directory") {
			return true
		}
	}

	return false

}

// Returns true if the program needs to exit
func RunArgs(args []string) bool {

	SkipNext := false

	for i := 0; i < len(args); i++ {
		if SkipNext {
			SkipNext = false
			continue
		}

		switch args[i] {

		case "-h", "--help":
			PrintHelp()
			return true

		case "-ns", "-nosound":
			AudioEnabled = false

		case "-ps", "-parec-sound":
			AudioEnabled = true

		case "-t", "-type":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					DefaultVideoType = val
					SkipNext = true
				}
			}

		case "-vc", "-vcodec":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					if val == "auto" {
						FFMPEGArgs = " "
					} else {
						FFMPEGArgs = " -vcodec " + val
					}
					SkipNext = true
				}
			}

		case "-ws", "-workspace":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					wsroot = val
					SkipNext = true
				}
			}
		case "-o", "-output-dir":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					EndFileDir = val
					SkipNext = true
				}
			}

		case "-p", "-port":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					HTTPPort = val
					SkipNext = true
				}
			}

		case "-s", "-safe":
			AudioEnabled = false
			DefaultVideoType = "mkv"
			FFMPEGArgs = " "

		case "-ffmpeg":
			if i+1 < len(args) {

				for j := i + 1; j < len(args); j++ {
					if len(args[j]) != 0 {
						FFMPEGArgs += " " + args[j]
					}
				}
				return false
			}
		}

	}

	return false
}

func PrintHelp() {

	println(`
	
	Simple, fast screen recorder written in Go.

Usage : 
	-h, --help: Prints this help
	
	-ns, -nosound: Disables system sound recording. (Default)
	-ps, -parec-sound -: Enable system sound recording using 'parec' . Disabled by default.

	-t {filetype}, -type {filetype} : Sets the output file type and file extention. Default is mkv.
				 
				 If you are using no-re-encode option, you must use mkv or mp4.
				 Please test this with your web browser first.

				 If you are re-encoding the video,
				  you can use any video format that supports codec level concatenation 
				 These formats are supported as we know of : mp4, mkv, mpeg


	-vc, -vcodec: Sets the video codec for ffmpeg re-encoding. Default is 'auto'.

	             Please check with your ffmpeg installation for supported codecs. 
				 Pass '-vcodec auto' to let ffmpeg decide the video codec.
				 Pass '-vcodec libx264' for H.264 video codec.

				 Use command 'ffmpeg -codecs' to see the list of supported codecs.
				 Make sure it is in the list and 'DE' is present on capabilities.

	-ws, -workspace: Sets the workspace folder. Default is '$HOME/Videos/screencorder'

	-o, -output-dir: Sets the output folder. Default is '$HOME/Videos/screencorder'

	-p, -port: Sets the port for the web server. Default is '49542'

    -s, -safe: Safe mode for better compatibility.
	             This is same as '-vcodec auto -ns -t mkv'
`)

}

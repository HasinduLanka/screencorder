package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var Excecutables map[string]string = map[string]string{"ffmpeg": "", "echo": "", "bash": "", "rm": "", "pacmd": "", "parec": "", "lame": "", "xdg-open": ""}

func InitExec() error {
	for key := range Excecutables {
		path, err := exec.LookPath(key)
		if err != nil {
			return err
		}
		Excecutables[key] = path
		println("Excecutable " + key + " found at " + path)
	}
	return nil
}

func ExcecCmd(command string) (string, error) {
	return ExcecProgram("bash", "-c", command)
}
func ExcecCmdToString(command string) (string, error) {
	return ExcecProgramToString("bash", "-c", command)
}

func ExcecProgram(program string, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	println("Excecute " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = wsroot
	// configure `Stdout` and `Stderr`
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	// run command
	if err != nil {
		fmt.Println("Error:", err)
	}

	// out := string(ret)
	return "Done Excecute " + program + " " + args, err
}

func ExcecCmdTask(command string, endTask chan bool) (string, error) {
	return ExcecTask("bash", endTask, "-c", command)
}

func ExcecTask(program string, endTask chan bool, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	println("Excecute Task " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = wsroot
	// configure `Stdout` and `Stderr`
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Start()
	// run command
	if err != nil {
		fmt.Println("Error:", err)
	}

	Kill := <-endTask

	if Kill {
		PrintError(cmd.Process.Signal(os.Kill))
	} else {
		PrintError(cmd.Process.Signal(os.Interrupt))
	}

	// out := string(ret)
	return "Done Excecute Task " + program + " " + args, err
}

func ExcecProgramToString(program string, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	println("Excecute " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = wsroot
	// configure `Stdout` and `Stderr`
	cmd.Stderr = os.Stdout
	ret, err := cmd.Output()

	out := string(ret)
	return out, err
}

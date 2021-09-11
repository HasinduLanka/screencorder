package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var regex_url *regexp.Regexp = regexp.MustCompile("http.*")

func LoadURI(uri string) ([]byte, error) {
	if regex_url.MatchString(uri) {
		return DownloadFileToBytes(uri)
	} else {
		return LoadFile(uri)
	}
}

func AppendFile(filename string, content []byte) bool {
	F, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	F.Write(content)
	F.Close()
	return !PrintError(err)
}

func WriteFile(filename string, content []byte) bool {
	F, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	F.Write(content)
	F.Close()
	return !PrintError(err)

}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true

}

func MakeDir(name string) bool {
	return !PrintError(os.MkdirAll(name, os.ModePerm))
}

func DeleteFiles(name string) bool {
	println("Deleting file " + name)
	return !PrintError(os.RemoveAll(name))
}

func LoadURIToString(uri string) (string, error) {
	B, err := LoadURI(uri)
	return string(B), err
}

func LoadFile(filename string) ([]byte, error) {
	file, ferr := ioutil.ReadFile(filename)
	return file, ferr
}

func LoadFileToString(filename string) (string, error) {
	F, err := LoadFile(filename)
	return string(F), err
}

func DownloadToFile(filepath string, url string) error {

	body, derr := DownloadFileToStream(url)

	if derr != nil {
		return derr
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, body)
	PrintError(body.Close())
	return err
}

func DownloadFileToStream(url string) (io.ReadCloser, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp.Body, nil
}

func DownloadFileToBytes(url string) ([]byte, error) {
	str, err := DownloadFileToStream(url)

	if err != nil {
		return nil, err
	}
	out := StreamToByte(str)
	PrintError(str.Close())
	return out, nil
}

func DownloadFileToString(url string) (string, error) {
	str, err := DownloadFileToStream(url)

	if err != nil {
		return "", err
	}
	out := StreamToString(str)
	PrintError(str.Close())
	return out, nil
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

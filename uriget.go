package main

import (
	"bytes"
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

func AppendFile(filename string, content []byte) {
	F, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	CheckError(err)
	F.Write(content)
	F.Close()
}

func WriteFile(filename string, content []byte) {
	F, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	CheckError(err)
	F.Write(content)
	F.Close()
}

func MakeDir(name string) {
	os.MkdirAll(name, os.ModePerm)
}

func DeleteFiles(name string) {
	os.RemoveAll(name)
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
	CheckError(body.Close())
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
	CheckError(str.Close())
	return out, nil
}

func DownloadFileToString(url string) (string, error) {
	str, err := DownloadFileToStream(url)

	if err != nil {
		return "", err
	}
	out := StreamToString(str)
	CheckError(str.Close())
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

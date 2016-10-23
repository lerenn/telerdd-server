package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// MIME type
const (
	jpegMIME = "image/jpeg"
	pngMIME  = "image/png"
	gifMIME  = "image/gif"
)

func saveImage(host string, id int) (string, error) {
	var name string
	var imageJSON map[string]string

	// Get web response
	body, err := getJSON(fmt.Sprintf("http://"+host+"/messages/%d/image", id))
	if err != nil {
		fmt.Println("Can't read post list: " + err.Error())
		return name, err
	}

	// Make it form json to struct
	err = json.Unmarshal(body, &imageJSON)
	if err != nil {
		fmt.Println("Error when reading message list json:", err)
		return name, err
	}

	// Get mime and base64 img
	if imageJSON["img"] == "" {
		return name, errors.New("No img detected in response")
	}
	b64Img, mime := parseDataURL(imageJSON["img"])

	// Decode
	img, err := base64.StdEncoding.DecodeString(b64Img)
	if err != nil {
		return "", err
	}

	// Save image
	return writeImg(img, mime, id)
}

func parseDataURL(dataURL string) (string, string) {
	// Remove infos
	infos, b64Img := splitString(dataURL, ";base64,")
	_, mime := splitString(infos, ":")

	// Format mimes
	if mime == "image/jpg" {
		mime = jpegMIME
	}

	return b64Img, mime
}

func splitString(line, separator string) (string, string) {
	index := strings.Index(line, separator)
	if index < 0 {
		return line, ""
	}
	return line[:index], line[index+len(separator):]
}

func writeImg(img []byte, mime string, id int) (string, error) {
	ext := mime[6:]

	// Create directory if needed
	dir := "downloads/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	// Open file
	fn := dir + fmt.Sprintf("%d", id) + "." + ext
	f, err := os.Create(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Save in file
	_, err = f.Write(img)
	if err != nil {
		return "", err
	}

	return fn, nil
}

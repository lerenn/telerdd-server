package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Get addr
	host, err := getHostAddr()
	if err != nil {
		fmt.Println("Error when reading host address: " + err.Error())
		return
	}
	fmt.Println("Host addr: " + host)

	// Get post list
	list, err := getPostList(host)
	if err != nil {
		fmt.Println("Can't get post list: " + err.Error())
		return
	}

	// Get confirmation
	var confirmation string
	fmt.Println(fmt.Sprintf("There is %d messages (with or without images)", len(list["messages"])))
	fmt.Print("Do you want to save images from them ? [y/N] ")
	_, err = fmt.Scanf("%s", &confirmation)
	if err != nil || confirmation != "y" {
		return
	}

	// Get images
	for _, id := range list["messages"] {
		presence, err := checkImage(host, int(id))
		if err != nil {
			fmt.Println(fmt.Sprintf("Error for post #%d", int(id)))
			continue
		}

		if presence {
			fmt.Println(fmt.Sprintf("There is an image for post #%d", int(id)))
			name, err := saveImage(host, int(id))
			if err != nil {
				fmt.Println(fmt.Sprintf("Error when saving image from post #%d:", int(id)) + err.Error())
				continue
			}
			fmt.Println(fmt.Sprintf("Image from post #%d saved as %s", int(id), name))
		} else {
			fmt.Println(fmt.Sprintf("No image for post #%d", int(id)))
		}
	}
}

func getHostAddr() (string, error) {
	var host string
	fmt.Print("Enter your api host address (ex: \"api.example.com\"): ")
	nbr, err := fmt.Scanf("%s", &host)
	if err == nil && nbr != 1 {
		return "", errors.New("There should be only one address")
	}

	return host, err
}

func getPostList(host string) (map[string][]float64, error) {
	var msgsList map[string][]float64

	// Get web response
	body, err := getJSON("http://" + host + "/messages?status=all")
	if err != nil {
		fmt.Println("Can't read post list: " + err.Error())
		return msgsList, err
	}

	// Make it form json to struct
	err = json.Unmarshal(body, &msgsList)
	if err != nil {
		fmt.Println("Error when reading message list json:", err)
		return msgsList, err
	}

	return msgsList, nil
}

func getJSON(addr string) ([]byte, error) {
	var body []byte

	// Get web response
	resp, err := http.Get(addr)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	// Read web response
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Can't read post list: " + err.Error())
		return body, err
	}

	return body, nil
}

func checkImage(host string, id int) (bool, error) {
	var post map[string]string

	// Get web response
	body, err := getJSON(fmt.Sprintf("http://"+host+"/messages/%d", id))
	if err != nil {
		fmt.Println("Can't read post list: " + err.Error())
		return false, err
	}

	// Make it form json to struct
	err = json.Unmarshal(body, &post)
	if err != nil {
		fmt.Println("Error when reading message list json:", err)
		return false, err
	}

	if post["img"] == "true" {
		return true, nil
	}
	return false, nil
}

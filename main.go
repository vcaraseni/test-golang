package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var apiUrl = "https://app.asana.com/api/1.0"

type UData struct {
	Data []User `json:"data"`
}

type User struct {
	Gid string `json:"gid"`
}

type PData struct {
	Data []Project `json:"data"`
}

type Project struct {
	Gid string `json:"gid"`
}

func main() {
	processData()
	for i := range time.Tick(time.Second * 1) {
		if i.Second() == 30 || i.Second() == 0 { // each 30 second
			processData()
		}
	}
}

func processData() {
	// Users
	allUsers := doRequest(apiUrl + "/users")

	udata := UData{}
	if err := json.Unmarshal(allUsers, &udata); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}

	for _, gid := range udata.Data {
		usrInfo := getUserInfo(string(gid.Gid))
		saveToFile(gid.Gid+".json", usrInfo)
	}

	// Projects
	allProjects := doRequest(apiUrl + "/projects")

	pdata := PData{}
	if err := json.Unmarshal(allProjects, &pdata); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}

	for _, gid := range pdata.Data {
		usrInfo := getProjectInfo(string(gid.Gid))
		saveToFile(gid.Gid+".json", usrInfo)
	}
}

func getUserInfo(gid string) []byte {
	usrUrlPth := "/users/"
	return doRequest(apiUrl + usrUrlPth + gid)
}

func getProjectInfo(gid string) []byte {
	prUrlPath := "/projects/"
	return doRequest(apiUrl + prUrlPath + gid)
}

func doRequest(url string) []byte {
	bearer := "Bearer 2/1207527720536398/1207527735271853:7039039c59e186dd866ed29ea41d7641"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body
}

func saveToFile(fileName string, toSave []byte) {
	os.WriteFile(fileName, toSave, 0666)
}

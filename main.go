package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/samalba/dockerclient"
)

var Addresses []string

func parseInfo(info *dockerclient.Info) (role string, primary string) {
	for _, item := range info.DriverStatus {
		k, v := item[0], item[1]
		if strings.Contains(k, "Role") {
			role = v
		}
		if strings.Contains(k, "Primary") {
			primary = v
		}
	}
	return
}

func getClientInfo(client dockerclient.Client) (string, string, error) {
	info, err := client.Info()
	if err != nil {
		return "", "", err
	}
	role, primary := parseInfo(info)
	return role, primary, nil
}

type Status struct {
	Role    string
	Primary string
	Error   bool
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]Status)
	for _, a := range Addresses {
		client, err := dockerclient.NewDockerClient(a, nil)
		if err != nil {
			panic(err)
		}
		role, primary, err := getClientInfo(client)
		if err != nil {
			result[a] = Status{Error: true}
		} else {
			result[a] = Status{role, primary, false}
			fmt.Printf("%s role: %s | primary: %s\n", a, role, primary)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("The list of managers must be provided as arguments")
		os.Exit(1)
	}

	Addresses = os.Args[1:]

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/status", StatusHandler)

	log.Println("Listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// Ensure the "GoGoFiles" directory exists
	fmt.Println("Hi Chris")
	fmt.Println("Creating Directory GoGoFiles next to this program, Place your files in here.")
	fmt.Println("Warning, everything in this directory will be accessible to the network")

	createFolderIfNotExist("GoGoFiles")

	// Start the server
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./GoGoFiles/"))
	mux.Handle("/", http.StripPrefix("/", fileserver))

	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to get local IP address: %v", err)
	}

	fmt.Printf("The folder is now available on http://%s:8082\n", ip)
	fmt.Println("Note if you have multiple addresses, you may need to change the IP address in the URL. Sorry Chris.")

	err = http.ListenAndServe(":8082", mux)
	log.Fatal(err)
}

// Make the folder we are serving the files from
func createFolderIfNotExist(folderName string) {
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create folder: %v", err)
		}
	}
}

// getLocalIP returns the non-loopback local IP address of the host
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no non-loopback address found")
}

package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Embed the templates directory
//
//go:embed templates/*
var templateFiles embed.FS

// Embed the static directory
//
//go:embed static/*
var staticFiles embed.FS

type FileInfo struct {
	Name         string
	Size         int64
	LastModified string
	IsDir        bool
	Path         string
}

func main() {
	// Ensure the "GoGoFiles" directory exists
	fmt.Println("Hi Chris")
	fmt.Println("Creating Directory GoGoFiles next to this program, Place your files in here.")
	fmt.Println("Warning, everything in this directory will be accessible to the network")

	createFolderIfNotExist("GoGoFiles")

	// Start the server
	mux := http.NewServeMux()

	// Serve static files from embedded filesystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Direct file downloads
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./GoGoFiles"))))

	// Home page with file listing
	mux.HandleFunc("/", handleHome)

	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to get local IP address: %v", err)
	}

	fmt.Printf("The folder is now available on http://%s:8082\n", ip)
	fmt.Println("Note if you have multiple addresses, you may need to change the IP address in the URL. Sorry Chris.")

	err = http.ListenAndServe(":8082", mux)
	log.Fatal(err)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files, err := listFiles("GoGoFiles")
	if err != nil {
		http.Error(w, "Error listing files", http.StatusInternalServerError)
		return
	}

	// Get all local IP addresses
	localIPs, err := getLocalIPs()
	if err != nil {
		log.Printf("Warning: Failed to get local IP addresses: %v", err)
		localIPs = []string{}
	}

	// Check if the current host is localhost or 127.0.0.1
	host := r.Host
	isLocalhost := false
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		isLocalhost = true
	}

	// Create template functions map
	funcMap := template.FuncMap{
		"formatBytes": formatBytes,
	}

	// Parse template from embedded filesystem
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFS(templateFiles, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseURL := fmt.Sprintf("http://%s:8082", r.Host)
	data := struct {
		Files       []FileInfo
		BaseURL     string
		IsLocalhost bool
		LocalIPs    []string
	}{
		Files:       files,
		BaseURL:     baseURL,
		IsLocalhost: isLocalhost,
		LocalIPs:    localIPs,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// formatBytes converts bytes to human-readable formats (KB, MB, GB, etc.)
func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func listFiles(dir string) ([]FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:         file.Name(),
			Size:         file.Size(),
			LastModified: file.ModTime().Format(time.RFC1123),
			IsDir:        file.IsDir(),
			Path:         filepath.Join(dir, file.Name()),
		})
	}

	return fileInfos, nil
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

// getLocalIPs returns the non-loopback local IP addresses of the host
func getLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no non-loopback address found")
	}

	return ips, nil
}

// For backwards compatibility - returns the first IP address
func getLocalIP() (string, error) {
	ips, err := getLocalIPs()
	if err != nil {
		return "", err
	}
	return ips[0], nil
}

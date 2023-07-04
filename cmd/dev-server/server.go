package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/prometheus"

	"github.com/gorilla/mux"
)

const (
	templateFolder  = "./web/templates"
	resourceFolder  = "./web/static"
	localPort       = ":8080"
	refreshInterval = 1 * time.Second
	openOn          = ""
)

func main() {
	log.Printf("Server started at http://localhost%s", localPort)
	needsRefresh := false
	// Folder watcher
	go func() {
		fileWatch := map[string]time.Time{}
		for {
			time.Sleep(3 * time.Second)
			f, err := os.Open(templateFolder)
			if err != nil {
				log.Fatal("Couldnt open templateFolder", err)
			}
			files, _ := f.Readdir(-1)
			f.Close()

			for _, file := range files {
				info, _ := os.Stat(f.Name() + "/" + file.Name())
				if err != nil {
					// TODO: handle errors (e.g. file not found)
					log.Fatal("Couldnt stat file", err)
				}
				if info.ModTime().After(fileWatch[file.Name()]) {
					fileWatch[file.Name()] = info.ModTime()
					needsRefresh = true
				}
			}
		}
	}()

	router := mux.NewRouter().StrictSlash(true)
	staticRouter := router.PathPrefix("/static/")
	staticRouter.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(resourceFolder))))

	fs := flag.NewFlagSet("dev-server", flag.ExitOnError)
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	prometheusApp := prometheus.New(prometheusConfig)

	router.Handle("/metrics", prometheusApp.Handler()).Methods(http.MethodGet)

	router.HandleFunc("/reload", func(wr http.ResponseWriter, req *http.Request) {
		if needsRefresh {
			log.Println("Forcing reload")
			wr.WriteHeader(http.StatusUpgradeRequired)
			return
		}
		// Weird quirk but with an empty response and status code for no content, Chrome still views it
		// as a failed load
		fmt.Fprint(wr, "{}")
	}).Methods(http.MethodGet)
	router.HandleFunc("/reload.js", func(wr http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(wr, fmt.Sprintf(`
		setInterval(function() {
			fetch("/reload")
			.then(function (res) {
				if (res.status == 426) {
					window.location.reload(true);
				}
			})
		}, %d);
		`, refreshInterval.Milliseconds()))
		needsRefresh = false
	}).Methods(http.MethodGet)

	router.HandleFunc("/", loadMain).Methods(http.MethodGet)

	router.HandleFunc("/{route_name}", func(wr http.ResponseWriter, req *http.Request) {
		templateName := mux.Vars(req)["route_name"]

		// Define the global layout template
		layoutPath := filepath.Join(templateFolder, "layout.tpl.html")
		layoutFileContents, err := ioutil.ReadFile(layoutPath)
		if err != nil {
			fmt.Fprintf(wr, "Can't read file '%s' - %s", layoutPath, err)
			return
		}

		// Create a new template with the layout
		tmpl, err := template.New(templateName).Parse(string(layoutFileContents))
		if err != nil || tmpl == nil {
			fmt.Fprintf(wr, "Can't create template with layout '%s' - %s", templateName, err)
			return
		}

		// Parse additional templates
		_, err = tmpl.ParseGlob(filepath.Join(templateFolder, "*.tpl.html"))
		if err != nil {
			fmt.Fprintf(wr, "Can't parse more templates '%s' - %s", templateName, err)
			return
		}

		// Execute the template with the specific content
		err = tmpl.ExecuteTemplate(wr, templateName, nil)
		if err != nil {
			fmt.Fprintf(wr, "Can't execute template '%s' - %s", templateName, err)
			return
		}
	}).Methods(http.MethodGet)

	// Disable Caching of results
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
				wr.Header().Set("Cache-Control", "max-age=0, must-revalidate")
				next.ServeHTTP(wr, req)
			})
		},
	)

	devSrv := http.Server{
		Addr:           localPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Open localhost and start server
	go open(fmt.Sprintf("http://localhost%s/%s", localPort, openOn))
	devSrv.ListenAndServe()
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func loadMain(wr http.ResponseWriter, req *http.Request) {
	templateName := "index.tpl.html"
	templatePath := filepath.Join(templateFolder, templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(wr, fmt.Sprintf("Can't read file '%s' - %s", templatePath, err), http.StatusInternalServerError)
		return
	}

	// Load additional templates
	_, err = tmpl.ParseGlob(filepath.Join(templateFolder, "*.tpl.html"))
	if err != nil {
		http.Error(wr, fmt.Sprintf("Can't load additional templates - %s", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(wr, nil)
	if err != nil {
		http.Error(wr, fmt.Sprintf("Can't execute template '%s' - %s", templateName, err), http.StatusInternalServerError)
		return
	}
}

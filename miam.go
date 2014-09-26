package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const resourceDir = "resources"
const staticDir = "static"
const templateDir = "templates"

type element struct {
	Index  int
	Valeur string
}

type liste struct {
	Static_dir     string
	Title          string
	Raw_body       string
	Processed_body map[string][]element
}

type index struct {
	Static_dir string
	Title      string
	T_names    []string
}

// main Handle qui me retourne une page avec la gestion de la liste des pages existantes
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Utilise le template pour faire une redirection vers cette page.
	switch {
	case strings.Contains(r.URL.Path, "style.css"):
		chttp.ServeHTTP(w, r)
	case r.URL.Path == "/":
		log.Println(r.URL.Path)
		names, err := listeFiles()
		if err != nil {
			http.Error(w, "Could not retrieve list of files", http.StatusInternalServerError)
			log.Fatal("Could not retrieve list of files")
		}
		s_index := index{Static_dir: staticDir, Title: "Liste", T_names: names}
		templates.ExecuteTemplate(w, "index.html", s_index)
	default:
		regex := regexp.MustCompile("/([^/]*)(\\..*)*$")
		matches := regex.FindStringSubmatch(r.URL.Path)
		s_liste := liste{Static_dir: staticDir}
		s_liste.Title = matches[1]
		err := loadListe(&s_liste)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		processBody(&s_liste)
		templates.ExecuteTemplate(w, "liste.html", s_liste)
		//http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func processBody(l *liste) {
	var carte = make(map[string][]element)
	var menu string
	var idx int = 0
	for _, line := range strings.Split(l.Raw_body, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 1 && line[0] != '#' && strings.TrimSpace(line) != "" {
			if len(line) > 2 && line[0] == '=' {
				menu = strings.TrimSpace(line[1:])
				carte[menu] = make([]element, 0)
			} else if len(line) > 2 && line[0] == '-' {
				carte[menu] = append(carte[menu], element{Index: idx, Valeur: strings.TrimSpace(line[1:])})
				idx += 1
			}
		}
	}
	l.Processed_body = carte
}

func loadListe(l *liste) error {
	filename := resourceDir + "/" + l.Title + ".txt"
	stat, err := os.Stat(filename)
	if err != nil {
		return (err)
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	buf := make([]byte, stat.Size())
	if _, err := file.Read(buf); err != nil {
		return err
	}
	l.Raw_body = fmt.Sprintf("%s", buf)
	return nil
}

func listeFiles() ([]string, error) {
	dir, err := os.Open(resourceDir)
	if err != nil {
		log.Fatalf("Could not open dir %v for listing", resourceDir)
		return nil, err
	}
	fileListe, err := dir.Readdirnames(0) //liste all elements from resourceDir
	if err != nil {
		log.Fatalf("Could not retrieve list of files")
		return nil, err
	}
	var out []string = make([]string, 0) // Création d'une liste vide
	for _, file := range fileListe {
		if ind := strings.LastIndex(file, "."); ind != -1 {
			out = append(out, file[0:ind])
		}
	}
	return out, nil
}

var chttp = http.NewServeMux()
var templates = template.Must(template.ParseFiles(templateDir+"/index.html", templateDir+"/liste.html"))

func main() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", mainHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Fail to listen on port 8080")
	}
}

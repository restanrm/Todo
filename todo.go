package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"path/filepath"
)

type configuration struct {
	resourceDir, staticDir, templateDir string
	basePath                            string
	listenAddress                       string
	templates                           *template.Template
}

var conf configuration // variable globale de configuration de l'application

type element struct {
	Index  int
	Valeur string
}

type index struct {
	Static_dir string
	Title      string
	T_names    []string
}

// main Handle qui me retourne une page avec la gestion de la liste des pages existantes
func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "style.css"):
		chttp.ServeHTTP(w, r)
	case r.URL.Path == "/":
		names, err := listFiles(conf.resourceDir)
		if err != nil {
			http.Error(w, "Could not retrieve list of files", http.StatusInternalServerError)
			log.Fatal("Could not retrieve list of files")
		}
		s_index := index{Static_dir: conf.staticDir, Title: "Liste", T_names: names}
		conf.templates.ExecuteTemplate(w, "index.html", s_index)
	case r.FormValue("liste") != "":
		s_liste := getTitle(r)
		s_liste.Raw_body = r.FormValue("liste")
		if err := s_liste.saveListe(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		s_liste.processBody()
		conf.templates.ExecuteTemplate(w, "liste.html", s_liste)
		/*
		* Idée pour l'intégration de la modification pour le status:
		* Intégration d'un nouveau fichier permettant le stockage des données de status.
		* En cas de création de nouvelle page (enregistrement) -> réinitialisation de ce fichier avec les nouvelles données.
		*
		* Utilisation du format JSON, abandon de la gestion par fichier brut comme on le fait maintenant
		*
		* Utilisation d'un autre format plus complexe nécessitant plus de caractères de gestion
		 */
	default:
		s_liste := getTitle(r)
		err := s_liste.loadListe()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return 
		}
		s_liste.processBody()
		conf.templates.ExecuteTemplate(w, "liste.html", s_liste)
	}
}

func getTitle(r *http.Request) liste {
	s_liste := liste{Static_dir: conf.staticDir}
	matches := regex_title_page.FindStringSubmatch(r.URL.Path)
	s_liste.Title = matches[1]
	return s_liste
}
func listFiles(dirpath string) ([]string, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		log.Fatalf("Could not open dir %v for listing", dirpath)
		return nil, err
	}
	fileListe, err := dir.Readdirnames(0) //liste all elements from conf.resourceDir
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
var regex_title_page = regexp.MustCompile("/([^/]*)(\\..*)*$")

func main() {
	conf = configuration{resourceDir: "resources", templateDir: "templates", staticDir: "static"}
	var confDir = flag.String("config", "", "Dossier de données permettant le fonctionnement du service")
	var adresse = flag.String("adresse", ":8080", "Adresse d'écoute pour proposer le service")
	flag.Parse()
	conf.basePath = *confDir
	conf.listenAddress = *adresse
	if conf.basePath == "" {
		var err error
		conf.basePath, err = filepath.Abs(".")
		if err != nil {
			log.Fatal("Could not set working dir as main configuration dir ", err)
		}
	}

	// Go into configuration directory
	f_confdir, err := os.Open(conf.basePath)
	if err != nil {
		log.Fatal("Could not open directory: ", err)
	}
	if err := f_confdir.Chdir(); err != nil {
		log.Fatal("Could not go in configuration directory: ", err)
	}

	// load templates in cache 
	conf.templates = template.Must(template.ParseFiles(conf.templateDir+"/index.html", conf.templateDir+"/liste.html"))

	// Start webService
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", mainHandler)
	err = http.ListenAndServe(conf.listenAddress, nil)
	if err != nil {
		log.Fatal("Fail to listen on ", conf.listenAddress)
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"../github.com/gorilla/mux"
)

func main() {
	serveWeb()
}

var themeName = getThemeName()
var staticPages = populateStaticPages()

type gopher struct {
	Name string
	Age  int
	Cars []string
}

func serveWeb() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/", serveHome)
	gorillaRoute.HandleFunc("/index", serveGopherInfo)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}

func serveHome(resWriter http.ResponseWriter, r *http.Request) {
	staticPage := staticPages.Lookup("home.html")
	if staticPage == nil {
		resWriter.WriteHeader(404)
	}
	staticPage.Execute(resWriter, nil)
}

func serveGopherInfo(resWriter http.ResponseWriter, r *http.Request) {
	staticPage := staticPages.Lookup("index.html")
	if staticPage == nil {
		resWriter.WriteHeader(404)
	}

	gopherJSON := []byte(`{"Name":"Lamia","Age":20,"Cars": ["Ford", "BMW", "Fiat"] }`)

	var newGopher gopher

	err := json.Unmarshal(gopherJSON, &newGopher)
	if err != nil {
		panic(err)
	}

	staticPage.Execute(resWriter, newGopher)
}

func getThemeName() string {
	return "bootstrap4"
}

func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)
	basePath := "pages"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	for _, pathInfo := range templatePathsRaw {
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	basePath = "themes"
	templateFolder, _ = os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ = templateFolder.Readdir(-1)
	for _, pathInfo := range templatePathsRaw {
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	result.ParseFiles(*templatePaths...)
	return result
}

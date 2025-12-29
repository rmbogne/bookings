package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rmbogne/bookings/pkg/config"
	"github.com/rmbogne/bookings/pkg/models"
)

//var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Rendertemplate using html pages
// Right approach is to look at the template on disk the first time, then
// save it in the cache for future usage
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//TODO 1: Create a template cache
	//Get the template cache from the app Config
	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//TOD2: get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td) // this is set to catch error

	//TODO3: Render the template
	_, err := buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}

}

// TODo1: Create a template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	//anothe r way to declare a map is
	myCache := map[string]*template.Template{}

	//1. get all of the files named *.page.gohtml from ./templates
	pages, err := filepath.Glob("./templates/*.page.gohtml")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		//get only the file name here
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		//now let's look at the layout
		matches, err := filepath.Glob("./templates/*.layout.gohtml")

		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}

/*
// This approach has the constraint where you need to always comme and add layout in the render
// looking to avoid reading the template from the disc before parsing it
var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	//save the layout in a data structure, and find from that structure to get it parsed
	var tmpl *template.Template
	var err error

	//Check if we already have the template in our cache
	_, inMap := tc[t]
	if !inMap {
		//Need t ocreate the template, read from disc and parse it
		log.Println("creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.gohtml",
	}

	//parse the template
	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		return err
	}

	// add template to cache (map)
	tc[t] = tmpl

	return nil

}
*/

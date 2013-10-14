package cmds

import (
	"fmt"
	"github.com/drone/routes"
	"github.com/hoisie/mustache"
	"log"
	"net/http"
	"os"
	"path"
)

var ServeDoc = `
Serve the folkpocket website. The port and pocket appId can be configured from config.toml
`

func ServeAction() error {
	mux := routes.New()
	mux.Get("/", homeAction)
	mux.Get("/stories", homeAction)
	mux.Get("/confirmation", homeAction)

	// serve the public directory
	pwd, _ := os.Getwd()
	mux.Static("/public", pwd)

	// listen and serve
	http.Handle("/", mux)
	return http.ListenAndServe(":4000", logger(http.DefaultServeMux))

	// precompile all of the mustache templates

	return nil
}

// logger logs every http access
func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// render renders a view with the default layout
func render(view string) string {
	tplName := fmt.Sprintf("%s.html.mustache", view)
	return mustache.RenderFileInLayout(path.Join("views", tplName), path.Join("views", "layout.html.mustache"), nil)
}

// homeAction renders the home page
func homeAction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, render("home"))
}

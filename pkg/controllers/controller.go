package controllers

import (
	"net/http"

	"github.com/CloudyKit/jet"
	ghastApp "github.com/bradcypert/ghast/pkg/app"
	ghastContainer "github.com/bradcypert/ghast/pkg/container"
)

// GhastController should be embedded into consumer controllers
// and provides helper functions for working with the responseWriter, etc.
type GhastController struct{}

// Container returns the DI container associated with the given controller/request pairing.
func (c GhastController) Container() *ghastContainer.Container {
	return ghastApp.AppContext.Value("ghast/container").(*ghastContainer.Container)
}

// PathParam Get a Path Parameter from a given request and key
func (c GhastController) PathParam(r *http.Request, key string) interface{} {
	return r.Context().Value(key)
}

// View executes a view from the app templates
func (c GhastController) View(name string, w http.ResponseWriter, vars jet.VarMap, contextualData interface{}) {
	tmpl, err := ghastApp.GetApp(c.Container()).GetViewSet().GetTemplate(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err = tmpl.Execute(w, vars, contextualData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

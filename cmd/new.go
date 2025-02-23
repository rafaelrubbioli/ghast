package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// newCmd cobra command to help generate a new Ghast project
var newCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a new Ghast project",
	Long:  `Creates a new Ghast project based off the provided project name.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		runDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Unable to get working directory when creating new ghast app")
		}
		fmt.Print("Please enter your root package name: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		pkgName := strings.Replace(text, "\n", "", -1)

		// make relevant directories
		os.Mkdir(projectName, 0777)
		os.Mkdir(projectName+"/views", 0777)
		os.Mkdir(projectName+"/controllers", 0777)

		type pkg struct {
			Pkg string
		}
		// make mod file
		modFileTemplate := template.Must(template.New("mod").Parse(modTemplate))
		os.Chdir("./" + projectName)
		f, err := os.Create("./go.mod")
		if err != nil {
			panic("Unable to create new Ghast application controller")
		}
		modFileTemplate.Execute(f, pkg{pkgName})
		f.Close()
		os.Chdir(runDir)

		// make initial controller
		controllerTemplate := template.Must(template.New("controller").Parse(demoControllerTemplate))
		os.Chdir("./" + projectName + "/controllers")
		f, err = os.Create("./HomeController.go")
		if err != nil {
			panic("Unable to create new Ghast application controller")
		}
		controllerTemplate.Execute(f, nil)
		f.Close()
		os.Chdir(runDir)

		// make initial view
		viewTemplate := template.Must(template.New("view").Parse(viewTemplate))
		os.Chdir("./" + projectName + "/views")
		f, err = os.Create("./template.jet")
		if err != nil {
			panic("Unable to create new Ghast application template")
		}
		viewTemplate.Execute(f, nil)
		f.Close()
		os.Chdir(runDir)

		// make main file
		mainTemplate := template.Must(template.New("main").Parse(mainTemplate))
		f, err = os.Create(fmt.Sprintf("./%s/main.go", projectName))
		if err != nil {
			panic("Unable to create new Ghast application")
		}
		mainTemplate.Execute(f, pkg{pkgName})
		f.Close()

		os.Chdir("./" + projectName)
		// generate a YAML config
		yml := template.Must(template.New("yml").Parse(yamlTemplate))
		f, err = os.Create("./config.yml")
		if err != nil {
			panic("Unable to create ghast's config.yml application")
		}
		yml.Execute(f, pkg{projectName})
		f.Close()

		// fetch the go modules we need for ghast.
		goExecutable, err := exec.LookPath("go")
		cmdGoGet := &exec.Cmd{
			Path:   goExecutable,
			Args:   []string{goExecutable, "get", "-u", "./..."},
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
		}

		// A successful run of this has an exit code of 2
		cmdGoGet.Run()
		// if err = cmdGoGet.Run(); err != nil {
		//  panic("Unable to fetch go modules")
		// }

		fmt.Printf("Successfully created a new Ghast project in ./%s", projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

var modTemplate = `module {{.Pkg}}

go 1.13
`

var yamlTemplate = `
# modify freely, remove at your own risk
ghast:
  config:
    port: 9000

app:
  config:
    appName: {{.Pkg}}
`

var demoControllerTemplate = `
package controllers

import (
    "net/http"
    "github.com/CloudyKit/jet"
    ghastController "github.com/bradcypert/ghast/pkg/controllers"
)

type HomeController struct {
    ghastController.GhastController
}

func (c HomeController) Index(w http.ResponseWriter, r *http.Request) {
    vars := make(jet.VarMap)
    appName := c.Config("@app.config.appName").(string)
    vars.Set("AppName", appName)
    c.View("template.jet", w, vars, nil)
}
`

var viewTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Hello from Ghast!</title>
  </head>
  <body>
    <h1>Hello from {{ "{{ AppName }}" }}!</h1>
    <p>You've successfully scaffolded out your first Ghast application! There are tons of places to take this project from here! I've provided a few suggestions below:</p>
    <ul>
      <li>Building just an API? Delete the views folder if you don't need it.</li>
      <li>No interest in MVC? Toss the controllers directory and just use the router and app directly.</li>
      <li>Non-commital? No worries, Ghast's router handlers and controller functions all conform to <a href="https://golang.org/pkg/net/http/#HandleFunc">Go's standard HTTP Handler Func</a>. Moving away from Ghast isn't painless (nothing ever is) but shouldn't be too difficult either.</li>
      <li>Building something huge? Keep the existing patterns in place and use the CLI to generate new controllers, models, or migrations</li>
    </ul>
    <p>Don't forget! You can find an up-to-date readme here: <a href="https://www.github.com/bradcypert/ghast">Ghast's Readme</a>. Happy Ghasting!</p>
  </body>
</html>

`

var mainTemplate = `package main

import (
    "fmt"
    "net/http"

    ghastApp "github.com/bradcypert/ghast/pkg/app"
    ghastRouter "github.com/bradcypert/ghast/pkg/router"
    "{{.Pkg}}/controllers"
)

func main() {
    router := ghastRouter.Router{}

    // We can use controllers. Generate more using "ghast make controller MyControllerName" from your terminal
    router.Get("/", controllers.HomeController{}.Index)

    // Or we can just use standard Go HTTP handler funcs
    router.Get("/:name", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Hello "+r.Context().Value("name").(string))
    })
    
    app := ghastApp.NewApp()
    app.SetRouter(router)
    app.Start()
}
`

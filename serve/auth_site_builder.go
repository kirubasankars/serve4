package serve

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kirubasankars/metal"

	html "html/template"
)

type AuthSiteBuilder struct{}

func (au *AuthSiteBuilder) Build(module Module) {

	GetPath := func(ctx ServeContext, w http.ResponseWriter, r *http.Request) string {
		if ctx.Module.AppEnabled {
			return filepath.Join(ctx.Module.Path, "app")
		}
		return filepath.Join(ctx.Module.Path)
	}

	module.handlers["/"] = func(ctx ServeContext, w http.ResponseWriter, r *http.Request) {

		model := metal.NewMetal()
		templatePath := filepath.Join(GetPath(ctx, w, r), "index.html")
		tpl, _ := ioutil.ReadFile(templatePath)
		template, _ := html.New(templatePath).Parse(string(tpl))
		w.Header().Set("Content-Type", "text/html")

		formsAuth := new(FormsAuthentication)

		if r.Method == "POST" {

			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			formData := make(map[string]string)
			for key, values := range r.Form { // range over map
				for _, value := range values { // range over []string
					formData[key] = value
				}
			}

			if formsAuth.Authenticate(formData["username"], formData["password"]) {
				formsAuth.SetAuthCookie(formData["username"], ctx, w)
				http.Redirect(w, r, formData["redirectUrl"], http.StatusFound)
			} else {
				model.Set("ok", false)
				model.Set("message", "login failed")
				template.Execute(w, model.Raw())
			}
			return
		}

		if r.Method == "GET" {
			template.Execute(w, model.Raw())
		}
	}

	module.handlers["/signout"] = func(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
		fa := new(FormsAuthentication)
		fa.Signout(ctx, w, r)
	}

}

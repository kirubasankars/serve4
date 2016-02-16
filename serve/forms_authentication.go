package serve

import (
	"fmt"
	"net/http"
)

type FormsAuthentication struct{}

func (fa *FormsAuthentication) Authenticate(username string, password string) bool {
	if username == "admin" && password == "admin" {
		return true
	}
	return false
}

func (fa *FormsAuthentication) SetAuthCookie(username string, ctx ServeContext, w http.ResponseWriter) {

	cookieValue := make(map[string]string)
	cookieValue["username"] = username

	jar, path := ctx.GetJar()

	if encoded, err := jar.Encode("_auth", cookieValue); err == nil {
		cookie := &http.Cookie{
			Name:  "_auth",
			Value: encoded,
			Path:  path,
		}
		http.SetCookie(w, cookie)
	} else {
		fmt.Println(err)
	}
}

func (fa *FormsAuthentication) Signout(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
	_, p := ctx.GetJar()
	cookie := &http.Cookie{
		Name:   "_auth",
		Value:  "",
		Path:   p,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func (fa *FormsAuthentication) RedirectToLogin(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
	_, path := ctx.GetJar()
	http.Redirect(w, r, path+"/_auth/?redirectUrl="+r.URL.Path, http.StatusFound)
}

func (fa *FormsAuthentication) Validate(ctx *ServeContext, w http.ResponseWriter, r *http.Request) bool {
	jar, _ := ctx.GetJar()
	if cookie, err := r.Cookie("_auth"); err == nil {
		value := make(map[string]string)
		if err = jar.Decode("_auth", cookie.Value, &value); err == nil {
			return true
		}
	}
	return false
}

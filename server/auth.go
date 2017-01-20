package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

func initAuth(addr string) {
	gomniauth.SetSecurityKey("lorem-ipsum")
	gomniauth.WithProviders(
		google.New(
			"989480970338-1k7ssev1iap2n82jvmo5h1ucvj13k0fh.apps.googleusercontent.com",
			"7sUrCa7X9rraA2RxsjMm6Nb4",
			fmt.Sprintf("http://%s/auth/callback", addr),
		),
	)

}
type authHandler struct {
	next http.Handler
}

//ServeHTTP handle authentication
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated
		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}




func loginHandler(w http.ResponseWriter, r *http.Request) {
log.Println("login handler")
	act := mux.Vars(r)["action"]
	log.Println("act is", act)
	switch act {
	case "login" :
		t := templateHandler{filename: "login.html"}
		t.ServeHTTP(w,r)
	case "google" :
		provider, err := gomniauth.Provider("google")
		if err != nil {
			log.Fatal("Error getting provider", provider, err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatal("error trying to GeetBeginAuthURL for", provider, err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider("google")
		if err != nil {
			log.Fatal("error when trygin to get provider", provider, err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatal("errror trying to complet auth", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatal("error getting the user ", err)
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		userId := fmt.Sprintf("%x", m.Sum( []byte(user.Email())))
		m = md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		avatar := fmt.Sprintf("//gravatar.com/avatar/%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			"name":   user.Name(),
			"avatar": avatar,
			"email":  user.Email(),
			"userid": userId,
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "logon":
		email := r.FormValue("username")
		pass := r.FormValue("password")
		if pass != "password" {
			w.Header()["Location"] = []string{"/auth/login"}
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(email))
		avatar := fmt.Sprintf("//gravatar.com/avatar/%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			"name":   "some user",
			"avatar":  avatar,
			"email":  email,
			"userid": email,
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	}


}
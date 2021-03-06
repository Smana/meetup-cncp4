package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
)

func checkAuth(w http.ResponseWriter, r *http.Request, user string, pass string) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == user && pair[1] == pass
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get env variables created with the kubernetes secret
		cncpUser := os.Getenv("CNCP_USER")
		cncpPass := os.Getenv("CNCP_PASS")
		if cncpUser == "" || cncpPass == "" {
			log.Fatal("CNCP_PASS and CNCP_USER env vars not found !")
		}
		if checkAuth(w, r, cncpUser, cncpPass) {
			w.Write([]byte("Welcome to Cloud Native Computing Paris !! \n"))
			// myOriginalHandler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="CNCP REALM"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})

	http.ListenAndServe(":8080", nil)
}

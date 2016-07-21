package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func root(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Write(routesByte())
}

func showTime(w http.ResponseWriter, r *http.Request, p map[string]string) {
	t := time.Now()
	if len(p) == 1 {
		location, err := time.LoadLocation(p["in"])
		if err != nil {
			httpError(w, r.URL.Path, 422)
			return
		}
		t = t.In(location)
	}
	res := []byte("{\"Time\":\"" + t.String() + "\",\"Location\":\"" + t.Location().String() + "\"}")
	w.Write(res)
}

func headerInfo(w http.ResponseWriter, r *http.Request, p map[string]string) {
	r.Write(w)
}

func geoLocateInfo(w http.ResponseWriter, r *http.Request, p map[string]string) {
	address := strings.Split(r.RemoteAddr, ":")[0]
	forwardedFor := r.Header["X-Forwarded-For"]
	if len(forwardedFor) > 0 {
		address = forwardedFor[len(forwardedFor)-1]
	}

	res, err := http.Get("https://freegeoip.net/json/" + address)
	if err != nil {
		httpError(w, r.URL.Path, 500)
		return
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		httpError(w, r.URL.Path, 500)
		return
	}

	w.Write(content)
}

func main() {
	addRoute("/", []string{}, root)
	addRoute("/time", []string{"in"}, showTime)
	addRoute("/who", []string{}, headerInfo)
	addRoute("/where", []string{}, geoLocateInfo)

	port := os.Getenv("PORT")
	if port == "" {
		logger.err.Println("PORT cannot be set to \"\"")
		return
	}
	logger.info.Printf("Starting server on port: %v\n", port)
	startServer(":" + port)
}

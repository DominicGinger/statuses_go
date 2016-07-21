package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func root(w http.ResponseWriter, r *http.Request, p map[string]string) {
	json.NewEncoder(w).Encode(routesMap())
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
	res := map[string]string{"Time:": t.String(), "Location:": t.Location().String()}
	json.NewEncoder(w).Encode(res)
}

func headerInfo(w http.ResponseWriter, r *http.Request, p map[string]string) {
	r.Write(w)
}

func geoLocateInfo(w http.ResponseWriter, r *http.Request, p map[string]string) {
	res, err := http.Get("https://freegeoip.net/json/" + r.RemoteAddr)
	if err != nil {
		httpError(w, r.URL.Path, 500)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		httpError(w, r.URL.Path, 500)
		return
	}

	m := map[string]string{}
	json.Unmarshal(body, &m)
	json.NewEncoder(w).Encode(m)
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

package main

import (
	"encoding/json"
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

func main() {
	addRoute("/", []string{}, root)
	addRoute("/time", []string{"in"}, showTime)
	addRoute("/who", []string{}, headerInfo)

	port := os.Getenv("PORT")
	logger.info.Printf("Starting server on port: %v\n", port)
	startServer(":" + port)
}

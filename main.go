package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func root(w http.ResponseWriter, r *http.Request, p map[string]string) {
	res := map[string]string{"Page:": "/"}
	json.NewEncoder(w).Encode(res)
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

func main() {
	routes["/"] = route{root, []string{}, "/"}
	routes["/time"] = route{showTime, []string{"in"}, "/time"}
	logger.info.Println("Starting server")
	startServer(":" + os.Getenv("PORT"))
}

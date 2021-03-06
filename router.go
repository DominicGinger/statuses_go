package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

var errCodes = map[int]string{
	404: "Not Found",
	422: "Unprocessable Entity",
	500: "Internal Error",
}

type router struct{}

type route struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
	params  []string
	url     string
}

var routes = map[string]route{}

func routesByte() []byte {
	b := []byte("[")
	for _, v := range routes {
		str := "{\"url\":\"" + v.url + "\",\"params\":\"" + strings.Join(v.params[:], "\",\"") + "\"},"
		b = append(b, []byte(str)...)
	}
	b = append(b, []byte("]")...)
	return b
}

func addRoute(url string, params []string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	routes[url] = route{handler, params, url}
}

func processURL(uri *url.URL) (string, map[string]string) {
	url := uri.EscapedPath()
	query := uri.Query()
	params := map[string]string{}

	for k, v := range query {
		params[k] = v[0]
	}
	return url, params
}

func validParams(allowedParams []string, params map[string]string) bool {
	found := false
	for k := range params {
		for _, v := range allowedParams {
			if k == v {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func httpError(w http.ResponseWriter, url string, code int) {
	logger.err.Printf("%v [%v] %v \n", url, code, errCodes[code])
	http.Error(w, errCodes[code], code)
}

func resolveURL(uri *url.URL) (rte route, params map[string]string, code int) {
	url, params := processURL(uri)
	rte, ok := routes[url]
	code = 200
	if !ok {
		code = 404
	}
	if code == 200 && !validParams(rte.params, params) {
		code = 422
	}
	return
}

func (*router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rte, p, code := resolveURL(r.URL)
	if code != 200 {
		httpError(w, rte.url, code)
		return
	}

	start := time.Now()
	rte.handler(w, r, p)
	elapsed := time.Since(start)
	logger.info.Printf("%v - Elapsed: %v\n", rte.url, elapsed)
}

func startServer(port string) {
	http.ListenAndServe(port, &router{})
}

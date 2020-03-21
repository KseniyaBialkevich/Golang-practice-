package main

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForURLShortening(router chi.Router, format *render.Render) {

	mMap := make(map[string]string)

	//http://localhost:8080/process?url=https://en.wikipedia.org/wiki/URL_shortening
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")

		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)

		mMap[hashString] = url
		//mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"] = "https://en.wikipedia.org/wiki/URL_shortening"

		resultStringURL := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringURL)
	})

	//http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash") // 3c4d0ce2967a743e5ee1f2c4cb31e29e

		orignURL, _ := mMap[hashPath]
		//"https://en.wikipedia.org/wiki/URL_shortening", _ := mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"]

		http.Redirect(write, request, orignURL, 301) // перенаправление
	})

	// //http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e // либо через "внутренний" роутер intRouter
	// router.Route("/{hash}", func(intRouter chi.Router) {
	// 	intRouter.Get("/", func(write http.ResponseWriter, request *http.Request) {
	// 		hashPath := chi.URLParam(request, "hash")

	// 		orignUrl, _ := mMap[hashPath]

	// 		http.Redirect(write, request, orignURL, 301)
	// 	})
}

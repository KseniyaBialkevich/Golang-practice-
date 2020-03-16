package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/unrolled/render"
)

var format = render.New()

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(crs.Handler)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "public")
	FileServer(router, "", "/public", http.Dir(filesDir))

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		format.JSON(w, 200, "ok")
	})

	//******************************************************************************
	// http://localhost:8080/add?a=1&b=2
	router.Get("/add", func(write http.ResponseWriter, request *http.Request) {
		a := request .URL.Query().Get("a")
		ai, _ := strconv.Atoi(a)
		b := request .URL.Query().Get("b")
		bi, _ := strconv.Atoi(b)
		sum := ai + bi
		fmt.Fprintf(write, "Sum = %d", sum) // Fprintf пишет значения в write, a write это то куда отпр-ся результат
		//format.JSON(write, 200, sum)
		format.Text(write, 200, strconv.Itoa(sum))
	})
	//******************************************************************************
	// http://localhost:8080/addextra
	count := 0
	router.Get("/addextra", func(write http.ResponseWriter, request *http.Request) {
		count++
		//fmt.Fprintf(write, "Count: %d", count)
		format.JSON(write, 200, count)
	})
	//******************************************************************************
	// http: //localhost:8080/add?a=1&b=2
	router.Get("/add", func(write http.ResponseWriter, request *http.Request) {
		a := request.URL.Query().Get("a")
		ai, _ := strconv.Atoi(a)
		b := request.URL.Query().Get("b")
		bi, _ := strconv.Atoi(b)
		sum := ai + bi
		fmt.Fprintf(write, "Sum = %d", sum)
	})
	//******************************************************************************
	// http://localhost:8080/addsum
	countt := 0
	router.Get("/addsum", func(write http.ResponseWriter, request *http.Request) {
		countt = countt + 1
		result := struct {
			Count int `json:"count"`
		}{count}
		format.JSON(write, 200, result)
	})
	//******************************************************************************
	// http://localhost:8080/vote?entity=1
	type AnimalsCounter struct {
		dogCounter, catCounter, gopherCounter int
	}
	animalCounter := AnimalsCounter{0, 0, 0}
	router.Get("/vote", func(write http.ResponseWriter, request *http.Request) {
		var entity string = request.URL.Query().Get("entity")
		if entity == "1" {
			animalCounter.dogCounter++
		} else if entity == "2" {
			animalCounter.catCounter++
		} else if entity == "3" {
			animalCounter.gopherCounter++
		}
		format.JSON(write, 200, "ok")
	})
	// http://localhost:8080/addCandidate?name=bird&id=4
	type VoteInfo struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	voteCounter := make(map[int]VoteInfo)
	// http://localhost:8080/result
	type AddCandidate struct {
		ID    int `json:"id"`
		Count int `json:"count"`
	}
	router.Get("/addCandidate", func(write http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		sID := request.URL.Query().Get("id")
		id, _ := strconv.Atoi(sID)
		voteInfo, exist := voteCounter[id]
		if exist {
			voteInfo.Count++
			voteCounter[id] = voteInfo
		} else {
			voteCounter[id] = VoteInfo{name, 1}
		}
		voteInfo, _ = voteCounter[id]
		count := voteInfo.Count
		result := AddCandidate{id, count}
		format.JSON(write, 200, result)
	})
	router.Get("/result", func(write http.ResponseWriter, request *http.Request) {
		result := "{\n"
		isFirst := true
		for _, voteInfo := range voteCounter {
			if isFirst {
				result += fmt.Sprintf(" %s: %d", voteInfo.Name, voteInfo.Count)
				isFirst = false
			} else {
				result += fmt.Sprintf(",\n %s: %d", voteInfo.Name, voteInfo.Count)
			}
		}
		result += "\n}"
		fmt.Fprintf(write, result)
		//format.JSON(write, 200, result)
	})
	//******************************************************************************

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
		http.Redirect(write, request, orignURL, 301)
	})

	http.ListenAndServe(":8080", router)
}

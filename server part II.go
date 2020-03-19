package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
																																																																																		

	type Link struct {
		orignURL string
		accessTime []time.Time //срез accessTime типа time.Time => Time структура в пакете time
	}

	mMap := make(map[string]Link)

	//http://localhost:8080/process
	//?url=https://en.wikipedia.org/wiki/URL_shortening
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")

		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)

		times := []time.Time{} //инициализация переменной times типа срез, где каждый элемент имеет тип time.Time

		mMap[hashString] = Link{url, times} // инициализация структуры Link, присваиваем полям значения переменных url и times
		//mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"] = "https://en.wikipedia.org/wiki/URL_shortening", "times"

		resultStringURL := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringURL)
	})


	//http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash") // 3c4d0ce2967a743e5ee1f2c4cb31e29e

		copyOfLink, _ := mMap[hashPath] // возвращаем пер-й copyOfLink значение по ключу mMap, в copyOfLink хранится копия теперь

		copyOfLink.accessTime = append(copyOfLink.accessTime, time.Now()) //в поле (accessTime) структуры помещаем текущее время (time.Now())

		mMap[hashPath] = copyOfLink // теперь копию ложим в оригинал

		orignURL := copyOfLink.orignURL // в переменную помещаем значение (ссылку) поля (oringURL) структуры
		//"https://en.wikipedia.org/wiki/URL_shortening", _ := mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"]

		http.Redirect(write, request, orignURL, 301) //перенаправление
	})

	//http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e/times
	router.Get("/{hash}/times", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash")

		link, _ := mMap[hashPath] // возвращаем переменной link значение по ключу mMap, теперь link - это структура

		times := link.accessTime ////в переменную типа массив помещаем значение поля структуры link (accessTime)

		count := len(times) // длина среза times

		type JSONResult struct { // создание анонимной структуры
			Count int `json:"count"`
			Times []string `json:"times"`
		}

		strTimes := timesArrayToStringArray(times) // инициализируем переменную значением функции timesArrayToStringArray


		jsonResult := JSONResult{count, strTimes} // инициализируем сруктуру JSONResult,  присваиваем полям значения переменных count и strTimes
		format.JSON(write, 200, jsonResult)
	})


	http.ListenAndServe(":8080", router)
}

func timesArrayToStringArray(times []time.Time) []string { // создание функции с одним аргументом, тип - массив, элементы которого типа time.Time; возвращаемый тип


	strTimes := make([]string, len(times)) // создаем срез strTimes, элементы которого типа string, длина среза равна длине

	for idx, time := range times {
		strTime := time.Format("2006-01-02 15:04:05") // инициализируем переменную значением time
		strTimes[idx] = strTime // присваиваем значение элементам среза по каждому индексу
	}

	return strTimes // возвращаем массив типа int
}

// struct{
// count int
// times []string
// }{len(..), {"2020-03-18 22:24"}}

/*
testTime := time.Now()
fmt.Printf("%s\n", testTime.Format("2006-01-02 15:04:05"))

res := timesArrayToStringArray([]time.Time{time.Now(), time.Now(), time.Now()})
fmt.Printf("%s\n", res)
*/

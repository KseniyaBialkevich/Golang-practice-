package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func timesArrayToString(times []time.Time) []string {

	strTimes := make([]string, len(times))

	for idx, time := range times {
		strTime := time.Format("2006-01-02 15:04:05")
		strTimes[idx] = strTime
	}
	return strTimes
}

func check(err error, write http.ResponseWriter, format *render.Render) {
	log.Println(err)
	format.Text(write, 503, "Unable to save data.")
}

func found(isFound bool, write http.ResponseWriter, format *render.Render) {
	format.Text(write, 404, "No data was found for this hash.")
}

type Link struct {
	OrignURL   string      `json:"orignURL"`
	AccessTime []time.Time `json:"accessTime"`
}

type JSONResult struct {
	Count int      `json:"count"`
	Times []string `json:"times"`
}

func handlersForURLFile(router chi.Router, format *render.Render) {

	mMap := map[string]Link{}

	path, err := os.Getwd() //найти путь к текущему каталогу
	if err != nil {
		log.Fatalln(err)
	}
	pathToFile := fmt.Sprintf("%s/history.txt", path) //путь файла

	data, err := ioutil.ReadFile(pathToFile) //чтение файла
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(data, &mMap) //переобразование байтовых данные json в map

	//http://localhost:8080/user/process
	//url=https://en.wikipedia.org/wiki/URL_shortening expire=3m
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")

		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)

		times := []time.Time{}

		mMap[hashString] = Link{url, times}

		dataResult, err := json.Marshal(mMap) //преобразование данных map в байтовые данные/в json
		if err != nil {
			check(err, write, format)
			return
		}

		err = ioutil.WriteFile(pathToFile, dataResult, 0666) //запись данных в файл
		if err != nil {
			check(err, write, format)
			return
		}

		resultStringURL := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringURL)
	})

	//http://localhost:8080/user/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash")

		copyOfLink, isFound := mMap[hashPath]
		if !isFound {
			found(isFound, write, format)
			return
		}

		copyOfLink.AccessTime = append(copyOfLink.AccessTime, time.Now())

		mMap[hashPath] = copyOfLink

		dataResult, err := json.Marshal(mMap) //преобразование
		if err != nil {
			check(err, write, format)
			return
		}

		err = ioutil.WriteFile(pathToFile, dataResult, 0666) //запись данных в файл
		if err != nil {
			check(err, write, format)
			return
		}

		orignURL := copyOfLink.OrignURL

		http.Redirect(write, request, orignURL, 301) // перенаправление
	})

	//http://localhost:8080/user/3c4d0ce2967a743e5ee1f2c4cb31e29e/times
	router.Get("/{hash}/times", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash")

		link, isFound := mMap[hashPath]
		if !isFound {
			found(isFound, write, format)
			return
		}

		times := link.AccessTime

		count := len(times)

		strTimes := timesArrayToString(times)

		jsonResult := JSONResult{count, strTimes}

		format.JSON(write, 200, jsonResult)
	})
}

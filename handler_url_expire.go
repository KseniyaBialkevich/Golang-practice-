package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func timesArrayToStringI(times []time.Time) []string {

	strTimes := make([]string, len(times))

	for idx, time := range times {
		strTime := time.Format("2006-01-02 15:04:05")
		strTimes[idx] = strTime
	}
	return strTimes
}

func checkI(err error, write http.ResponseWriter, format *render.Render) {
	log.Println(err)
	format.Text(write, 503, "Unable to save data.")
}

func foundI(isFound bool, write http.ResponseWriter, format *render.Render) {
	format.Text(write, 404, "No data was found for this hash.")
}

func calcExpireTime(expire string) (time.Time, bool) { //вычисление, когда истечет время от времени создания сылки

	if expire == "" || len(expire) < 1 {
		log.Print()
		return time.Time{}, false
	}

	len := len(expire) //длина строки

	number := expire[0 : len-1] //substring
	unit := expire[len-1:]      //substring

	numberInt, err := strconv.Atoi(number) //int
	if err != nil {
		log.Print(err)
		return time.Time{}, false //вернуть пустое время и false
	}

	createTime := time.Now() //время создания ссылки

	var expireTime time.Time //истеченное время

	switch unit {
	case "d":
		expireTime = createTime.Add(time.Hour * 24 * time.Duration(numberInt))
	case "h":
		expireTime = createTime.Add(time.Hour * time.Duration(numberInt))
	case "m":
		expireTime = createTime.Add(time.Minute * time.Duration(numberInt))
	case "s":
		expireTime = createTime.Add(time.Second * time.Duration(numberInt))
	default:
		return time.Time{}, false
	}
	return expireTime, true
}

type LinkI struct {
	OrignURL       string      `json:"orignURL"`
	AccessTime     []time.Time `json:"accessTime"`
	ExpireLinkTime time.Time   `json:"expireLinkTime"`
}

type JSONResultI struct {
	Count int      `json:"count"`
	Times []string `json:"times"`
}

func handlersForURLExpire(router chi.Router, format *render.Render) {

	mMap := map[string]LinkI{}

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

	//http://localhost:8080/expire/process
	//url=https://en.wikipedia.org/wiki/URL_shortening expire=3m
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")
		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)

		expire := request.FormValue("expire")

		expireTime, ok := calcExpireTime(expire) //вызов функции
		if !ok {
			format.Text(write, 404, "Unit of time not found.")
			return
		}

		times := []time.Time{}

		mMap[hashString] = LinkI{url, times, expireTime}

		dataResult, err := json.Marshal(mMap) //преобразование данных map в байтовые данные/в json
		if err != nil {
			checkI(err, write, format)
			return
		}

		err = ioutil.WriteFile(pathToFile, dataResult, 0666) //запись данных в файл
		if err != nil {
			checkI(err, write, format)
			return
		}

		resultStringURL := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringURL)
	})

	//http://localhost:8080/expire/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash")

		copyOfLink, isFound := mMap[hashPath]
		if !isFound {
			foundI(isFound, write, format)
			return
		}

		timeRequestOfLink := time.Now() //время запроса ссылки

		if timeRequestOfLink.After(copyOfLink.ExpireLinkTime) { //если время запроса ссылки больше истекшего времени
			format.Text(write, 404, "Page Not Found.")
			return
		}

		copyOfLink.AccessTime = append(copyOfLink.AccessTime, time.Now())

		mMap[hashPath] = copyOfLink

		dataResult, err := json.Marshal(mMap) //преобразование
		if err != nil {
			checkI(err, write, format)
			return
		}

		err = ioutil.WriteFile(pathToFile, dataResult, 0666) //запись данных в файл
		if err != nil {
			checkI(err, write, format)
			return
		}

		orignURL := copyOfLink.OrignURL

		http.Redirect(write, request, orignURL, 301) // перенаправление
	})

	//http://localhost:8080/expire/3c4d0ce2967a743e5ee1f2c4cb31e29e/times
	router.Get("/{hash}/times", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash")

		linkI, isFound := mMap[hashPath]
		if !isFound {
			foundI(isFound, write, format)
			return
		}

		times := linkI.AccessTime

		count := len(times)

		strTimes := timesArrayToStringI(times)

		jsonResult := JSONResultI{count, strTimes}

		format.JSON(write, 200, jsonResult)
	})
}

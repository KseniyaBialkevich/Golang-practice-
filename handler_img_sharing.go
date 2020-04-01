package main

import (
	"crypto/md5"
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

func timesArrayToStringII(times []time.Time) []string {

	strTimes := make([]string, len(times))

	for idx, time := range times {
		strTime := time.Format("2006-01-02 15:04:05")
		strTimes[idx] = strTime
	}
	return strTimes
}

func checkII(err error, write http.ResponseWriter, format *render.Render) {
	log.Println(err)
	format.Text(write, 500, "Unable to save data.")
}

func foundII(isFound bool, write http.ResponseWriter, format *render.Render) {
	format.Text(write, 404, "No data was found for this hash.")
}

func calcExpireTimeI(expire string) (time.Time, error) { //вычисление, когда истечет время от времени создания сылки

	if expire == "" || len(expire) < 1 {
		return time.Time{}, fmt.Errorf("unit of time not found")
	}

	len := len(expire) //длина строки

	number := expire[0 : len-1] //substring
	unit := expire[len-1:]      //substring

	numberInt, err := strconv.Atoi(number) //int
	if err != nil {
		log.Print(err)
		return time.Time{}, fmt.Errorf("unit of time not found")
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
		return time.Time{}, fmt.Errorf("unit of time not found")
	}
	return expireTime, nil
}

//общая структура
type ImgTime struct {
	AccessTime []time.Time `json:"accessTime"`
	ExpireTime time.Time   `json:"expireTime"`
}

func handlersForImageSharing(router chi.Router, format *render.Render) {

	mapImgTime := map[string]ImgTime{}

	path, err := os.Getwd() //найти путь к текущему каталогу
	if err != nil {
		log.Fatalln(err)
	}

	//http://localhost:8080/image
	//upload = image.png expire = 10s
	router.Post("/image", func(write http.ResponseWriter, request *http.Request) {
		expire := request.FormValue("expire")

		request.ParseMultipartForm(10 << 20) //загрузка 10 мб файлов

		file, handler, err := request.FormFile("upload") // возвращает файл по ключу "upload"| имя файла/заголовок/размер | ошибку
		if err != nil {
			log.Println(err)
			format.Text(write, 500, "Error Retrieving the File")
			return
		}
		defer file.Close()

		mimeType := handler.Header.Get("Content-Type") //тип содержимого

		if !((mimeType == "image/jpeg") || (mimeType == "image/png")) {
			format.Text(write, 500, "The format file is not valid.")
			return
		}

		hashNameFile := md5.Sum([]byte(handler.Filename))
		stringNameFile := fmt.Sprintf("%x", hashNameFile)

		var imgHashAdress string

		if mimeType == "image/jpeg" {
			imgHashAdress = stringNameFile + ".jpeg"
		} else if mimeType == "image/png" {
			imgHashAdress = stringNameFile + ".png"
		}

		expireTime, err := calcExpireTimeI(expire) //вызов функции
		if err != nil {
			format.Text(write, 404, "Unit of time not found.")
			return
		}

		times := []time.Time{} // срез

		mapImgTime[imgHashAdress] = ImgTime{times, expireTime}

		// newFile, err := os.Create(path + "/public/imgs/" + imgHashAdress) //создание нового файла
		// if err != nil {
		// 	log.Println(err)
		// 	format.Text(write, 500, "Unable to create new file.")
		// 	return
		// }
		// defer newFile.Close()

		// _, err = io.Copy(newFile, file) //копирование из источника в новый файл
		// if err != nil {
		// 	log.Println(err)
		// 	format.Text(write, 500, "Cannot copy from source file to new file.")
		// 	return
		// }

		// ioutil.ReadAll    _OR_    io.Copy(newFile, file)

		data, err := ioutil.ReadAll(file) // чтение
		if err != nil {
			log.Fatal(err)
			return
		}

		pathFile := path + "/public/imgs/" + imgHashAdress // путь к изображению

		err = ioutil.WriteFile(pathFile, data, 0666) // запись
		if err != nil {
			checkII(err, write, format)
			return
		}

		format.Text(write, 200, "File uploaded successfully!")
		imageAdress := fmt.Sprintf("\nThe address of your uploaded image:\nhttp://localhost:8080/public/imgs/%s", imgHashAdress)
		format.Text(write, 200, imageAdress)
	})

	//http://localhost:8080/public/imgs/92177896a4998aec4800fe54c1e71f10.jpeg
	router.Get("/public/imgs/{imgHashName}", func(write http.ResponseWriter, request *http.Request) {
		imgHashName := chi.URLParam(request, "imgHashName")

		copyOfAdress, isFound := mapImgTime[imgHashName]
		if !isFound {
			foundII(isFound, write, format)
			return
		}

		timeRequestOfAdress := time.Now() //время запроса ссылки

		if timeRequestOfAdress.After(copyOfAdress.ExpireTime) { //если время запроса ссылки больше истекшего времени
			format.Text(write, 404, "page not found")
			return
		}

		copyOfAdress.AccessTime = append(copyOfAdress.AccessTime, time.Now())

		mapImgTime[imgHashName] = copyOfAdress

		pathFile := path + "/public/imgs/" + imgHashName // путь к изображению

		http.ServeFile(write, request, pathFile)
	})

	//http://localhost:8080/public/imgs/92177896a4998aec4800fe54c1e71f10.jpeg/history
	router.Get("/public/imgs/{imgHashName}/history", func(write http.ResponseWriter, request *http.Request) {
		imgHashName := chi.URLParam(request, "imgHashName")

		adress, isFound := mapImgTime[imgHashName]
		if !isFound {
			foundII(isFound, write, format)
			return
		}

		times := adress.AccessTime
		count := len(times)

		type JSONResult struct {
			Count int      `json:"count"`
			Times []string `json:"times"`
		}

		strTimes := timesArrayToStringII(times)

		jsonResult := JSONResult{count, strTimes}

		format.JSON(write, 200, jsonResult)
	})
}

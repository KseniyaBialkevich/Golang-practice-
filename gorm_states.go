package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// stucts for SHARING URL

// Таблица для хранения ID, оригинальной ссылки, ее хеш и время истекания доступа
type LinkInfo struct {
	ID         int `gorm:"primary_key"`
	OrignURL   string
	Hash       string `gorm:"unique;not null"`
	ExpireTime time.Time
}

// Таблица для хранения ID, ID хеша из таблицы LinkInfo и время посещения хеш-ссылки
type AccessTimeURL struct {
	ID         int `gorm:"primary_key"`
	LinkID     int
	AccessTime time.Time
}

// stucts for SHARING IMAGES

// Таблица для хранения ID, хеш-имени картинки и время истекания доступа
type ImgInfo struct {
	ID         int    `gorm:"primary_key"`
	Hash       string `gorm:"unique;not null"`
	ExpireTime time.Time
}

// Таблица для хранения ID, ID хеша из таблицы ImgInfo и время просмотра изображения по пути с хеш-именем файла
type AccessTimeIMG struct {
	ID         int `gorm:"primary_key"`
	HashID     int
	AccessTime time.Time
}

//Базу данных открываем/закрываем в файле main.go

func gormForStates(router chi.Router, format *render.Render) {

	//------------------------------SHORTING-URLS------------------------------//

	// Миграция схем
	db.AutoMigrate(&LinkInfo{}, &AccessTimeURL{})

	//http://localhost:8080/gorm/process
	//url=https://en.wikipedia.org/wiki/URL_shortening expire=3m
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")

		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)

		expire := request.FormValue("expire")

		expireTime, err := calcExpireTimeII(expire) //вызов функции
		if err != nil {
			format.Text(write, 500, "unit of time not found")
			return
		}

		// Создание (объявление и инициализация экземпляра структуры, перечисление имен полей с их значениями)
		linkInfo := LinkInfo{OrignURL: url, Hash: hashString, ExpireTime: expireTime}

		// Сохранение поля в БД
		db.Save(&linkInfo)

		resultStringURL := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringURL)
	})

	//http://localhost:8080/gorm/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hash := chi.URLParam(request, "hash")

		// Объявление экземпляра структуры в linkInfo
		linkInfo := LinkInfo{}

		// Запрос - получить запись по конкретному хешу и положить в linkInfo
		db.Where("hash = ?", hash).First(&linkInfo)

		timeRequestOfLink := time.Now() //время запроса ссылки

		if timeRequestOfLink.After(linkInfo.ExpireTime) { //если время запроса ссылки больше истекшего времени
			format.Text(write, 404, "Page Not Found.")
			return
		}

		// Объявлении экземпляра структуры в accessTimeUrl
		accessTimeURL := AccessTimeURL{LinkID: linkInfo.ID, AccessTime: time.Now()}

		// Сохранение поля в БД
		db.Save(&accessTimeURL)

		// Берем оригинальную ссылку из выбранного поля по хешу
		orignURL := linkInfo.OrignURL

		http.Redirect(write, request, orignURL, 301) // перенаправление
	})

	// //http://localhost:8080/gorm/3c4d0ce2967a743e5ee1f2c4cb31e29e/times
	router.Get("/{hash}/times", func(write http.ResponseWriter, request *http.Request) {
		hash := chi.URLParam(request, "hash")

		// Объявление экземпляра структуры в linkInfo
		linkInfo := LinkInfo{}

		// Запрос - получить запись по конкретному
		db.Where("hash = ?", hash).First(&linkInfo)

		// Создание среза с типом структура AccessTimeURL (для хранения структур времени доступов)
		var accessTimeSlice []AccessTimeURL

		// Запрос - получить запись по конкретному ID хеша и положить все значения в accessTimeSlice
		db.Where("link_id = ?", linkInfo.ID).Find(&accessTimeSlice)

		// Длина среза
		count := len(accessTimeSlice)

		// Инициализируем переменную значением функции timesArrayToStringForShortUrl
		strTimes := timesArrayToStringForShortURL(accessTimeSlice)

		type JSONOutput struct {
			Count int      `json:"count"`
			Times []string `json:"times"`
		}

		jsonResult := JSONOutput{count, strTimes}

		format.JSON(write, 200, jsonResult)
	})

	//------------------------------SHARING-IMAGES------------------------------//

	//Миграция схем
	db.AutoMigrate(&ImgInfo{}, &AccessTimeIMG{})

	//http://localhost:8080/gorm/image
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

		newFile, err := os.Create(path + "/public/imgs/" + imgHashAdress)
		if err != nil {
			log.Println(err)
			format.Text(write, 500, "Unable to create new file.")
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Println(err)
			format.Text(write, 500, "Cannot copy from source file to new file.")
			return
		}

		imgInfo := ImgInfo{Hash: imgHashAdress, ExpireTime: expireTime}

		db.Save(&imgInfo)

		format.Text(write, 200, "File uploaded successfully!")
		imageAdress := fmt.Sprintf("\nThe address of your uploaded image:\nhttp://localhost:8080/gorm/public/imgs/%s", imgHashAdress)
		format.Text(write, 200, imageAdress)
	})

	//http://localhost:8080/gorm/public/imgs/92177896a4998aec4800fe54c1e71f10.jpeg
	router.Get("/public/imgs/{imgHashName}", func(write http.ResponseWriter, request *http.Request) {
		imgHashName := chi.URLParam(request, "imgHashName")

		imgInfo := ImgInfo{}

		db.Where("hash = ?", imgHashName).First(&imgInfo)

		timeRequestOfAdress := time.Now() //время запроса пути

		if timeRequestOfAdress.After(imgInfo.ExpireTime) { //если время запроса ссылки больше истекшего времени
			format.Text(write, 404, "page not found")
			return
		}

		accessTimeIMG := AccessTimeIMG{HashID: imgInfo.ID, AccessTime: time.Now()}

		db.Save(&accessTimeIMG)

		pathFile := path + "/public/imgs/" + imgHashName // путь к изображению

		http.ServeFile(write, request, pathFile)
	})

	//http://localhost:8080/gorm/public/imgs/92177896a4998aec4800fe54c1e71f10.jpeg/history
	router.Get("/public/imgs/{imgHashName}/history", func(write http.ResponseWriter, request *http.Request) {
		imgHashName := chi.URLParam(request, "imgHashName")

		imgInfo := ImgInfo{}

		db.Where("hash = ?", imgHashName).First(&imgInfo)

		var accessTimeSlice []AccessTimeIMG

		db.Where("hash_id = ?", imgInfo.ID).Find(&accessTimeSlice)

		count := len(accessTimeSlice)

		strTimes := timesArrayToStringForShortIMG(accessTimeSlice)

		type JSONResult struct {
			Count int      `json:"count"`
			Times []string `json:"times"`
		}

		jsonResult := JSONResult{count, strTimes}

		format.JSON(write, 200, jsonResult)
	})
}

// Функция с одним арг. типа - срез, элементы которого типа типа структура AccessTimeURL
func timesArrayToStringForShortURL(accesses []AccessTimeURL) []string {
	// Срез типа string c длиной среза accesses
	strTimesSlice := make([]string, len(accesses))

	for idx, value := range accesses {
		accessTime := value.AccessTime                      // AccessTime - назв. поля структуры AccessTimeURL
		strTime := accessTime.Format("2006-01-02 15:04:05") // Инициализируем переменную значением accessTime типа time
		strTimesSlice[idx] = strTime                        // Присваиваем значение элементам среза типа string по каждому индексу
	}
	return strTimesSlice // Возвращаем срез, в котором содержится время посещений ссылки
}

// Функция с одним арг. типа - срез, элементы которого типа типа структура AccessTimeIMG
func timesArrayToStringForShortIMG(accesses []AccessTimeIMG) []string {

	strTimesSlice := make([]string, len(accesses))

	for idx, value := range accesses {
		accessTime := value.AccessTime
		strTime := accessTime.Format("2006-01-02 15:04:05")
		strTimesSlice[idx] = strTime
	}
	return strTimesSlice
}

func checkIII(err error, write http.ResponseWriter, format *render.Render) {
	log.Println(err)
	format.Text(write, 500, "Unable to save data.")
}

func foundIII(isFound bool, write http.ResponseWriter, format *render.Render) {
	format.Text(write, 404, "No data was found for this hash.")
}

// Вычисление, когда истечет время от времени создания сылки
func calcExpireTimeII(expire string) (time.Time, error) {

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

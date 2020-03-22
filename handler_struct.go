package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForStruct(router chi.Router, format *render.Render) {

	var lastUser string

	//http://localhost:8080/struct/user   заполнение структуры новыми значениями и возвращение их в виде строки (всех пользователей)
	//id=101 name=Barry surname=White age=59 sex=male
	router.Post("/user", func(write http.ResponseWriter, request *http.Request) {
		id := request.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		name := request.FormValue("name")
		surname := request.FormValue("surname")
		age := request.FormValue("age")
		ageInt, _ := strconv.Atoi(age)
		sex := request.FormValue("sex")

		UserInfo := User{idInt, name, surname, ageInt, sex}

		UserInfoResult := fmt.Sprintf("id: %d\nname: %s\nsurname: %s\nage: %d\nsex: %s\n", UserInfo.ID, UserInfo.Name, UserInfo.Surname, UserInfo.Age, UserInfo.Sex)

		lastUser = UserInfoResult

		format.Text(write, 200, UserInfoResult)
	})

	// http://localhost:8080/struct/user/last   вернуть последнего сохраненного пользователя
	router.Get("/user/last", func(write http.ResponseWriter, request *http.Request) {
		format.Text(write, 200, lastUser)
	})

	var userInfo string

	//http://localhost:8080/struct/user/name/Ada вернуть через путь параметр "name" пользователя
	router.Get("/user/name/{name}", func(write http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")
		userInfo = name

		format.Text(write, 200, userInfo)
	})

	//http://localhost:8080/struct/user/surname/Lovelace   вернуть через путь параметр "surname" пользователя
	router.Get("/user/surname/{surname}", func(write http.ResponseWriter, request *http.Request) {
		surname := chi.URLParam(request, "surname")
		userInfo = userInfo + " " + surname

		format.Text(write, 200, userInfo)
	})

	//http://localhost:8080/struct/user/age/37   вернуть через путь параметр "age" пользователя
	router.Get("/user/age/{age}", func(write http.ResponseWriter, request *http.Request) {
		age := chi.URLParam(request, "age")
		userInfo = userInfo + " " + age

		format.Text(write, 200, userInfo)
	})

	//http://localhost:8080/struct/user/sex/female   вернуть через путь параметр "sex" пользователя
	router.Get("/user/sex/{sex}", func(write http.ResponseWriter, request *http.Request) {
		sex := chi.URLParam(request, "sex")
		userInfo = userInfo + " " + sex

		format.Text(write, 200, userInfo)
	})

}

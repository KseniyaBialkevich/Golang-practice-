package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForMap(router chi.Router, format *render.Render) {

	mapUser := make(map[int]User)

	//http://localhost:8080/map/user //обработчик принимает характеристики пользователя и добавляет их в map
	//id=111 name=Winston surname=Churchill age=91 sex=male
	router.Post("/user", func(write http.ResponseWriter, request *http.Request) {
		id := request.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		name := request.FormValue("name")
		surname := request.FormValue("surname")
		age := request.FormValue("age")
		ageInt, _ := strconv.Atoi(age)
		sex := request.FormValue("sex")

		_, ok := mapUser[idInt]

		if ok { //проверка на существующего пользователя
			result := fmt.Sprintf("User with id %d already exist!", idInt)
			format.Text(write, 404, result)
		} else {

			mapUser[idInt] = User{idInt, name, surname, ageInt, sex}

			var result []string
			for _, value := range mapUser {
				resultString := value.String()
				//resultString := fmt.Sprintf("ID: %d\nName: %s\nSurname: %s\nAge: %d\nSex: %s\n", value.ID, value.Name, value.Surname, value.Age, value.Sex)
				result = append(result, resultString)
			}

			resultUsers := strings.Join(result, "\n")
			format.Text(write, 200, resultUsers)
		}
	})

	//http://localhost:8080/map/user/id/111 //возвращает пользователя по id из map
	router.Get("/user/id/{id}", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		value, ok := mapUser[idInt]

		if ok {
			format.Text(write, 200, value.String())
		} else {
			format.Text(write, 404, "User is not found.\n")
		}
	})

	//http://localhost:8080/map/user/name/Winston //возвращает пользователей по name из map
	router.Get("/user/name/{name}", func(write http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")

		searchUserByName := func(users map[int]User, name string) ([]string, bool) {
			isFound := false
			result := make([]string, 0)

			for _, value := range mapUser {
				if name == value.Name {
					result = append(result, value.String())
					isFound = true
				}
			}
			return result, isFound
		}

		userArray, isFound := searchUserByName(mapUser, name)

		if isFound {
			result := strings.Join(userArray, "\n")
			format.Text(write, 200, result)
		} else {
			format.Text(write, 404, "User is not found.\n")
		}
	})

	//http://localhost:8080/map/user/put/id/111 //обновляет переданные параметры по id, принятые в теле запроса
	//surname=Jerome age=92
	router.Put("/user/put/id/{id}", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		surnameForm := request.FormValue("surname")
		ageForm := request.FormValue("age")
		ageFormInt, _ := strconv.Atoi(ageForm)

		value, ok := mapUser[idInt]

		if ok {
			value.Surname = surnameForm
			value.Age = ageFormInt
			mapUser[idInt] = value
			format.Text(write, 200, "Data has been updated successfully!\n\n")
			format.Text(write, 200, value.String())
		} else {
			format.Text(write, 404, "User is not found.\n")
		}
	})

	//http://localhost:8080/map/user/ids/111,222,999 //принимает строку с идентификаторами пользователей разделенными запятыми
	router.Get("/user/ids/{ids}", func(write http.ResponseWriter, request *http.Request) {
		ids := chi.URLParam(request, "ids")

		idsArray := strings.Split(ids, ",")

		for indx, value := idsArray

	})
}

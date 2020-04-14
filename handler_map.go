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
					result = append(result, value.String()) // создание среза пользователей с одинаковыми именами
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

	//http://localhost:8080/map/user/ids?arg=111,222,999 //принимает строку с идентификаторами пользователей разделенными запятыми
	router.Get("/user/ids", func(write http.ResponseWriter, request *http.Request) {
		ids := request.URL.Query().Get("arg")

		idsArray := strings.Split(ids, ",")

		foundIds := make([]string, 0)
		notFoundIds := make([]string, 0)

		for _, elem := range idsArray {
			elemInt, _ := strconv.Atoi(elem)

			value, ok := mapUser[elemInt]
			if ok {
				result := value.String()
				foundIds = append(foundIds, result)
			} else {
				notFoundIds = append(notFoundIds, elem)
			}
		}

		if len(foundIds) > 0 {
			result := strings.Join(foundIds, "\n")
			format.Text(write, 200, result)
		}
		if len(notFoundIds) > 0 { //если есть ненайденные пользователи
			resultString := strings.Join(notFoundIds, ",")
			result := fmt.Sprintf("\nId(s): %s not found!", resultString)
			format.Text(write, 404, result)
		}
	})

	//http://localhost:8080/map/user/names?arg=Winston,Ada,Barry //принимат список имен, разделенных запятой
	router.Get("/user/names", func(write http.ResponseWriter, request *http.Request) {
		names := request.URL.Query().Get("arg")

		namesArray := strings.Split(names, ",")

		resultArray := make([]string, 0)

		for _, elem := range namesArray {

			for _, value := range mapUser {

				if elem == value.Name {
					result := value.String()
					resultArray = append(resultArray, result)
				}
			}
		}
		result := strings.Join(resultArray, "\n")
		format.Text(write, 200, result)
	})

	//http://localhost:8080/map/user/and?name=Winston&age=91 //возвращает пользователей у которых имя И возраст соответствуют
	router.Get("/user/and", func(write http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		age := request.URL.Query().Get("age")
		ageInt, _ := strconv.Atoi(age)

		resultUsersArray := make([]string, 0)

		for _, value := range mapUser {
			if name == value.Name && ageInt == value.Age {
				result := value.String()
				resultUsersArray = append(resultUsersArray, result)
			}
		}

		result := strings.Join(resultUsersArray, "\n")
		format.Text(write, 200, result)

		if len(resultUsersArray) < 1 { //если совпадения не найдены
			format.Text(write, 404, "No matches found!")
		}
	})

	//http://localhost:8080/map/user/or?name=Winston&age=91 //возвращает пользователей у которых имя ИЛИ возраст соответствуют
	router.Get("/user/or", func(write http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		age := request.URL.Query().Get("age")
		ageInt, _ := strconv.Atoi(age)

		resultUsersArray := make([]string, 0)

		for _, value := range mapUser {
			if name == value.Name || ageInt == value.Age {
				result := value.String()
				resultUsersArray = append(resultUsersArray, result)
			}
		}

		result := strings.Join(resultUsersArray, "\n")
		format.Text(write, 200, result)

		if len(resultUsersArray) < 1 { //если совпадения не найдены
			format.Text(write, 404, "No matches found!")
		}
	})

	//http://localhost:8080/map/user/friend //принимает спискок из id других пользователей

	mapUserWithFriends := make(map[int]UserWithFriends)

	router.Post("/user/friend", func(write http.ResponseWriter, request *http.Request) {
		id := request.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		name := request.FormValue("name")
		surname := request.FormValue("surname")
		age := request.FormValue("age")
		ageInt, _ := strconv.Atoi(age)
		sex := request.FormValue("sex")
		friend := request.FormValue("friend")
		friendArray := strings.Split(friend, ",")

		friendsArray := make([]int, 0)
		for _, elem := range friendArray {
			friendInt, _ := strconv.Atoi(elem)
			friendsArray = append(friendsArray, friendInt)
		}

		_, ok := mapUserWithFriends[idInt]

		if ok { //проверка на существующего пользователя
			result := fmt.Sprintf("User with id %d already exist!", idInt)
			format.Text(write, 404, result)
		} else {

			mapUserWithFriends[idInt] = UserWithFriends{idInt, name, surname, ageInt, sex, friendsArray}

			var result []string

			for _, value := range mapUserWithFriends {
				resultString := value.ToString()
				result = append(result, resultString)
			}

			resultUsersWithFriends := strings.Join(result, "\n")
			format.Text(write, 200, resultUsersWithFriends)
		}
	})

	//http://localhost:8080/map/user/{id}/friend  //находит в map других пользователей которые записаны в поле Friend
	router.Get("/user/{id}/friend", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		_, ok := mapUserWithFriends[idInt]

		if ok {
			userStruct := mapUserWithFriends[idInt]

			existIndx := make([]int, 0)
			deleteIndx := make([]int, 0)

			for _, valueID := range userStruct.Friend {

				_, ok := mapUserWithFriends[valueID]

				if ok {
					existIndx = append(existIndx, valueID) // составляем новый срез только с найденными пользователями
				} else { // составляем новый срез только с ненайденными пользователями
					deleteIndx = append(deleteIndx, valueID)
				}
			}

			for _, valueID := range existIndx {
				resultStruct := mapUserWithFriends[valueID]
				result := fmt.Sprintf("Friend:\nID: %d\nName: %s\nSurname: %s\n\n", resultStruct.ID, resultStruct.Name, resultStruct.Surname)
				format.Text(write, 200, result)
			}

			for _, valueID := range deleteIndx {
				result := fmt.Sprintf("Friend with ID %d was not found and was deleted.\n", valueID)
				format.Text(write, 200, result)
			}

			userStruct.Friend = existIndx
			resultStr := userStruct.ToString()
			mapUserWithFriends[idInt] = userStruct                       // ложим копию в оригинал
			result := fmt.Sprintf("\nUpdated user data:\n%s", resultStr) // выводим обновленные данные пользователя
			format.Text(write, 200, result)

		} else {
			result := fmt.Sprintf("User with id %d is not found!", idInt)
			format.Text(write, 404, result)
		}
	})
}

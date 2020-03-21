package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForArray(router chi.Router, format *render.Render) {

	userArray := make([]User, 0)

	//http://localhost:8080/array/user   создать список(массив) из пользователей и добавлять пользователя в этот список
	router.Post("/user", func(write http.ResponseWriter, request *http.Request) {
		id := request.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		name := request.FormValue("name")
		surname := request.FormValue("surname")
		age := request.URL.Query().Get("age")
		ageInt, _ := strconv.Atoi(age)
		sex := request.URL.Query().Get("sex")

		UserInfo := User{idInt, name, surname, ageInt, sex}
		userArray = append(userArray, UserInfo)

		var result []string
		for _, value := range userArray {
			resultString := value.String()
			//resultString := fmt.Sprintf("ID: %d\nName: %s\nSurname: %s\nAge: %d\nSex: %s\n", value.ID, value.Name, value.Surname, value.Age, value.Sex)
			result = append(result, resultString)
		}

		resultUsers := strings.Join(result, "\n")
		format.Text(write, 200, resultUsers)
	})

	//http://localhost:8080/array/user/last   вернуть полследнего пользователя из списка, в виде строки
	router.Get("/user/last", func(write http.ResponseWriter, request *http.Request) {

		len := len(userArray)
		if len < 1 {
			format.Text(write, 404, "User is not found.\n")
		} else {
			userLast := userArray[len-1] //struct
			result := userLast.String()
			//result := fmt.Sprintf("ID: %d\nName: %s\nSurname: %s\nAge: %d\nSex: %s\n", userLast.ID, userLast.Name, userLast.Surname, userLast.Age, userLast.Sex)
			format.Text(write, 200, result)
		}
	})

	//http://localhost:8080/array/user/first   вернуть первого пользователя из списка, в виде строки
	router.Get("/user/first", func(write http.ResponseWriter, request *http.Request) {

		len := len(userArray)
		if len < 1 {
			format.Text(write, 404, "User is not found.\n")
		} else {
			userFirst := userArray[0] //struct
			result := userFirst.String()
			format.Text(write, 200, result)
		}
	})

	//http://localhost:8080/array/user/count вернуть количество пользователей, которых мы сохранили
	router.Get("/user/count", func(write http.ResponseWriter, request *http.Request) {

		len := len(userArray)
		lenString := strconv.Itoa(len)

		format.Text(write, 200, lenString)
	})

	//http://localhost:8080/array/user/3   вернуть пользователя который лежит в ячейке массива с указанным id в path parameter
	router.Get("/user/{id}", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		len := len(userArray)
		if idInt > len {
			format.Text(write, 404, "User is not found.\n")
		} else {
			userID := userArray[idInt]
			result := userID.String()
			format.Text(write, 200, result)
		}
	})

	//http://localhost:8080/array/user/id/111   поиск пользователя по номеру ID
	router.Get("/user/id/{id}", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		searchUserByID := func(users []User, id int) (User, bool) { //функция
			result := User{}
			isFound := false
			for _, value := range userArray {
				if value.ID == idInt { // ищем по значению поля
					isFound = true
					result = value // ложим структуру, т.е. найденного пользователя
				}
			}

			return result, isFound // возвращается структура и булевое значение
		}

		user, isFound := searchUserByID(userArray, idInt) //вызов функции

		if isFound { //true
			strUser := user.String()
			format.Text(write, 200, strUser) //возвращаем данные найденного пользователя по ID
		} else {
			format.Text(write, 404, "User is not found.\n")
		}
	})

	//http://localhost:8080/array/user/last   удаление последнего добавленного пользователя
	router.Delete("/user/last", func(write http.ResponseWriter, request *http.Request) {

		len1 := len(userArray)

		if len1 < 1 {
			format.Text(write, 404, "The user list is empty, can't be deleted.\n")
		} else {
			userPenult := userArray[:len1]
			len2 := len(userPenult)
			index := len2 - 1
			result := fmt.Sprintf("User id %d has been deleted.", index)
			format.Text(write, 200, result) //возвращаем индекс в массиве удаленного пользователя
		}
	})

	//http://localhost:8080/array/user/id/3   удаление пользователя по переданному id
	router.Delete("/user/id/{id}", func(write http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		idInt, _ := strconv.Atoi(id)

		len := len(userArray)

		if len < 1 {
			format.Text(write, 404, "The user list is empty, can't be deleted.\n")
		} else {
			userDelArray1 := userArray[:idInt]
			userDelArray2 := userArray[idInt-1:]
			userDelArray1 = append(userDelArray1, userDelArray2...)
			index := idInt - 1
			result := fmt.Sprintf("User id %d has been deleted.", index)
			format.Text(write, 200, result) //возвращаем индекс в массиве удаленного пользователя
		}
	})

	//http://localhost:8080/array/user/id/num/111   удаление пользователя по номеру ID
	router.Delete("/user/id/num/{idNum}", func(write http.ResponseWriter, request *http.Request) {
		idNum := chi.URLParam(request, "idNum")
		idNumInt, _ := strconv.Atoi(idNum)

		searchUserByID := func(users []User, id int) (int, bool) {
			indx := 0
			isFound := false
			for i, value := range userArray {
				if value.ID == idNumInt {
					indx = i
					isFound = true
				}
			}
			return indx, isFound
		}

		indx, isFound := searchUserByID(userArray, idNumInt)
		if isFound {
			userArrayDel1 := userArray[:indx]
			userArrayDel2 := userArray[indx+1:]
			userArrayDel1 = append(userArrayDel1, userArrayDel2...)
			result := fmt.Sprintf("User id %d has been deleted.", idNumInt)
			format.Text(write, 200, result)
		} else {
			format.Text(write, 404, "User is not found.\n")
		}
	})

}

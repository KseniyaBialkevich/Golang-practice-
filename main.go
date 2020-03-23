package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func main() {
	format := render.New()
	router := chi.NewRouter()

	handlersForPractice(router, format)

	router.Route("/test", func(methodRouter chi.Router) {
		handlersForMethods(methodRouter, format)
	})

	handlersForURLShortening(router, format)
	handlersForURLTime(router, format)

	router.Route("/struct", func(structRouter chi.Router) {
		handlersForStruct(structRouter, format)
	})

	router.Route("/array", func(arrayRouter chi.Router) {
		handlersForArray(arrayRouter, format)
	})

	router.Route("/map", func(mapRouter chi.Router) {
		handlersForMap(mapRouter, format)
	})

	fmt.Println("Server is running!")
	http.ListenAndServe(":8080", router)
}

//общая структура, в которой хранятся данные пользователя
type User struct {
	ID      int
	Name    string
	Surname string
	Age     int
	Sex     string
}

//метод String(), форматирующий структуру User в тип string
func (arg User) String() string {
	result := fmt.Sprintf("ID: %d\nName: %s\nSurname: %s\nAge: %d\nSex: %s\n", arg.ID, arg.Name, arg.Surname, arg.Age, arg.Sex)
	return result
}

//структура, в которой хранятся данные пользователя с дополнительным полем Friend
type UserWithFriends struct {
	ID      int
	Name    string
	Surname string
	Age     int
	Sex     string
	Friend  []int //поле является списком идентbфикаторов других пользователей
}

//метод ToString(), форматирующий структуру UserWithFriends в тип string
func (arg UserWithFriends) ToString() string {
	result := fmt.Sprintf("ID: %d\nName: %s\nSurname: %s\nAge: %d\nSex: %s\nFriends: %d\n", arg.ID, arg.Name, arg.Surname, arg.Age, arg.Sex, arg.Friend)
	return result
}

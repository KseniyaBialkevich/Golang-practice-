package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForPractice(router chi.Router, format *render.Render) {

	// http://localhost:8080/add?a=1&b=2
	router.Get("/add", func(write http.ResponseWriter, request *http.Request) {
		a := request.URL.Query().Get("a")
		ai, _ := strconv.Atoi(a)

		b := request.URL.Query().Get("b")
		bi, _ := strconv.Atoi(b)

		sum := ai + bi
		//fmt.Fprintf(write, "Sum = %d", sum) // Fprintf пишет значения в write, a write это то куда отпр-ся результат
		//format.JSON(write, 200, sum)
		format.Text(write, 200, strconv.Itoa(sum))
	})

	//------------------------------------------------------------------------------------

	// http://localhost:8080/addextra
	count := 0
	router.Get("/addextra", func(write http.ResponseWriter, request *http.Request) {
		count++

		//fmt.Fprintf(write, "Count: %d", count)
		format.JSON(write, 200, count)
	})

	//------------------------------------------------------------------------------------

	// http: //localhost:8080/add?a=1&b=2
	router.Get("/add", func(write http.ResponseWriter, request *http.Request) {
		a := request.URL.Query().Get("a")
		ai, _ := strconv.Atoi(a)

		b := request.URL.Query().Get("b")
		bi, _ := strconv.Atoi(b)

		sum := ai + bi
		fmt.Fprintf(write, "Sum = %d", sum)
	})

	//------------------------------------------------------------------------------------

	// http://localhost:8080/addsum
	countt := 0
	router.Get("/addsum", func(write http.ResponseWriter, request *http.Request) {

		countt = countt + 1

		result := struct {
			Count int `json:"count"`
		}{countt}

		format.JSON(write, 200, result)
	})

	//------------------------------------------------------------------------------------

	// http://localhost:8080/vote?entity=1
	type AnimalsCounter struct {
		dogCounter, catCounter, gopherCounter int
	}

	animalCounter := AnimalsCounter{0, 0, 0}

	router.Get("/vote", func(write http.ResponseWriter, request *http.Request) {
		var entity string = request.URL.Query().Get("entity")

		if entity == "1" {
			animalCounter.dogCounter++
		} else if entity == "2" {
			animalCounter.catCounter++
		} else if entity == "3" {
			animalCounter.gopherCounter++
		}

		format.JSON(write, 200, "ok")
	})

	// http://localhost:8080/addCandidate?name=bird&id=4

	type VoteInfo struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	voteCounter := make(map[int]VoteInfo)

	// http://localhost:8080/result
	type AddCandidate struct {
		ID    int `json:"id"`
		Count int `json:"count"`
	}

	router.Get("/addCandidate", func(write http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		sID := request.URL.Query().Get("id")
		id, _ := strconv.Atoi(sID)

		voteInfo, exist := voteCounter[id]

		if exist { //true
			voteInfo.Count++
			voteCounter[id] = voteInfo
		} else { //false
			voteCounter[id] = VoteInfo{name, 1}
		}

		voteInfo, _ = voteCounter[id]

		count := voteInfo.Count

		result := AddCandidate{id, count}

		format.JSON(write, 200, result)
	})

	router.Get("/result", func(write http.ResponseWriter, request *http.Request) {

		result := "{\n"
		isFirst := true

		for _, voteInfo := range voteCounter {
			if isFirst {
				result += fmt.Sprintf(" %s: %d", voteInfo.Name, voteInfo.Count)
				isFirst = false
			} else {
				result += fmt.Sprintf(",\n %s: %d", voteInfo.Name, voteInfo.Count)
			}
		}

		result += "\n}"

		fmt.Fprintf(write, result)
		//format.JSON(write, 200, result)
	})
}

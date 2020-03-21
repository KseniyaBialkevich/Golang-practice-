package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

func handlersForMethods(router chi.Router, format *render.Render) {

	// The GET method requests ---------------------------------------------------------------------------

	//http://localhost:8080/test/get/value?value=my_hero   //обработчик, который возвращает переданное значение в value
	router.Get("/get/value", func(writer http.ResponseWriter, request *http.Request) {
		value := request.URL.Query().Get("value")
		format.Text(writer, 200, value)
		//my_hero
	})

	//http://localhost:8080/test/get/random   //обработчик, который возвращает случайное число
	router.Get("/get/random", func(writer http.ResponseWriter, request *http.Request) {
		tNow := time.Now().UnixNano()
		rand.Seed(tNow)
		randNum := rand.Intn(110)
		stringRandNum := fmt.Sprintf("%d", randNum)
		format.Text(writer, 200, stringRandNum)
		//~84
	})

	//http://localhost:8080/test/get/sum?value1=2&value2=3   //обработчик, переводит оба параметра в тип число, делает их суму в переменную sum и возвращает строку
	router.Get("/get/sum", func(writer http.ResponseWriter, request *http.Request) {
		value1 := request.URL.Query().Get("value1")
		intValue1, _ := strconv.Atoi(value1)
		value2 := request.URL.Query().Get("value2")
		intValue2, _ := strconv.Atoi(value2)
		sum := intValue1 + intValue2
		result := fmt.Sprintf("\"%d + %d = %d\"", intValue1, intValue2, sum)
		format.Text(writer, 200, result)
		//"2 + 3 = 5"
	})

	//http://localhost:8080/test/get/many_arg?arg1=apple&arg2=banana&arg3=orange&arg4=avocado&arg5=passion fruit   //возвращает последовательность, разделенную запятой
	router.Get("/get/many_arg", func(writer http.ResponseWriter, request *http.Request) {
		arg1 := request.URL.Query().Get("arg1")
		arg2 := request.URL.Query().Get("arg2")
		arg3 := request.URL.Query().Get("arg3")
		arg4 := request.URL.Query().Get("arg4")
		arg5 := request.URL.Query().Get("arg5")
		handRecord := fmt.Sprintf("%s, %s, %s, %s, %s", arg1, arg2, arg3, arg4, arg5) //apple, banana, orange, avocado, passion fruit
		libRecord := []string{arg1, arg2, arg3, arg4, arg5}
		libraryRecord := strings.Join(libRecord, ", ") //apple, banana, orange, avocado, passion fruit
		result := fmt.Sprintf("Hand record: %s;\nLibrary record: %s.", handRecord, libraryRecord)
		format.Text(writer, 200, result)
		//Hand record: apple, banana, orange, avocado, passion fruit;
		//Library record: apple, banana, orange, avocado, passion fruit.
	})

	//http://localhost:8080/test/get/sumti_sum?arg=1,2,3,4,5,6   //выделить из этого аргуметра все числа, сложить их и вывести все числа и сумму
	router.Get("/get/sumti_sum", func(writer http.ResponseWriter, request *http.Request) {
		arg := request.URL.Query().Get("arg")
		arrayString := strings.Split(arg, ",")
		len := len(arrayString)
		arrayInt := make([]int, len)
		for i := 0; i < len; i++ {
			arrayInt[i], _ = strconv.Atoi(arrayString[i])
		}
		sum := 0
		for i := 0; i < len; i++ {
			sum = sum + arrayInt[i]
		}
		seqNum := strings.Join(arrayString, ", ")    // 1, 2, 3, 4, 5, 6
		exercise := strings.Join(arrayString, " + ") //1 + 2 + 3 + 4 + 5 + 6
		result := fmt.Sprintf("%s\n%s = %d", seqNum, exercise, sum)
		format.Text(writer, 200, result)
		//1, 2, 3, 4, 5, 6
		//1 + 2 + 3 + 4 + 5 + 6 = 21
	})

	//http://localhost:8080/test/get/sumti_sum_rand_separator?args=1*2*3&separator=*   //числа разделены символом который указан в separator, вывести сумму
	router.Get("/get/sumti_sum_rand_separator", func(writer http.ResponseWriter, request *http.Request) {
		args := request.URL.Query().Get("args")
		separator := request.URL.Query().Get("separator")
		arrayArg := strings.Split(args, separator)
		len := len(arrayArg)
		arrayInt := make([]int, len)
		for i := 0; i < len; i++ {
			arrayInt[i], _ = strconv.Atoi(arrayArg[i])
		}
		sum := 0
		for i := 0; i < len; i++ {
			sum = sum + arrayInt[i]
		}
		exercise := strings.Join(arrayArg, " + ")
		result := fmt.Sprintf("%s = %d", exercise, sum)
		format.Text(writer, 200, result)
		//1 + 2 + 3 = 6
	})

	//http://localhost:8080/test/get/replace?args=cucumber,broccoli,tomato,avocado&separator=- // заменить запятую на переданный в separator символ и вернуть строку
	router.Get("/get/replace", func(writer http.ResponseWriter, request *http.Request) {
		args := request.URL.Query().Get("args")
		separator := request.URL.Query().Get("separator")
		arrayArg := strings.Split(args, ",")
		result := strings.Join(arrayArg, separator)
		format.Text(writer, 200, result)
		//cucumber broccoli tomato avocado
	})

	//http://localhost:8080/test/get/replace_by?args=1a*2b*3c&separator_orign=*&separator_new=-   //заменить separator_orign на переданный в separator_new символ
	router.Get("/get/replace_by", func(writer http.ResponseWriter, request *http.Request) {
		args := request.URL.Query().Get("args")
		separatorOrigin := request.URL.Query().Get("separator_orign")
		separatorNew := request.URL.Query().Get("separator_new")
		arrayArg := strings.Split(args, separatorOrigin)
		result := strings.Join(arrayArg, separatorNew)
		format.Text(writer, 200, result)
		//1a-2b-3c
	})

	// The POST method requests ---------------------------------------------------------------------------

	//http://localhost:8080/test/post/same   //обработчик принимает параметр с название value в теле запроса и возвращает его
	//value=Aloha!
	router.Post("/post/same", func(writer http.ResponseWriter, request *http.Request) {
		value := request.FormValue("value")
		format.Text(writer, 200, value)
		//Aloha!
	})

	//http://localhost:8080/test/post/concat   //принимает два параметра в теле запроса и возвращает их конкатенацию
	//val1=butter val2=fly
	router.Post("/post/concat", func(writer http.ResponseWriter, request *http.Request) {
		val1 := request.FormValue("val1")
		val2 := request.FormValue("val2")
		format.Text(writer, 200, val1+val2)
		//butterfly
	})

	//http://localhost:8080/test/post/sum   //переводит в числовой тип, делает сумму и возвращает строку
	//int1=3 int2=6
	router.Post("/post/sum", func(writer http.ResponseWriter, request *http.Request) {
		int1 := request.FormValue("int1")
		int2 := request.FormValue("int2")
		int1Int, _ := strconv.Atoi(int1)
		int2Int, _ := strconv.Atoi(int2)
		sum := int1Int + int2Int
		result := fmt.Sprintf("%s + %s = %d", int1, int2, sum)
		format.Text(writer, 200, result)
		//3 + 6 = 9
	})

	//http://localhost:8080/test/post/multi_for   //принимает параметры с одинаковым именем (arg), возвращает последовательность, разделенную запятыми
	//arg=Monday arg=Tuesday arg=Wednesday arg=Thursday arg=Friday
	router.Post("/post/multi_for", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		argsMap := request.Form //map[arg:[Monday Tuesday Wednesday Thursday Friday]]
		argsArray := argsMap["arg"]
		result := strings.Join(argsArray, ", ")
		format.Text(writer, 200, result)
		//Monday, Tuesday, Wednesday, Thursday, Friday
	})

	//http://localhost:8080/test/post/multi_for/many_keys   //принимает параметры с разными именами, возвращает последовательность, разделенную запятыми
	//arg1=Monday arg2=Tuesday arg3=Wednesday arg4=Thursday arg5=Friday
	router.Post("/post/multi_for", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		argsMap := request.Form //map[arg1:[Monday] arg2:[Tuesday] arg3:[Wednesday] arg4:[Thursday] arg5:[Friday]]
		argsArray := []string{}
		for _, arrayValue := range argsMap {
			for _, value := range arrayValue {
				argsArray = append(argsArray, value)
			}
		}
		result := strings.Join(argsArray, ", ")
		format.Text(writer, 200, result)
		//Monday, Tuesday, Wednesday, Thursday, Friday
	})

	// Path (Путь) -----------------------------------------------------------------------------------------

	//http://localhost:8080/test/path_parameter/Minsk   //возвращает значение которое было передано
	router.Get("/path_parameter/{paramPath}", func(writer http.ResponseWriter, request *http.Request) {
		paramPath := chi.URLParam(request, "paramPath")
		format.Text(writer, 200, paramPath)
		//Minsk
	})

	//http://localhost:8080/test/path_parameter/sum/30/69   //переводит параметры в число, делает их суму и возвращает строку
	router.Get("/path_parameter/sum/{paramPath1}/{paramPath2}", func(writer http.ResponseWriter, request *http.Request) {
		paramPath1 := chi.URLParam(request, "paramPath1")
		paramPath2 := chi.URLParam(request, "paramPath2")
		param1Int, _ := strconv.Atoi(paramPath1)
		param2Int, _ := strconv.Atoi(paramPath2)
		sum := param1Int + param2Int
		result := fmt.Sprintf("%d + %d = %d", param1Int, param2Int, sum)
		format.Text(writer, 200, result)
		//30 + 69 = 99
	})

	//http://localhost:8080/test/path_parameter/many_args/Who/is/John/Gold   //возвращает последовательность, разделенную пробелами
	router.Get("/path_parameter/many_args/{paramPath1}/{paramPath2}/{paramPath3}/{paramPath4}", func(writer http.ResponseWriter, request *http.Request) {
		paramsPath1 := chi.URLParam(request, "paramPath1")
		paramsPath2 := chi.URLParam(request, "paramPath2")
		paramsPath3 := chi.URLParam(request, "paramPath3")
		paramsPath4 := chi.URLParam(request, "paramPath4")
		paramsArray := []string{paramsPath1, paramsPath2, paramsPath3, paramsPath4}
		paramsString := strings.Join(paramsArray, " ")
		format.Text(writer, 200, paramsString+"?")
		//Who is John Gold?
	})

}

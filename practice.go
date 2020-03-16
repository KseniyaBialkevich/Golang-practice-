package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

var format = render.New()

func main() {
	router := chi.NewRouter()

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		result := "ok"
		format.Text(writer, 200, result)
	})

	mMap := make(map[string]string)
	//http://localhost:8080/process
	//?url=https://en.wikipedia.org/wiki/URL_shortening
	router.Post("/process", func(write http.ResponseWriter, request *http.Request) {
		url := request.FormValue("url")
		hash := md5.Sum([]byte(url))
		hashString := fmt.Sprintf("%x", hash)
		mMap[hashString] = url
		//mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"] = "https://en.wikipedia.org/wiki/URL_shortening"
		resultStringUrl := fmt.Sprintf("http://localhost:8080/%s", hashString)
		format.Text(write, 200, resultStringUrl)
	})

	// //http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e
	// router.Route("/{hash}", func(intRouter chi.Router) {
	// 	intRouter.Get("/", func(write http.ResponseWriter, request *http.Request) {
	// 		hashPath := chi.URLParam(request, "hash")
	// 		orignUrl, _ := mMap[hashPath]
	// 		format.Text(write, 200, orignUrl)
	// 	})
	// })

	//http://localhost:8080/3c4d0ce2967a743e5ee1f2c4cb31e29e
	router.Get("/{hash}", func(write http.ResponseWriter, request *http.Request) {
		hashPath := chi.URLParam(request, "hash") // 3c4d0ce2967a743e5ee1f2c4cb31e29e
		orignUrl, _ := mMap[hashPath]
		//"https://en.wikipedia.org/wiki/URL_shortening", _ := mMap["3c4d0ce2967a743e5ee1f2c4cb31e29e"]
		//format.Text(write, 200, orignUrl)
		http.Redirect(write, request, orignUrl, 301)
	})

	// ***********************************************************************
	// The GET method requests

	//http://localhost:8080/test/get/value?value=my_hero
	router.Get("/test/get/value", func(writer http.ResponseWriter, request *http.Request) {
		value := request.URL.Query().Get("value")
		format.Text(writer, 200, value)
		//my_hero
	})

	//http://localhost:8080/test/get/random
	router.Get("/test/get/random", func(writer http.ResponseWriter, request *http.Request) {
		tNow := time.Now().UnixNano()
		rand.Seed(tNow)
		randNum := rand.Intn(110)
		stringRandNum := fmt.Sprintf("%d", randNum)
		format.Text(writer, 200, stringRandNum)
		//~84
	})

	//http://localhost:8080/test/get/sum?value1=2&value2=3
	router.Get("/test/get/sum", func(writer http.ResponseWriter, request *http.Request) {
		value1 := request.URL.Query().Get("value1")
		intValue1, _ := strconv.Atoi(value1)
		value2 := request.URL.Query().Get("value2")
		intValue2, _ := strconv.Atoi(value2)
		sum := intValue1 + intValue2
		result := fmt.Sprintf("\"%d + %d = %d\"", intValue1, intValue2, sum)
		format.Text(writer, 200, result)
		//"2 + 3 = 5"
	})

	//http://localhost:8080/test/get/many_arg?arg1=apple&arg2=banana&arg3=orange&arg4=avocado&arg5=passion fruit
	router.Get("/test/get/many_arg", func(writer http.ResponseWriter, request *http.Request) {
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

	//http://localhost:8080/test/get/sumti_sum?arg=1,2,3,4,5,6
	router.Get("/test/get/sumti_sum", func(writer http.ResponseWriter, request *http.Request) {
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

	//http://localhost:8080/test/get/sumti_sum_rand_separator?args=1*2*3&separator=*
	router.Get("/test/get/sumti_sum_rand_separator", func(writer http.ResponseWriter, request *http.Request) {
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

	//http://localhost:8080/test/get/replace?args=cucumber,broccoli,tomato,avocado&separator=-
	router.Get("/test/get/replace", func(writer http.ResponseWriter, request *http.Request) {
		args := request.URL.Query().Get("args")
		separator := request.URL.Query().Get("separator")
		arrayArg := strings.Split(args, ",")
		result := strings.Join(arrayArg, separator)
		format.Text(writer, 200, result)
		//cucumber broccoli tomato avocado
	})

	//http://localhost:8080/test/get/replace_by?args=1a*2b*3c&separator_orign=*&separator_new=-
	router.Get("/test/get/replace_by", func(writer http.ResponseWriter, request *http.Request) {
		args := request.URL.Query().Get("args")
		separatorOrigin := request.URL.Query().Get("separator_orign")
		separatorNew := request.URL.Query().Get("separator_new")
		arrayArg := strings.Split(args, separatorOrigin)
		result := strings.Join(arrayArg, separatorNew)
		format.Text(writer, 200, result)
		//1a-2b-3c
	})

	//************************************************************************
	// The POST method requests

	//http://localhost:8080/test/post/same
	//?value=Aloha!
	router.Post("/test/post/same", func(writer http.ResponseWriter, request *http.Request) {
		value := request.FormValue("value")
		format.Text(writer, 200, value)
	})
	//Aloha!

	//http://localhost:8080/test/post/concat
	//?val1=butter&val2=fly
	router.Post("/test/post/concat", func(writer http.ResponseWriter, request *http.Request) {
		val1 := request.FormValue("val1")
		val2 := request.FormValue("val2")
		format.Text(writer, 200, val1+val2)
		//butterfly
	})

	//http://localhost:8080/test/post/sum
	//?int1=3&int2=6
	router.Post("/test/post/sum", func(writer http.ResponseWriter, request *http.Request) {
		int1 := request.FormValue("int1")
		int2 := request.FormValue("int2")
		int1Int, _ := strconv.Atoi(int1)
		int2Int, _ := strconv.Atoi(int2)
		sum := int1Int + int2Int
		result := fmt.Sprintf("%s + %s = %d", int1, int2, sum)
		format.Text(writer, 200, result)
		//3 + 6 = 9
	})

	//http://localhost:8080/test/post/multi_for
	//?arg=Monday&arg=Tuesday&arg=Wednesday&arg=Thursday&arg=Friday
	router.Post("/test/post/multi_for", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		argsMap := request.Form //map[arg:[Monday Tuesday Wednesday Thursday Friday]]
		argsArray := argsMap["arg"]
		result := strings.Join(argsArray, ", ")
		format.Text(writer, 200, result)
		//Monday, Tuesday, Wednesday, Thursday, Friday
	})

	//http://localhost:8080/test/post/multi_for/many_keys
	//?arg1=Monday&arg2=Tuesday&arg3=Wednesday&arg4=Thursday&arg5=Friday
	router.Post("/test/post/multi_for", func(writer http.ResponseWriter, request *http.Request) {
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

	//************************************************************************
	// Path (Путь)

	//http://localhost:8080/test/path_parameter/Minsk
	router.Get("/test/path_parameter/{paramPath}", func(writer http.ResponseWriter, request *http.Request) {
		paramPath := chi.URLParam(request, "paramPath")
		format.Text(writer, 200, paramPath)
		//Minsk
	})

	//http://localhost:8080/test/path_parameter/sum/30/69
	router.Get("/test/path_parameter/sum/{paramPath1}/{paramPath2}", func(writer http.ResponseWriter, request *http.Request) {
		paramPath1 := chi.URLParam(request, "paramPath1")
		paramPath2 := chi.URLParam(request, "paramPath2")
		param1Int, _ := strconv.Atoi(paramPath1)
		param2Int, _ := strconv.Atoi(paramPath2)
		sum := param1Int + param2Int
		result := fmt.Sprintf("%d + %d = %d", param1Int, param2Int, sum)
		format.Text(writer, 200, result)
		//30 + 69 = 99
	})

	//http://localhost:8080/test/path_parameter/many_args/Who/is/John/Gold
	router.Get("/test/path_parameter/many_args/{paramPath1}/{paramPath2}/{paramPath3}/{paramPath4}", func(writer http.ResponseWriter, request *http.Request) {
		paramsPath1 := chi.URLParam(request, "paramPath1")
		paramsPath2 := chi.URLParam(request, "paramPath2")
		paramsPath3 := chi.URLParam(request, "paramPath3")
		paramsPath4 := chi.URLParam(request, "paramPath4")
		paramsArray := []string{paramsPath1, paramsPath2, paramsPath3, paramsPath4}
		paramsString := strings.Join(paramsArray, " ")
		format.Text(writer, 200, paramsString+"?")
		//Who is John Gold?
	})

	//************************************************************************

	http.ListenAndServe(":8080", router)
}


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	endpnts "gitgub.com/my/repo/endpoints"
)

func home_page(w http.ResponseWriter, r *http.Request) {

	// проверка, что переход только по разрешенным эндпоинтам
	if !endpnts.Check_endpoint(r) {
		infoResp := endpnts.InfoResponse{Msg: "Forbidden"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)                  // на все некорректные uri отвечаем 403
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		return
	}

	//http.NotFound(w, r)

	//fmt.Println(r.URL)

	var show_data endpnts.ShowData

	var s string

	if key := r.FormValue("key_get"); key != "" { // если хотим получить данные
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/get_key?key=%s", key), nil) // то делаем запрос

		resp, _ := http.DefaultClient.Do(req)

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		s = string(body)
		show_data.GetKeyMsg = s
		fmt.Printf("Get response: %s", body)

	}

	if key := r.FormValue("key_set"); key != "" { // если хотим добавить данные

		val := r.FormValue("value")

		fmt.Println(key, val)

		data := endpnts.Response{Key: key, Value: val}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(data)

		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest("POST", "http://localhost:8080/set_key", &buf) // то делаем запрос
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		s = string(body)
		show_data.SetKeyMsg = s
		fmt.Printf("Set response: %s", body)

	}

	if key := r.FormValue("key_del"); key != "" { // если хотим добавить данные

		data := endpnts.Response{Key: key}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(data)

		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest("DELETE", "http://localhost:8080/del_key", &buf) // то делаем запрос
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		s = string(body)

		show_data.DelKeyMsg = s

		//fmt.Println(body)
		fmt.Printf("Delete response: %s", body)

	}

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil { // если получили ошибку то выводим на страницу какую
		fmt.Fprintf(w, err.Error())
	}

	if s != "" {
		tmpl.Execute(w, show_data)
	} else {
		tmpl.Execute(w, nil)
	}

}

func without_interface(w http.ResponseWriter, r *http.Request) {
	// проверка, что переход только по разрешенным эндпоинтам
	if !endpnts.Check_endpoint(r) { //
		infoResp := endpnts.InfoResponse{Msg: "Forbidden"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)                  // на все некорректные uri отвечаем 403
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		return
	}
	fmt.Fprintf(w, "<h1>App is running!</h1>")
}

func main() {

	//fmt.Println(len(os.Args), os.Args[1]) // считываем аргумент при запуске

	if len(os.Args) > 1 && os.Args[1] == "yes" { // если да, то запускаем с веб-интерфейсом
		fmt.Println("App is running with web-interface!")
		http.HandleFunc("/", home_page)
	} else { // если не переданы аргументы или передан другой, то запуск только для curl запросов
		fmt.Println("App is running!")
		http.HandleFunc("/", without_interface)
	}

	endpnts.HandleRequest()

}

package endpoints

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/redis/go-redis/v9"
)

type Config struct { // структура для конфига
	Pass    string // пароль редиса
	Address string // ip редиса
}

func get_conf() Config {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	password := viper.GetString("pass")
	address := viper.GetString("address")

	return Config{Pass: password, Address: address}
}

var addr string = get_conf().Address // передача параметров из конфига
var pass string = get_conf().Pass

type Response struct {
	Key   string `json: "key"`
	Value string `json: "value"`
}

type InfoResponse struct {
	Msg string `json: "msg"`
}

type ShowData struct { // структура для вывода информации на главной странице
	GetKeyMsg string
	SetKeyMsg string
	DelKeyMsg string
}

func get_key(w http.ResponseWriter, r *http.Request) { // функция обработки эндпоинта /get_key
	fmt.Printf("got /get_key connection \n")

	key := r.URL.Query().Get("key")

	cert, err := tls.LoadX509KeyPair("tests/tls/redis.crt", "tests/tls/redis.key")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile("tests/tls/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := redis.NewClient(&redis.Options{
		//Addr: "redis:6379",
		Addr:     addr,
		Password: pass,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			//ServerName: "Generic-cert",
			ServerName: "localhost",
		},
		DB: 0,
	})

	ctx := context.Background()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		infoResp := InfoResponse{Msg: "There is no such key."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)                  // если такого ключа нет в редисе, то 404
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		//http.NotFound(w, r)

		return
	}

	res := Response{Key: key, Value: val}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)             // ресурс найден 200 ок
	json.NewEncoder(w).Encode(res) // возвращаем json

}

func set_key(w http.ResponseWriter, r *http.Request) { // функция обработки эндпоинта /set_key
	fmt.Printf("got /set_key connection \n")

	var res Response
	err := json.NewDecoder(r.Body).Decode(&res) // получаем json инфу переданную в body и записываем в res

	//fmt.Println(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // ошибка отправленных данных
		infoResp := InfoResponse{Msg: "Invalid key or value type! Must be string."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		return
	}

	cert, err := tls.LoadX509KeyPair("tests/tls/redis.crt", "tests/tls/redis.key")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile("tests/tls/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			ServerName:   "localhost",
		},
		DB: 0,
	})

	ctx := context.Background()

	// fmt.Println(res.Key, res.Value) полученные ключ и значение

	response_string := "Key setted."

	val, err := client.Get(ctx, res.Key).Result() // проверяем есть ли ключ в БД
	if val != "" {                                // если есть то пишем, что обновили
		response_string = "Key updated."
	}

	err = client.Set(ctx, res.Key, res.Value, 0).Err() // устанавливаем значение по ключу
	if err != nil {
		panic(err)
	}

	var empty_value string

	if res.Value == "" { // если не передано значение ключа, то добавляем сообщение об этом
		empty_value = "Notice, that value was empty!"
	}

	infoResp := InfoResponse{Msg: fmt.Sprintf("%s %s", response_string, empty_value)} // формируем json ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(infoResp) // возвращаем json
}

func del_key(w http.ResponseWriter, r *http.Request) { // функция обработки эндпоинта /del_key
	fmt.Printf("got /del_key connection \n")

	var res Response
	err := json.NewDecoder(r.Body).Decode(&res) // получаем json инфу переданную в body и записываем в res

	//fmt.Println(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // ошибка отправленных данных
		infoResp := InfoResponse{Msg: "Invalid key type! Must be string."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		return
	}

	cert, err := tls.LoadX509KeyPair("tests/tls/redis.crt", "tests/tls/redis.key")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile("tests/tls/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		//Addr:     "redis:6379",
		Password: pass,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			//ServerName: "Generic-cert",
			ServerName: "localhost",
		},
		DB: 0,
	})

	ctx := context.Background()
	key := res.Key

	_, err = client.Get(ctx, key).Result()

	if err != nil { // если такого ключа нет в БД, то ошибку
		infoResp := InfoResponse{Msg: "There is no such key."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(infoResp) // возвращаем json
		return
	}

	err = client.Del(ctx, res.Key).Err() // удаляем данные по ключу
	if err != nil {
		panic(err)
	}

	infoResp := InfoResponse{Msg: "Key sucessfully deleted!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(infoResp) // возвращаем json

}

func Check_endpoint(r *http.Request) bool {
	url := html.EscapeString(r.URL.Path)
	if url != "/" {
		return false
	}
	return true
}

func HandleRequest() {
	http.HandleFunc("/get_key", get_key)
	http.HandleFunc("/set_key", set_key)
	http.HandleFunc("/del_key", del_key)
	//htpp.HandleFunc("/get_handle", get_handle)
	// http.HandleFunc("/create", create)
	// http.HandleFunc("/save_data", save_data)
	http.ListenAndServe(":8080", nil)
}

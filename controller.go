package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/cagge/lvl0/getter"
	_ "github.com/lib/pq"
	stan "github.com/nats-io/stan.go"
)

type Content struct {
	Title string
	Text  string
}

var Memory [][]byte

func home(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Привет из Snippetbox"))
	ts, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
func control() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

	// http.HandleFunc("/", generateResponse)
	// err := http.ListenAndServe("localhost:8080", nil)

	// if err != nil {
	// 	log.Fatalln("Случилась какая-то ошибка:", err)
	// }

	sc, err := stan.Connect("test-cluster", "clientID", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return err
	}
	go func() {
		for {
			jsonByte, _ := getter.RunData()
			sc.Publish("order", []byte(jsonByte))
			Memory = append(Memory, jsonByte)
			time.Sleep(3 * time.Second)
		}

	}()
	_, err = sc.Subscribe("order", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(Memory[0]))
		table, err := getter.Under(m.Data)
		if err != nil {
			return
		}
		//Memory = append(m.Data)
		err = connectBD(table)
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}

	select {}
	return nil
}

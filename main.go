package main

import (
	"flag"
	"fmt"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	//只允许访问根路径
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	//只允许GET请求
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	// 获取传入的参数
	addr := flag.String("addr", "0.0.0.0:9991", "http service address")
	flag.Parse()
	hub := NewHub()
	go hub.Run()
	//注册每种请求对应的处理函数
	http.HandleFunc("/", serveHome)
	// 接收http请求，根据入参通知指定的client
	http.HandleFunc("/upload", func(rw http.ResponseWriter, r *http.Request) {
		send(hub, rw, r)
	})
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(hub, rw, r)
	})
	fmt.Println("service start on " + *addr)
	//如果启动成功，该行会一直阻塞，hub.run()会一直运行
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Printf("start http service error: %s\n", err)
	}
}

//go run main.go --port 3434

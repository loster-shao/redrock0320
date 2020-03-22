package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("start now")
	router := mux.NewRouter()//创建路由
	//ws控制器不断处理管道数据，进行同步数据(就是加个协程）
	go h.run()
	router.HandleFunc("/ws", myws)//指定ws回调函数
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("err:", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

//websocket  连接结构体
type connection struct {
	ws   *websocket.Conn  //ws连接器
	sc   chan []byte      //管道
	data *Data            //数据
}

//websocket升级
var wu = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true }}

func myws(w http.ResponseWriter, r *http.Request) {
	ws, err := wu.Upgrade(w, r, nil)//获取ws对象
	if err != nil {
		return
	}
	//创建连接对象去做事情
	//初始化连接对象
	c := &connection{sc: make(chan []byte, 256), ws: ws, data: &Data{}}
	//在ws中注册一下
	h.register <- c
	go c.writer()//写入
	c.reader()//读出
	defer func() {
		c.data.Type = "logout"
		user_list = del(user_list, c.data.User)//用户列表删除
		c.data.UserList = user_list
		c.data.Content = c.data.User
		data_b, _ := json.Marshal(c.data)//数据序列化，让所有人看到某人下线
		h.broadcsat <- data_b
		h.register <- c
	}()
}

func (c *connection) writer() {
	//从管道遍历数据
	for message := range c.sc {
		c.ws.WriteMessage(websocket.TextMessage, message)  //数据写出
	}
	c.ws.Close()  //关闭                                                                                                           -
}

var user_list = []string{}

func (c *connection) reader() {
	for {//不断读websocket
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			h.register <- c
			break//如果读不了，则删除
		}
		json.Unmarshal(message, &c.data)//读取数据
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content//弹出窗口
			c.data.From = c.data.User
			user_list = append(user_list, c.data.User)//登录后，将用户加入到用户列表
			c.data.UserList = user_list//每个用户都登录了的列表
			data_b, _ := json.Marshal(c.data)//数据序列化
			h.broadcsat <- data_b
		case "user":// 普通状态
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)//序列化
			h.broadcsat <- data_b
		case "logout":
			c.data.Type = "logout"
			user_list = del(user_list, c.data.User)
			data_b, _ := json.Marshal(c.data)//数据序列化，所有人都应该看到某人下线
			h.broadcsat <- data_b
			h.register <- c
		default:
			fmt.Println("========default================")
		}
	}
}

func del(slice []string, user string) []string {
	count := len(slice)//严谨判断
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}//定义新的返回切片
	//删除传入切片中的指定用户，其他用户放到新的切片
	for i := range slice {
		//利用索引删用户
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}

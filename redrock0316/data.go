package main

//将连接中传输的数据抽象出对象
type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`//用户
	From     string   `json:"from"`//哪个用户发言
	Type     string   `json:"type"`//标识信息类型 登录 握手
	Content  string   `json:"content"`//传输内容
	UserList []string `json:"user_list"`//用户列表
}

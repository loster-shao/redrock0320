package main

import "encoding/json"

var h = hub{
	connection: make(map[*connection]bool),//应该是是否连接
	unregister: make(chan *connection),    //销毁请求
	broadcsat:  make(chan []byte),         //从连接器发送的信息
	register:   make(chan *connection),    //注册请求
}

//抽象ws连接器
//处理ws中逻辑
type hub struct {
	connection map[*connection]bool  //是否注册连接器
	broadcsat chan []byte  //从连接器发送的信息
	register chan *connection  //从连接器注册请求
	unregister chan *connection  //销毁请求
}

func (h *hub) run() {
	//监听数据管道，在后端不断处理管道数据
	for {//通过不同的数据管道，处理不同的逻辑
		select {
		//注册
		case c := <-h.register:
			h.connection[c] = true                //标志注册了
			c.data.Ip = c.ws.RemoteAddr().String()//组装data数据
			c.data.Type = "handshake"             //更新类型 表示不知道这是什么类型
			c.data.UserList = user_list           //更新用户列表
			data_b, _ := json.Marshal(c.data)     //继续序列化成byte
			c.sc <- data_b                        //放入数据管道
		//登出
		case c := <-h.unregister:
			//if函数判断map是否含有需删的数据
			if _, ok := h.connection[c]; ok {
				delete(h.connection, c)            //删除注销连接
				close(c.sc)                        //关闭管道
			}
		//数据channel
		case data := <-h.broadcsat:
			//处理数据流转，将数据同步到所有用户（即插入所有用户的管道）
			//c是具体的每个连接
			for c := range h.connection {
				select {
				case c.sc <- data:
				default:
					//防止死循环
					delete(h.connection, c)
					close(c.sc)
				}
			}
		}
	}
}

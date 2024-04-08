package main

type Hub struct {
	clients    map[*Client]bool //维护所有的client
	broadcast  chan []byte      //广播消息
	register   chan *Client     //注册
	unregister chan *Client     //注销
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte), //同步管道，确保hub消息不堆积，同时多个client给hub发数据会阻塞
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			//client上线，注册
			hub.clients[client] = true
		case client := <-hub.unregister:
			//查询当前client是否存在
			if _, exists := hub.clients[client]; exists {
				//注销client 通道
				close(client.send)
				//删除注销的client
				delete(hub.clients, client)
			}
		case msg := <-hub.broadcast:
			//将message广播给每一位client
			for client := range hub.clients {
				select {
				case client.send <- msg:
				//异常client处理
				default:
					close(client.send)
					//删除异常的client
					delete(hub.clients, client)
				}
			}
		}

	}
}

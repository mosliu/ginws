package wsutils

// 代码 参考 https://blog.csdn.net/dipolar/article/details/51532231
type Hub struct {
    // 注册了的连接器
    Connections map[*Client]bool

    // 从连接器中发入的信息
    Broadcast chan []byte

    // 从连接器中注册请求
    Register chan *Client

    // 从连接器中注销请求
    Unregister chan *Client
}

var CommonHub = Hub{
   Broadcast:   make(chan []byte),
   Register:    make(chan *Client),
   Unregister:  make(chan *Client),
   Connections: make(map[*Client]bool),
}

func (h *Hub) Run() {
    for {
        select {
        case c := <-h.Register:
            h.Connections[c] = true
        case c := <-h.Unregister:
            if _, ok := h.Connections[c]; ok {
                delete(h.Connections, c)
                close(c.SendChan)
            }
        case m := <-h.Broadcast:
            for c := range h.Connections {
                select {
                case c.SendChan <- m:
                default:
                    delete(h.Connections, c)
                    close(c.SendChan)
                }
            }
        }
    }
}

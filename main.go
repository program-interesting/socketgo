package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024 * 10,
		// 写入存储空间大小
		WriteBufferSize: 1024 * 10,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 链接列表
var memeberList map[string]*websocket.Conn

// 检测离开
func salir() {
	for {
		for k, v := range memeberList {
			if err := v.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				// 发送失败一次
				v.Close()
				delete(memeberList, k)
				fmt.Println("离开了", k)
				fmt.Println("当前连接人数:", len(memeberList))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wbsCon *websocket.Conn
		err    error
		data   []byte
	)
	// 完成http应答，在httpheader中放下如下参数
	wbsCon, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级WS出错", err)
		return // 获取连接失败直接返回
	} else {
		memeberList[wbsCon.RemoteAddr().String()] = wbsCon
		// delete(memeberList, wbsCon.RemoteAddr().String())
		fmt.Println("来了新链接", wbsCon.RemoteAddr().String())
		fmt.Println("当前连接人数:", len(memeberList))
	}

	for {
		// 只能发送Text, Binary 类型的数据,下划线意思是忽略这个变量.
		if _, data, err = wbsCon.ReadMessage(); err != nil {
			goto ERR // 跳转到关闭连接
		}
		fmt.Println(wbsCon.RemoteAddr().String(), string(data[:]))

		// []byte(string)
		if err = wbsCon.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR // 发送消息失败，关闭连接
		}
	}
ERR:
	// 关闭连接
	wbsCon.Close()
}

func main() {
	fmt.Println("开始运行")
	memeberList = make(map[string]*websocket.Conn)
	go salir()
	// 当有请求访问ws时，执行此回调方法
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe("0.0.0.0:1777", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

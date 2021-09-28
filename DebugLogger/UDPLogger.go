package DebugLogger

import (
	"bufio"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/newServer/setup"
	"net"
	"strings"
	"time"
)

func workerLogger(socket net.Conn) {
	defer socket.Close()
	socket.SetReadDeadline(time.Time{})
	reader := bufio.NewReader(socket)
	logger.Info.Printf("Новое соединение для логирования %s", socket.RemoteAddr().String())
	for {
		c, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("При чтении лога с %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		c = strings.Replace(c, "\n", "", -1)
		c = strings.Replace(c, "\r", "", -1)
		fmt.Printf("%s-%s\n", socket.RemoteAddr().String(), string(c))
	}
}
func ListenUDP() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", setup.Set.CommServer.PortDebug))
	if err != nil {
		logger.Error.Print("Debug Logger %s", err.Error())
		return
	}
	fmt.Println("DebugLogger ready")
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка входа по порту %d", setup.Set.CommServer.PortDevice)
			continue
		}
		go workerLogger(socket)
	}
}

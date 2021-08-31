package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/extcon"
	"github.com/ruraomsk/newServer/setup"
	"net"
	"sync"
	"time"
)

var devices map[int]*DeviceControl
var mutex sync.Mutex

func secondServer(stop chan interface{}) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", setup.Set.CommServer.PortDevice))
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %d", setup.Set.CommServer.PortDevice)
		stop <- 1
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка входа по порту %d", setup.Set.CommServer.PortDevice)
			continue
		}
		go workerListenDevice(socket)
	}

}
func MainServer(stop chan interface{}) {
	devices = make(map[int]*DeviceControl)
	go secondServer(stop)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", setup.Set.CommServer.Port))
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %d", setup.Set.CommServer.Port)
		stop <- 1
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка входа по порту %d", setup.Set.CommServer.Port)
			continue
		}
		go workerSendDevice(socket)
	}
}
func workerListenDevice(socket net.Conn) {
	//_ = socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	//_ = socket.SetWriteDeadline(time.Now().Add(time.Second * 10))

	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	logger.Info.Printf("Новое соединение для прослушивания %s", socket.RemoteAddr().String())
	//читаем что это за устройство

	c, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Error.Printf("При чтении настроек второго канала с %s %s", socket.RemoteAddr().String(), err.Error())
		_ = socket.Close()
		return
	}
	c = c[:len(c)-1]
	logger.Info.Printf("Первая строка резервного %s", string(c))
	var md DeviceInfo
	err = json.Unmarshal(c, &md)
	if err != nil {
		logger.Error.Printf("При конвертации от %s строки %s ошибка %s", socket.RemoteAddr().String(), string(c), err.Error())
		_ = socket.Close()
		return
	}
	logger.Info.Printf("Первый  резервного %v", md)

	mutex.Lock()
	dev, is := devices[md.ID]
	if !is {
		logger.Error.Printf("Устройство %d неверная логика подключения канала с %s ", md.ID, socket.RemoteAddr().String())
		_ = socket.Close()
		mutex.Unlock()
		return
	}
	dev.socketSecond = socket
	dev.full = true
	go deviceToServer(dev)
	mutex.Unlock()
	c = append(newMessageServerInfo().Adopt(), '\n')
	logger.Info.Printf("Отправляем по резервному %s", string(c))
	_, err = writer.Write(c)
	if err != nil {
		logger.Error.Printf("Устройство %d при передаче на %s %s", md.ID, socket.RemoteAddr().String(), err.Error())
		dev.stop <- 1
		return
	}
}
func workerSendDevice(socket net.Conn) {
	context, err := extcon.NewContext(fmt.Sprintf("md%s", socket.RemoteAddr().String()))
	if err != nil {
		return
	}

	_ = socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	_ = socket.SetWriteDeadline(time.Now().Add(time.Second * 10))
	reader := bufio.NewReader(socket)
	logger.Info.Printf("Новое соединение %s", socket.RemoteAddr().String())
	for {
		//читаем что это за устройство

		c, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Error.Printf("При чтении настроек с %s %s", socket.RemoteAddr().String(), err.Error())
			_ = socket.Close()
			return
		}
		//logger.Info.Printf("Первая строка основного %s",string(c))
		c = c[:len(c)-1]
		var md DeviceInfo
		err = json.Unmarshal(c, &md)
		if err != nil {
			logger.Error.Printf("При конвертации от %s строки %s ошибка %s", socket.RemoteAddr().String(), string(c), err.Error())
			_ = socket.Close()
			return
		}
		//logger.Info.Printf("Первый  основного %v",md)
		mutex.Lock()
		old, is := devices[md.ID]
		if is {
			old.closeAll()
			delete(devices, md.ID)
		}
		dev := newDevice(md, socket)
		devices[md.ID] = dev
		mutex.Unlock()
		go serverToDevice(dev)
		dev.chanToMain <- newMessageServerInfo().Adopt()
		for {
			select {
			case <-dev.timer.C:
				//Прошло много времени со времени последнего обмена с устройством
				dev.closeAll()
				logger.Info.Printf("Очень долго не отвечает %d", dev.DeviceInfo.ID)
			case <-context.Done():
				dev.closeAll()
				logger.Info.Printf("Приказано умереть %d", dev.DeviceInfo.ID)
				return
			case <-dev.stop:
				dev.closeAll()
				logger.Info.Printf("Остановлено устройство %d", dev.DeviceInfo.ID)
				return
			case <-dev.errorChan:
				dev.closeAll()
				logger.Info.Printf("Устройство %d отключено из-за ошибок", dev.DeviceInfo.ID)
				return
			}
		}

	}
}

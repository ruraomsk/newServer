package server

import (
	"bufio"
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
var delay = 5
var cmds = []string{"give-me-SetupDK", "give-me-TCPMain", "give-me-TCPSec", "give-me-Cameras", "give-me-OneMonth,1",
	"give-me-OneWeek,2", "give-me-OneDay,3", "give-me-OnePK,4"} //,"give-me-Cameras"
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
		go workerListenDevice(&socket)
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
func workerListenDevice(socket *net.Conn) {
	//_ = socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	//_ = socket.SetWriteDeadline(time.Now().Add(time.Second * 10))

	reader := bufio.NewReader(*socket)
	writer := bufio.NewWriter(*socket)
	logger.Info.Printf("Новое соединение для второго канала %s", (*socket).RemoteAddr().String())

	c, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Error.Printf("При чтении настроек второго канала с %s %s", (*socket).RemoteAddr().String(), err.Error())
		_ = (*socket).Close()
		return
	}
	c = c[:len(c)-1]
	logger.Info.Printf("Первая строка второго %s", string(c))
	md := getInfoConnect(string(c))
	mutex.Lock()
	dev, is := devices[md.ID]
	if !is {
		logger.Error.Printf("Устройство %d неверная логика подключения канала с %s ", md.ID, (*socket).RemoteAddr().String())
		_ = (*socket).Close()
		mutex.Unlock()
		return
	}
	dev.socketSecond = socket
	dev.full = true
	go deviceToServer(dev)
	mutex.Unlock()

	s := fmt.Sprintf("ok,127,%d,%d\n", delay, time.Now().Unix())
	logger.Info.Printf("Отправляем по резервному %s", s)
	_, _ = writer.WriteString(s)
	err = writer.Flush()
	if err != nil {
		logger.Error.Printf("При передаче на %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
		dev.errorChan <- 1
		return
	}

	if err != nil {
		logger.Error.Printf("Устройство %d при передаче на %s %s", md.ID, (*socket).RemoteAddr().String(), err.Error())
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

		logger.Info.Printf("Первый  основного %s", string(c))
		md := getInfoConnect(string(c))
		deleteDevice(md.ID)
		dev := newDevice(md, &socket)
		devices[md.ID] = dev
		go serverToDevice(dev)
		dev.chanToMain <- fmt.Sprintf("ok,127,%d,%d\n", delay, time.Now().Unix())
		for {
			select {
			case <-dev.chanFromMain:
				//logger.Info.Printf("Пришло от главного %d %s", dev.DeviceInfo.ID, mess)
			case <-dev.chanFromSecond:
				//logger.Info.Printf("Пришло от резервного %d %s", dev.DeviceInfo.ID, mess)
				dev.chanToSecond <- "ok"
			case <-dev.timer1.C:
				logger.Info.Printf("KEEP Alive для %d", dev.DeviceInfo.ID)
				dev.chanToMain <- giveMeStatus()
				for _, cmd := range cmds {
					//logger.Info.Printf("Команда %s для %d",cmd, dev.DeviceInfo.ID)
					dev.chanToMain <- cmd
				}
				//Прошло много времени со времени последнего обмена с устройством
				//deleteDevice(md.ID)
			case <-dev.timer2.C:
				logger.Info.Printf("Очень долго не отвечает %d", dev.DeviceInfo.ID)
				//deleteDevice(md.ID)
				//return
			case <-context.Done():
				deleteDevice(md.ID)
				logger.Info.Printf("Приказано умереть %d", dev.DeviceInfo.ID)
				return
			case <-dev.stop:
				deleteDevice(md.ID)
				logger.Info.Printf("Остановлено устройство %d", dev.DeviceInfo.ID)
				return
			case <-dev.errorChan:
				deleteDevice(md.ID)
				logger.Info.Printf("Устройство %d отключено из-за ошибок", dev.DeviceInfo.ID)
				return
			}
		}

	}
}
func deleteDevice(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	d, is := devices[id]
	if is {
		d.closeAll()
		delete(devices, id)
	}
}

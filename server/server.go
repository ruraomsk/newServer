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
var delay = 15
var cmds = []string{"give-me-SetupDK", "give-me-TCPMain", "give-me-TCPSec", "give-me-Cameras", "give-me-OneMonth,1",
	"give-me-OneWeek,1", "give-me-OneDay,1", "give-me-OnePK,1", "give-me-OnePhase,0"} //,"give-me-Cameras"
var sets = []string{
	"change-begin",
	"set-you-OnePhase,0={\"num\":0,\"work\":true}",
	"set-you-SetupDK={\"dkn\":10234,\"tmaxf\":8,\"tminf\":5,\"tminmax\":7,\"dktype\":10,\"extn\":15678,\"tprom\":6,\"preset\":false}",
	"set-you-Cameras={\"cameras\":[{\"id\":1,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":2,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":3,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":4,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":5,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":6,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":7,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":8,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":9,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":10,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":11,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":12,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":13,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":14,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":15,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8},{\"id\":16,\"ip\":\"192.168.115.168\",\"login\":\"admin\",\"password\":\"admin\",\"port\":8441,\"chanels\":8}]}",
	"set-you-OneMonth,1={\"num\":1,\"days\":[1,2,3,4,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]}",
	"set-you-OneWeek,1={\"num\":1,\"days\":[1,2,3,4,5,6,7]}",
	"set-you-OneDay,1={\"num\":1,\"count\":2,\"lines\":[{\"npk\":1,\"hour\":0,\"min\":0},{\"npk\":2,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0},{\"npk\":0,\"hour\":0,\"min\":0}]}",
	"set-you-OnePK,1={\"pk\":1,\"tpu\":1,\"razlen\":true,\"tc\":0,\"shift\":0,\"twot\":false,\"sts\":[{\"line\":1,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":2,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":3,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":4,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":5,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":6,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":7,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":8,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":9,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":10,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":11,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":12,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":13,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":14,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":15,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":16,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":17,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":18,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":19,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":20,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":21,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":22,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":23,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0},{\"line\":24,\"start\":0,\"num\":0,\"tf\":0,\"stop\":0,\"plus\":false,\"trs\":false,\"dt\":0}]}",
	"change-end",
}

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
	fmt.Printf("%s Новое соединение %s\n", time.Now().String(), socket.RemoteAddr().String())
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
			case <-dev.timer1.C:
				logger.Info.Printf("KEEP Alive для %d", dev.DeviceInfo.ID)
				dev.chanToMain <- giveMeStatus()
				for _, cmd := range cmds {
					//logger.Info.Printf("Команда %s для %d",cmd, dev.DeviceInfo.ID)
					if !dev.deleted {
						dev.chanToMain <- cmd
					}
				}
				for _, set := range sets {
					if !dev.deleted {
						dev.chanToMain <- set
					}
				}
			case <-dev.timer2.C:
				logger.Info.Printf("Очень долго не отвечает %d", dev.DeviceInfo.ID)
				dev.stop <- 1
			case <-context.Done():
				dev.stop <- 1
				logger.Info.Printf("Приказано умереть %d", dev.DeviceInfo.ID)
			case <-dev.stop:
				deleteDevice(md.ID)
				logger.Info.Printf("Остановлено устройство %d", dev.DeviceInfo.ID)
				fmt.Printf("Остановлено устройство %d\n", dev.DeviceInfo.ID)
			case <-dev.errorChan:
				logger.Info.Printf("Устройство %d отключено из-за ошибок", dev.DeviceInfo.ID)
				fmt.Printf("%s Устройство %d отключено из-за ошибок\n", dev.DeviceInfo.ID)
				dev.stop <- 1
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

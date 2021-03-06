package server

import (
	"github.com/ruraomsk/TLServer/logger"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DeviceControl struct {
	mutex          sync.Mutex
	socketMain     *net.Conn        //Соединение по которому сервер шлет свои команды устройству
	socketSecond   *net.Conn        //Соединение по которому устройства шлет свои внеочередные сообщения
	stop           chan interface{} //По этому каналу приходит команды остановиться
	serverCMD      chan interface{} //По этому каналу приходят команды от системы
	errorChan      chan interface{} //Если программы ввода/вывода находят ошибку
	chanToMain     chan string
	chanFromMain   chan string
	chanToSecond   chan string
	chanFromSecond chan string
	interval       time.Duration //Интервал времени для проверки состояния устройства
	DeviceInfo     DeviceInfo
	timer1         *time.Timer //Для вычисления не выхода на связь
	timer2         *time.Timer //Для вычисления не выхода на связь

	full    bool
	deleted bool
}

func (d *DeviceControl) restartTimer1() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.timer1.Stop()
	d.timer1 = time.NewTimer(d.interval * time.Second)
}
func (d *DeviceControl) restartTimer2() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.timer2.Stop()
	d.timer2 = time.NewTimer((d.interval + 5) * time.Second)
}

func newDevice(devinfo DeviceInfo, socket *net.Conn) *DeviceControl {
	dev := new(DeviceControl)
	dev.chanToMain = make(chan string)
	dev.chanToSecond = make(chan string, 100)
	dev.chanFromMain = make(chan string, 100)
	dev.chanFromSecond = make(chan string, 100)
	dev.errorChan = make(chan interface{})
	dev.stop = make(chan interface{})
	dev.serverCMD = make(chan interface{})
	dev.DeviceInfo = devinfo
	dev.interval = time.Duration(delay)
	dev.socketMain = socket
	dev.full = false
	dev.deleted = false
	dev.timer1 = time.NewTimer(dev.interval * time.Second)
	dev.timer2 = time.NewTimer(10 * dev.interval * time.Second)
	return dev
}
func (d *DeviceControl) closeAll() {
	logger.Info.Printf("Обмен с устройством %d %s прекращается", d.DeviceInfo.ID, d.DeviceInfo.Type)
	if !d.deleted {
		close(d.chanToMain)
		close(d.chanToSecond)
		close(d.chanFromMain)
		close(d.chanFromSecond)
		//_=d.socketMain.SetWriteDeadline(time.Now().Add(1 * time.Second))
		//_,_=d.socketMain.Write([]byte("stop\n"))
		_ = (*d.socketMain).Close()
		if d.full {
			//_=d.socketSecond.SetWriteDeadline(time.Now().Add(1 * time.Second))
			//_,_=d.socketSecond.Write([]byte("stop\n"))
			_ = (*d.socketSecond).Close()
		}
		d.deleted = true

	}
	logger.Info.Printf("Обмен с устройством %d %s прекращен", d.DeviceInfo.ID, d.DeviceInfo.Type)
}

type OnBoard struct {
	GPRS GPRSInfo `json:"gprs,omitempty"`
	GPS  GPSInfo  `json:"gps,omitempty"`
	ETH  ETHInfo  `json:"eth,omitempty"`
	USB  USBInfo  `json:"usb,omitempty"`
}
type SoftWare struct {
	Version string `json:"version"`
}
type Secret struct {
	Key string `json:"key"`
}
type DebugInfo struct {
	Make bool   `json:"make"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type GPRSInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type GPSInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type ETHInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type USBInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}

func getInfoConnect(s string) DeviceInfo {
	ls := strings.Split(s, ",")
	if len(ls) != 3 {
		return DeviceInfo{ID: 0}
	}
	id, _ := strconv.Atoi(ls[1])
	return DeviceInfo{ID: id, Type: ls[2]}
}

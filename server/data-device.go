package server

import (
	"github.com/ruraomsk/TLServer/logger"
	"net"
	"sync"
	"time"
)

type DeviceControl struct {
	mutex          sync.Mutex
	socketMain     net.Conn         //Соединение по которому сервер шлет свои команды устройству
	socketSecond   net.Conn         //Соединение по которому устройства шлет свои внеочередные сообщения
	stop           chan interface{} //По этому каналу приходит команды остановиться
	serverCMD      chan interface{} //По этому каналу приходят команды от системы
	errorChan      chan interface{} //Если программы ввода/вывода находят ошибку
	chanToMain     chan MessageServer
	chanFromMain   chan MessageDevice
	chanToSecond   chan MessageServer
	chanFromSecond chan MessageDevice
	interval       time.Duration //Интервал времени для проверки состояния устройства
	DeviceInfo     DeviceInfo
	timer          *time.Timer //Для вычисления не выхода на связь
	full           bool
	deleted        bool
}

func (d *DeviceControl) restartTimer() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.timer.Stop()
	d.timer = time.NewTimer(d.interval * time.Second)
}
func newDevice(devinfo DeviceInfo, socket net.Conn) *DeviceControl {
	dev := new(DeviceControl)
	dev.chanToMain = make(chan MessageServer)
	dev.chanToSecond = make(chan MessageServer)
	dev.chanFromMain = make(chan MessageDevice)
	dev.chanFromSecond = make(chan MessageDevice)
	dev.errorChan = make(chan interface{})
	dev.stop = make(chan interface{})
	dev.serverCMD = make(chan interface{})
	dev.DeviceInfo = devinfo
	dev.interval = 4 * 60
	dev.socketMain = socket
	dev.full = false
	dev.deleted = true
	dev.timer = time.NewTimer(dev.interval * time.Second)
	return dev
}
func (d *DeviceControl) closeAll() {
	logger.Info.Printf("Обмен с устройством %d %s прекращается", d.DeviceInfo.ID, d.DeviceInfo.Type)
	if !d.deleted {
		close(d.chanToMain)
		close(d.chanToSecond)
		close(d.chanFromMain)
		close(d.chanFromSecond)
		_ = d.socketMain.Close()
		if d.full {
			_ = d.socketSecond.Close()
		}
		d.deleted = true
		mutex.Lock()
		delete(devices, d.DeviceInfo.ID)
		mutex.Unlock()

	}
	logger.Info.Printf("Обмен с устройством %d %s прекращен", d.DeviceInfo.ID, d.DeviceInfo.Type)
}

type DeviceInfo struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` //ETH или GPRS
	OnBoard   OnBoard   `json:"onboard,omitempty"`
	Software  SoftWare  `json:"software,omitempty"`
	Secret    Secret    `json:"secret,omitempty"`
	DebugInfo DebugInfo `json:"debug,omitempty"`
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

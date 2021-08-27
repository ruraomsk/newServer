package server

import (
	"bufio"
	"github.com/ruraomsk/TLServer/logger"
	"time"
)

func serverToDevice(dev *DeviceControl) {
	_ = dev.socketMain.SetReadDeadline(time.Now().Add(time.Second * 10))
	_ = dev.socketMain.SetWriteDeadline(time.Now().Add(time.Second * 10))
	reader := bufio.NewReader(dev.socketMain)
	writer := bufio.NewWriter(dev.socketMain)
	for {
		select {
		case <-dev.stop:
			return
		case s := <-dev.chanToMain:
			s = append(s, '\n')
			_, err := writer.Write(s)
			if err != nil {
				logger.Error.Printf("При передаче на %d по адресу %s %s", dev.DeviceInfo.ID, dev.socketMain.RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			dev.restartTimer()
			c, err := reader.ReadBytes('\n')
			if err != nil {
				logger.Error.Printf("При приеме от %d по адресу %s %s", dev.DeviceInfo.ID, dev.socketMain.RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			c = c[:len(c)-1]
			dev.chanFromMain <- c

		}
	}
}
func deviceToServer(dev *DeviceControl) {
	_ = dev.socketSecond.SetReadDeadline(time.Time{}) //Без таймаута чтение
	_ = dev.socketMain.SetWriteDeadline(time.Now().Add(time.Second * 10))
	reader := bufio.NewReader(dev.socketMain)
	writer := bufio.NewWriter(dev.socketMain)
	for {
		c, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Error.Printf("При приеме от %d по адресу %s %s", dev.DeviceInfo.ID, dev.socketMain.RemoteAddr().String(), err.Error())
			dev.errorChan <- 1
			return
		}
		c = c[:len(c)-1]
		dev.chanFromSecond <- c
		dev.restartTimer()
		select {
		case <-dev.stop:
			return
		case s := <-dev.chanToSecond:
			s = append(s, '\n')
			_, err = writer.Write(s)
			if err != nil {
				logger.Error.Printf("При передаче на %d по адресу %s %s", dev.DeviceInfo.ID, dev.socketMain.RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			dev.restartTimer()
		}
	}
}

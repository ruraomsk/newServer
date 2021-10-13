package server

import (
	"bufio"
	"github.com/ruraomsk/TLServer/logger"
	"time"
)

func serverToDevice(dev *DeviceControl) {
	reader := bufio.NewReader(*dev.socketMain)
	writer := bufio.NewWriter(*dev.socketMain)
	for {
		select {
		case <-dev.stop:
			return
		case s := <-dev.chanToMain:
			time.Sleep(250 * time.Millisecond)
			logger.Info.Printf("Отправляем по главному %s\n", s)
			_ = (*dev.socketMain).SetWriteDeadline(time.Now().Add(time.Second * 10))
			_, _ = writer.WriteString(s)
			_, _ = writer.WriteString("\n")
			err := writer.Flush()
			if err != nil {
				logger.Error.Printf("При передаче на %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			_ = (*dev.socketMain).SetReadDeadline(time.Now().Add(time.Second * 10))
			c, err := reader.ReadBytes('\n')
			if err != nil {
				logger.Error.Printf("При приеме ответа от %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			logger.Info.Printf("Получили по главному %s", string(c))
			dev.chanFromMain <- string(c[:len(c)-1])
			dev.restartTimer1()

		}
	}
}
func deviceToServer(dev *DeviceControl) {
	reader := bufio.NewReader(*dev.socketSecond)
	writer := bufio.NewWriter(*dev.socketSecond)
	for {
		_ = (*dev.socketSecond).SetReadDeadline(time.Time{}) //Без таймаута чтение
		c, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Error.Printf("При приеме от второго %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketSecond).RemoteAddr().String(), err.Error())
			dev.errorChan <- 1
			return
		}
		logger.Info.Printf("Получили по резервному %s", string(c))
		dev.chanFromSecond <- string(c[:len(c)-1])
		dev.restartTimer2()
		select {
		case <-dev.stop:
			return
		case s := <-dev.chanToSecond:
			logger.Info.Printf("Отправляем по второму %s", s)
			_ = (*dev.socketSecond).SetWriteDeadline(time.Now().Add(time.Second * 10))
			_, _ = writer.WriteString(s)
			_, _ = writer.WriteString("\n")
			err = writer.Flush()
			if err != nil {
				logger.Error.Printf("При передаче второго  %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
				dev.errorChan <- 1
				return
			}
			dev.restartTimer2()
		}
	}
}

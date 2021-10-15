package server

import (
	"bufio"
	"github.com/ruraomsk/TLServer/logger"
	"strings"
	"time"
)

func serverToDevice(dev *DeviceControl) {
	reader := bufio.NewReader(*dev.socketMain)
	writer := bufio.NewWriter(*dev.socketMain)
	for {
		select {
		case s := <-dev.chanToMain:
			for {
				time.Sleep(250 * time.Millisecond)
				if dev.deleted {
					return
				}
				logger.Info.Printf("-> %s", s)
				_ = (*dev.socketMain).SetWriteDeadline(time.Now().Add(time.Second * 10))
				_, _ = writer.WriteString(s)
				_, _ = writer.WriteString("\n")
				err := writer.Flush()
				if err != nil {
					logger.Error.Printf("При передаче на %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
					dev.errorChan <- 1
					return
				}
				var c []byte
				count := 1
				for {
					if dev.deleted {
						return
					}
					_ = (*dev.socketMain).SetReadDeadline(time.Now().Add(time.Second * 10))
					c, err = reader.ReadBytes('\n')
					if err != nil {
						count--
						logger.Error.Printf("При приеме ответа от %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketMain).RemoteAddr().String(), err.Error())
						if !strings.Contains(err.Error(), "i/o timeout") {
							dev.errorChan <- 1
							return
						}
						if count < 0 {
							dev.errorChan <- 1
							return
						}
						continue
					}
					break
				}
				if dev.deleted {
					return
				}
				logger.Info.Printf("<- %s", string(c))
				if strings.HasPrefix(string(c), "repeat again") {
					continue
				}
				dev.restartTimer1()
				dev.chanFromMain <- string(c[:len(c)-1])
				break
			}

		}
	}
}
func deviceToServer(dev *DeviceControl) {
	reader := bufio.NewReader(*dev.socketSecond)
	for {
		_ = (*dev.socketSecond).SetReadDeadline(time.Time{}) //Без таймаута чтение
		c, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Error.Printf("При приеме от второго %d по адресу %s %s", dev.DeviceInfo.ID, (*dev.socketSecond).RemoteAddr().String(), err.Error())
			dev.errorChan <- 1
			return
		}
		if dev.deleted {
			return
		}
		logger.Info.Printf("<= %s", string(c))
		dev.restartTimer2()
		if dev.deleted {
			return
		}
		dev.chanFromSecond <- string(c[:len(c)-1])
	}
}

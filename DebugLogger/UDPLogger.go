package DebugLogger

import (
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/newServer/setup"
	"net"
	"strings"
)

func ListenUDP() {
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", setup.Set.CommServer.PortDebug))
	if err != nil {
		logger.Error.Print("Debug Logger %s", err.Error())
		return
	}
	fmt.Println("DebugLogger ready")
	buffer := make([]byte, 512)
	defer pc.Close()
	for {
		for i := range buffer {
			buffer[i] = 0
		}
		_, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			continue
		}
		s := strings.Split(addr.String(), ":")
		instring := s[0] + ":" + string(buffer)
		instring = strings.Replace(instring, "\n\r", "", 1)
		fmt.Println(instring)
	}
}

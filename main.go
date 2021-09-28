package main

import (
	"fmt"
	"github.com/ruraomsk/newServer/server"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/extcon"
	"github.com/ruraomsk/newServer/setup"
	//pprof init

	_ "net/http/pprof"
)

var err error

//Секция инициализации программы
func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFile("config.toml", &setup.Set); err != nil {
		fmt.Println("Can't load config file - ", err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err = logger.Init("log")
	if err != nil {
		fmt.Println("Error opening logger subsystem ", err.Error())
		return
	}
	logger.Info.Println("Start new-server work...")
	fmt.Println("Start new-server work...")
	extcon.BackgroundInit()

	//go DebugLogger.ListenUDP()

	stop := make(chan interface{})
	go server.MainServer(stop)

	//ready := make(chan interface{})
	extcon.BackgroundWork(time.Duration(1*time.Second), stop)
	logger.Info.Println("Exit ag-server working...")
	fmt.Println("\nExit ag-server working...")
}

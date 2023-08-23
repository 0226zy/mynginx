package core

import (
	"fmt"
	"net"
	"strconv"
)

// NgxCycle: 全局信息
type NgxCycle struct {
	Modules []NgxModule
	// TODO: 对象内存池

	// pool * Pool
	Listening []*NgxListen
	Conf      *NgxConf
}

type NgxListen struct {
	port     int
	listener net.Listener
}

var cycle *NgxCycle

func CreateGlobalCycle() *NgxCycle {

	cycle = &NgxCycle{}
	return cycle
}

func GetGlobalCycle() *NgxCycle {
	return cycle
}

func (cycle *NgxCycle) CountHttpModule() int {
	count := 0
	for _, module := range cycle.Modules {
		if module.Type() != ENgxHttpModule {
			continue
		}
		count++
	}
	return count
}

func (cycle *NgxCycle) OpenListeningSockekts() {

	var err error
	for _, listen := range cycle.Listening {
		for tries := 0; tries < 5; tries++ {
			listen.listener, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(listen.port))
			if err != nil {
				fmt.Printf("listen port:%d failed:%v\n", listen.port, err)
				continue
			}
			break
		}
	}

}

func (cycle *NgxCycle) AddListening(port int) {
	cycle.Listening = append(cycle.Listening, &NgxListen{
		port:     port,
		listener: nil,
	})
}

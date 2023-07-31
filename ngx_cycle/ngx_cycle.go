package ngx_cycle

import (
	"github.com/0226zy/mynginx/ngx_conf"
)

// NgxCycle: 全局信息
type NgxCycle struct {

	// TODO: 对象内存池

	// pool * Pool
}

func NgxInitCycle() *NgxCycle {
	cycle := &NgxCycle{}

	// create conf
	NgxConf := ngx_conf.CreateNgxConf()

	// parse
	NgxConf.Parse()

	return cycle

}

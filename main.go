package main

import (
	"github.com/0226zy/mynginx/core"
)

func main() {

	// 初始化日志

	// 内存池初始化

	// 命令行参数处理

	// init cycle
	core.CreateGlobalCycle()
	ngxInitCycle()

	// 启动子进程

}

func ngxInitCycle() {

	//ngxCycle = core.GetGlobalCycle()

	// 执行所有模块的 create_conf

	// 解析 命令行参数配置

	// 解析配置文件
	ngxConf := core.CreateNgxConf()
	ngxConf.ParseFile("./conf.d/nginx.conf")

	// 执行所有模块的 init_conf

}

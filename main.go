package main

import (
	"fmt"
	"os"

	"github.com/0226zy/mynginx/core"
	httpModules "github.com/0226zy/mynginx/http/modules"
)

func main() {

	pid := os.Getpid()
	processCycle := core.NewProcessCycle()
	processCycle.ParseArgs()

	// 初始化日志

	// 内存池初始化

	// 命令行参数处理

	// init cycle
	core.CreateGlobalCycle()
	ngxInitCycle(processCycle)
	fmt.Printf("pid:%d end init cycle\n", pid)

	// 进程模型处理
	processCycle.ProcessCycle()
}

func ngxInitCycle(processCycle *core.ProcessCycle) {

	ngxCycle := core.GetGlobalCycle()
	ngxCycle.Modules = append(ngxCycle.Modules, core.GetNgxCoreModule())
	ngxCycle.Modules = append(ngxCycle.Modules, core.GetNgxEventsModule())
	ngxCycle.Modules = append(ngxCycle.Modules, core.GetNgxEventCoreModule())
	ngxCycle.Modules = append(ngxCycle.Modules, core.GetNgxHttpModule())
	ngxCycle.Modules = append(ngxCycle.Modules, core.GetNgxHttpCoreModule())
	ngxCycle.Modules = append(ngxCycle.Modules, httpModules.GetNgxHttpRewriteModule())

	// 执行所有模块的 create_conf

	// 解析 命令行参数配置

	// 解析配置文件
	ngxConf := core.CreateNgxConf()
	ngxConf.ModuleType = core.ENgxCoreModule
	ngxConf.CmdType = core.TNgxMainConf
	ngxConf.ParseFile("./conf.d/nginx.conf")

	// 执行所有模块的 init_conf

	// listen
	if processCycle.IsMasterProcess() {
		ngxCycle.OpenListeningSockekts()
	}
	ngxCycle.Conf = ngxConf

}

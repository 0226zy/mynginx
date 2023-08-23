package core

import (
	"fmt"
	"strconv"
)

// NgxCoreModule 基础模块
type NgxCoreModule struct {
	moduleType int
	moduleName string
	commands   []*NgxModuleCommand
	moduleConf ngxCoreModuleConf
}

type ngxCoreModuleConf struct {
	MasterProcess bool
	WorkProcess   int
	Daemon        bool
}

func WithWorkProcess(conf *ngxCoreModuleConf) CmdHandler {

	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {

		fmt.Printf("ngx_core_module Set daemonOpt ,args:%v\n", ngxConf.Args)
		num, err := strconv.Atoi(ngxConf.Args[1])
		if err != nil {
			panic("opt to int failed:" + ngxConf.Args[1])
		}
		conf.WorkProcess = num
	}
}

func WithMasterProcess(conf *ngxCoreModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("ngx_core_module Set masterPorcessOpt ,args:%v\n", ngxConf.Args)
		if ngxConf.Args[1] == "off" {
			conf.MasterProcess = false
		}
		if ngxConf.Args[1] == "on" {
			conf.MasterProcess = true
		}
	}
}

func WithDaemon(conf *ngxCoreModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		// TODO
		fmt.Printf("ngx_core_module Set WorkProcess ,args:%v\n", ngxConf.Args)
		if ngxConf.Args[1] == "off" {
			conf.Daemon = false
		}
		if ngxConf.Args[1] == "on" {
			conf.Daemon = true
		}
	}
}

var coreModule NgxCoreModule

func init() {

	coreModule = NgxCoreModule{
		moduleType: ENgxCoreModule,
		moduleName: "ngx_core_module",
		moduleConf: ngxCoreModuleConf{true, 0, true},
	}
	coreModule.commands = []*NgxModuleCommand{
		NewNgxModuleCommand("worker_processes", TNgxMainConf|TNgxConfFlag, WithWorkProcess(&coreModule.moduleConf)),
		NewNgxModuleCommand("master_process", TNgxMainConf|TNgxConfTake1, WithMasterProcess(&coreModule.moduleConf)),
		NewNgxModuleCommand("daemon", TNgxMainConf|TNgxDirectConf|TNgxConfTake1, WithDaemon(&coreModule.moduleConf))}
}

func (module *NgxCoreModule) GetCtx() interface{} {
	return nil
}

func (module *NgxCoreModule) InitMaster()  {}
func (module *NgxCoreModule) InitModule()  {}
func (module *NgxCoreModule) InitProcess() {}

func (module *NgxCoreModule) Index() int {
	return 0
}

func (module *NgxCoreModule) GetCommands() []*NgxModuleCommand { return module.commands }
func (module *NgxCoreModule) Name() string {
	return module.moduleName
}

func (module *NgxCoreModule) Type() int {
	return module.moduleType
}

func GetNgxCoreModule() NgxModule {
	return &coreModule
}

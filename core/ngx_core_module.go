package core

import "fmt"

const (
	ENgxMainConf int = 0
)

// NgxCoreModule 基础模块
type NgxCoreModule struct {
	moduleType int
	moduleName string
	commands   []NgxModuleCommand
}

// 模块配置设置 command
type workProcessOpt struct {
}

type masterPorcessOpt struct{}

type daemonOpt struct{}

func (opt *daemonOpt) Type() int {
	return ENgxMainConf
}

func (opt *daemonOpt) Name() string {
	return "daemon"
}
func (opt *masterPorcessOpt) Type() int {
	return ENgxMainConf
}

func (opt *masterPorcessOpt) Name() string {
	return "master_process"
}

func (opt *workProcessOpt) Type() int {
	return ENgxMainConf
}

func (opt *workProcessOpt) Name() string {
	return "worker_processes"
}

func (opt *daemonOpt) Set(conf *NgxConf, command *NgxModuleCommand, moudleConf *NgxModuleConf) {
	// TODO
	fmt.Printf("ngx_core_module Set daemonOpt ,args:%v\n", conf.Args)
}

func (opt *masterPorcessOpt) Set(conf *NgxConf, command *NgxModuleCommand, moudleConf *NgxModuleConf) {
	// TODO
	fmt.Printf("ngx_core_module Set masterPorcessOpt ,args:%v\n", conf.Args)
}

func (opt *workProcessOpt) Set(conf *NgxConf, command *NgxModuleCommand, moudleConf *NgxModuleConf) {
	// TODO
	fmt.Printf("ngx_core_module Set WorkProcess ,args:%v\n", conf.Args)
}

var coreModule NgxCoreModule

func init() {

	coreModule = NgxCoreModule{
		moduleType: ENgxCoreModule,
		moduleName: "ngx_core_module",
		commands: []NgxModuleCommand{
			&workProcessOpt{},
			&masterPorcessOpt{},
			&daemonOpt{},
		},
	}
}

func (module *NgxCoreModule) InitMaster()  {}
func (module *NgxCoreModule) InitModule()  {}
func (module *NgxCoreModule) InitProcess() {}

func (module *NgxCoreModule) Index() int {
	return 0
}

func (module *NgxCoreModule) GetCommands() []NgxModuleCommand { return module.commands }
func (module *NgxCoreModule) Name() string {
	return module.moduleName
}

func (module *NgxCoreModule) Type() int {
	return module.moduleType
}

func GetNgxCoreModule() NgxModule {
	return &coreModule
}

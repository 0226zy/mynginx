package core

import "fmt"

var eventsModule NgxEventsModule
var eventsModuleConf ngxEventsModuleConf

func init() {

	eventsModuleConf = ngxEventsModuleConf{}
	eventsModule = NgxEventsModule{
		moduleType: ENgxCoreModule,
		moduleName: "ngx_events_module",
		commands: []*NgxModuleCommand{
			NewNgxModuleCommand(
				"events",
				TNgxMainConf|TNgxConfBlock|TNgxConfNoArgs,
				WithEventsBlock(&eventsModuleConf)),
		},
	}
}

// NgxEventsModule: event 模块
type NgxEventsModule struct {
	moduleName string
	moduleType int
	commands   []*NgxModuleCommand
}

type ngxEventsModuleConf struct {
	WorkConnections int
}

// WithEventsBlock: 处理nginx.conf 中的 events 后面的配置块
func WithEventsBlock(conf *ngxEventsModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {

		// TODO:初始化 event 的其它 模块
		// 继续读取 events 后面的配置块
		preModuleType := ngxConf.ModuleType
		preCmdType := ngxConf.CmdType

		ngxConf.ModuleType = ENgxEventModule
		ngxConf.CmdType = TNgxEventConf
		fmt.Printf("Handle conf:Events Block args:%v\n", ngxConf.Args)
		ngxConf.Args = []string{}

		// 继续读取后面的配置
		ngxConf.ParseFile("")
		ngxConf.ModuleType = preModuleType
		ngxConf.CmdType = preCmdType

		// TODO:执行其它 event 模块的 init_conf

	}
}

func (module *NgxEventsModule) InitMaster()  {}
func (module *NgxEventsModule) InitModule()  {}
func (module *NgxEventsModule) InitProcess() {}

func (module *NgxEventsModule) Index() int {
	return 0
}

func (module *NgxEventsModule) GetCommands() []*NgxModuleCommand { return module.commands }
func (module *NgxEventsModule) Name() string {
	return module.moduleName
}

func (module *NgxEventsModule) Type() int {
	return module.moduleType
}

func GetNgxEventsModule() NgxModule {
	return &eventsModule
}

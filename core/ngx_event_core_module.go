package core

import (
	"fmt"
	"strconv"
)

var eventCoreModule NgxEventCoreModule

func init() {

	eventCoreModule = NgxEventCoreModule{
		moduleType: ENgxEventModule,
		moduleName: "ngx_event_core_module",
		commands: []*NgxModuleCommand{
			NewNgxModuleCommand(
				"worker_connections",
				TNgxEventConf|TNgxConfTake1,
				WithWokrerConnections(&eventsModuleConf)),
		},
	}
}

type NgxEventCoreModule struct {
	moduleName string
	moduleType int
	commands   []*NgxModuleCommand
}

func (module *NgxEventCoreModule) InitMaster()  {}
func (module *NgxEventCoreModule) InitModule()  {}
func (module *NgxEventCoreModule) InitProcess() {}

func (module *NgxEventCoreModule) Index() int {
	return 0
}

func (module *NgxEventCoreModule) GetCommands() []*NgxModuleCommand { return module.commands }
func (module *NgxEventCoreModule) Name() string {
	return module.moduleName
}

func (module *NgxEventCoreModule) Type() int {
	return module.moduleType
}

func GetNgxEventCoreModule() NgxModule {
	return &eventCoreModule
}

func WithWokrerConnections(conf *ngxEventsModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("conf handle,events block: worker_connnections args:%v\n", ngxConf.Args)
		num, err := strconv.Atoi(ngxConf.Args[1])
		if err != nil {
			panic("opt to int failed:" + ngxConf.Args[1])
		}
		conf.WorkConnections = num
	}
}

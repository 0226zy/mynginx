package core

import "fmt"

var httpModule NgxHttpModule
var httpModuleConf ngxHttpModuleConf

func init() {

	httpModuleConf = ngxHttpModuleConf{}

	httpModule = NgxHttpModule{
		moduleType: ENgxCoreModule,
		moduleName: "ngx_http_module",
		commands: []*NgxModuleCommand{
			NewNgxModuleCommand(
				"http",
				TNgxMainConf|TNgxMainConf|TNgxConfNoArgs,
				WithHttpBlock(&httpModuleConf)),
		},
	}
}

// NgxHttpModule: http 模块
type NgxHttpModule struct {
	moduleName string
	moduleType int
	commands   []*NgxModuleCommand
}

type NgxHttpConfCtx struct {
	mainConfs []*ngxHttpCoreMainConf
	srvConfs  []*ngxHttpCoreSrvConf
	locConfs  []*ngxHttpCoreLocConf
}

type ngxHttpCoreMainConf struct{}
type ngxHttpCoreSrvConf struct{}
type ngxHttpCoreLocConf struct{}

type ngxHttpModuleConf struct {
}

func (module *NgxHttpModule) InitMaster()  {}
func (module *NgxHttpModule) InitModule()  {}
func (module *NgxHttpModule) InitProcess() {}

func (module *NgxHttpModule) Index() int {
	return 0
}

func (module *NgxHttpModule) GetCommands() []*NgxModuleCommand { return module.commands }
func (module *NgxHttpModule) Name() string {
	return module.moduleName
}

func (module *NgxHttpModule) Type() int {
	return module.moduleType
}

func GetNgxHttpModule() NgxModule {
	return &httpModule
}

// WithHttpBlock: 处理 nginx.conf 中 http 的配置块
func WithHttpBlock(conf *ngxHttpModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		// TODO: 处理 http 模块

		// TODO: 处理其它 http 模块

		// TODO: 处理 nginx.conf 中 http 的配置块
		preModuleType := ngxConf.ModuleType
		preCmdType := ngxConf.CmdType
		fmt.Printf("handle conf:http block,args:%v\n", ngxConf.Args)
		ngxConf.Args = []string{}
		ngxConf.ModuleType = ENgxHttpModule
		ngxConf.CmdType = TNgxHttpMainConf
		ngxConf.ParseFile("")
		ngxConf.ModuleType = preModuleType
		ngxConf.CmdType = preCmdType

		// TODO: 处理 http 模块
	}
}

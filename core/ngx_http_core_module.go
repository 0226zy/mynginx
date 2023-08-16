package core

import "fmt"

var httpCoreModule NgxHttpCoreModule

func init() {

	httpCoreModule = NgxHttpCoreModule{
		moduleType: ENgxHttpModule,
		moduleName: "ngx_http_core_module",
		commands: []*NgxModuleCommand{
			NewNgxModuleCommand(
				"server",
				TNgxHttpMainConf|TNgxConfBlock|TNgxConfNoArgs,
				WithServerBlock(&httpModuleConf)),
			NewNgxModuleCommand(
				"location",
				TNgxHttpSRVConf|TNgxHttpLocConf|TNgxConfBlock|TNgxConfTake12,
				WithLocationBlock(&httpModuleConf)),
			NewNgxModuleCommand(
				"server_name",
				TNgxHttpSRVConf|TNgxConf1More,
				WithServerName(&httpModuleConf)),
		},
	}
}

// NgxHttpCoreModule: http 模块
type NgxHttpCoreModule struct {
	moduleName string
	moduleType int
	commands   []*NgxModuleCommand
}

func (module *NgxHttpCoreModule) InitMaster()  {}
func (module *NgxHttpCoreModule) InitModule()  {}
func (module *NgxHttpCoreModule) InitProcess() {}

func (module *NgxHttpCoreModule) Index() int {
	return 0
}

func (module *NgxHttpCoreModule) GetCommands() []*NgxModuleCommand { return module.commands }
func (module *NgxHttpCoreModule) Name() string {
	return module.moduleName
}

func (module *NgxHttpCoreModule) Type() int {
	return module.moduleType
}

func GetNgxHttpCoreModule() NgxModule {
	return &httpCoreModule
}

// WithHttpBlock: 处理 nginx.conf 中 http 的配置块
func WithServerBlock(conf *ngxHttpModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("handle conf:http.server block,args:%v\n", ngxConf.Args)
		// TODO: 处理其它 http 模块

		// TODO: 处理 nginx.conf 中 http 的配置块
		preCmdType := ngxConf.CmdType
		ngxConf.Args = []string{}
		ngxConf.CmdType = TNgxHttpSRVConf
		ngxConf.ParseFile("")
		ngxConf.CmdType = preCmdType

		// TODO: 处理 http 模块
	}
}

// WithServerName: sever_name
func WithServerName(conf *ngxHttpModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("handle conf:http.server.server_name args:%v\n", ngxConf.Args)
	}
}

// WithListen: listen
func WithListen(conf *ngxHttpModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("handle conf:http.server.listen args:%v\n", ngxConf.Args)
	}
}

// WithLocationBlock : location
func WithLocationBlock(conf *ngxHttpModuleConf) CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {
		fmt.Printf("handle conf:http.server.location block,args:%v\n", ngxConf.Args)
		// location 后面的 value
		if len(ngxConf.Args) == 3 {
			fmt.Printf("location value handle,args len=3\n")
		} else {
			fmt.Printf("location value handle,args len=2\n")

		}
		// TODO: 处理其它 http 模块

		//处理 location 后面的配置块
		preCmdType := ngxConf.CmdType
		ngxConf.Args = []string{}
		ngxConf.CmdType = TNgxHttpLocConf
		ngxConf.ParseFile("")
		ngxConf.CmdType = preCmdType

		// TODO: 处理 http 模块
	}
}

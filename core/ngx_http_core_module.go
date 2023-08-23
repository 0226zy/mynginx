package core

import (
	"fmt"
	"strconv"
)

var httpCoreModule NgxHttpCoreModule

func init() {

	httpCoreModule = NgxHttpCoreModule{
		moduleType: ENgxHttpModule,
		moduleName: "ngx_http_core_module",
		ctx:        &ngxHttpCoreModuleConfCtx{},
		commands: []*NgxModuleCommand{

			NewNgxModuleCommand(
				"listen",
				TNgxHttpMainConf|TNgxConfBlock|TNgxHttpSRVConf,
				WithListen()),

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
	ctx        INgxHttpConfCtx
	commands   []*NgxModuleCommand
}

type ngxHttpCoreModuleConfCtx struct {
}

func (ctx *ngxHttpCoreModuleConfCtx) PreConfiguration(cf *NgxConf) int {
	return 0
}

func (ctx *ngxHttpCoreModuleConfCtx) PostConfiguration(cf *NgxConf) int {
	return 0
}

func (ctx *ngxHttpCoreModuleConfCtx) CreateMainConf(cf *NgxConf) interface{} {
	return &ngxHttpCoreMainConf{serverConfs: []*ngxHttpCoreSrvConf{}, ports: []*ngxHttpConfPort{}}
}

func (ctx *ngxHttpCoreModuleConfCtx) InitMainConf(cf *NgxConf) interface{} {
	return nil
}

func (ctx *ngxHttpCoreModuleConfCtx) CreateSrvConf(cf *NgxConf) interface{} {
	return &ngxHttpCoreSrvConf{}

}
func (ctx *ngxHttpCoreModuleConfCtx) MergeSrvConf(cf *NgxConf) interface{} {
	return nil
}

func (ctx *ngxHttpCoreModuleConfCtx) CreateLocConf(cf *NgxConf) interface{} {
	return &ngxHttpCoreLocConf{}
}

func (ctx *ngxHttpCoreModuleConfCtx) MergeLocConf(cf *NgxConf) interface{} {
	return nil
}

func (module *NgxHttpCoreModule) GetCtx() interface{} {
	return module.ctx
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
func WithListen() CmdHandler {
	return func(ngxConf *NgxConf, ngxModule *NgxModuleCommand) {

		value := ngxConf.Args[1]
		// TODO: parse url,从 url 中获取 port
		port, err := strconv.Atoi(value)
		if err != nil {
			panic("Error:" + err.Error())
		}
		// TODO: parse opt,url 支持opt
		// TODO: 支持指定 port 范围
		// Add listen
		httpMainConf := httpConfCtx.mainConfs[0]
		srvConf := httpConfCtx.srvConfs[0]

		fmt.Printf("handle conf:http.server.listen args:%v\n", ngxConf.Args)
		portConf := &ngxHttpConfPort{addrs: []*ngxHttpConfAddr{}}
		portConf.family = 1
		portConf.portType = 1
		portConf.port = port
		addr := &ngxHttpConfAddr{servers: []*ngxHttpCoreSrvConf{}}
		addr.servers = append(addr.servers, srvConf)
		portConf.addrs = append(portConf.addrs, addr)
		httpMainConf.ports = append(httpMainConf.ports, portConf)
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

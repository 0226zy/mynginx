package core

import "fmt"

var httpModule NgxHttpModule
var httpModuleConf ngxHttpModuleConf
var httpConfCtx *NgxHttpConfCtx

func init() {

	httpConfCtx = &NgxHttpConfCtx{}
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

type INgxHttpConfCtx interface {
	PreConfiguration(cf *NgxConf) int
	PostConfiguration(cf *NgxConf) int
	CreateMainConf(cf *NgxConf) interface{}
	InitMainConf(cf *NgxConf) interface{}
	CreateSrvConf(cf *NgxConf) interface{}
	MergeSrvConf(cf *NgxConf) interface{}
	CreateLocConf(cf *NgxConf) interface{}
	MergeLocConf(cf *NgxConf) interface{}
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

type ngxHttpCoreMainConf struct {
	serverConfs []*ngxHttpCoreSrvConf
	ports       []*ngxHttpConfPort
}

type ngxHttpCoreSrvConf struct {
	confCtx         *NgxHttpConfCtx
	serverName      string
	locCons         []*ngxHttpCoreLocConf
	connPoolSize    uint64
	requestPoolSize uint64
	listen          bool // 是否设置了 listen
}

type ngxHttpListenOpt struct {
	// listen 支持一些可选参数
}

type ngxHttpConfPort struct {
	family   int
	portType int
	port     int
	addrs    []*ngxHttpConfAddr
}

type ngxHttpConfAddr struct {
	opt     ngxHttpListenOpt
	servers []*ngxHttpCoreSrvConf
}

type ngxHttpCoreLocConf struct {
	locName string
}

type ngxHttpModuleConf struct {
}

func (module *NgxHttpModule) GetCtx() interface{} {
	return nil
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
		// 构造 NgxHttpConfCtx
		cycle := GetGlobalCycle()
		countHttpModule := cycle.CountHttpModule()
		httpConfCtx.mainConfs = make([]*ngxHttpCoreMainConf, countHttpModule)
		httpConfCtx.srvConfs = make([]*ngxHttpCoreSrvConf, countHttpModule)
		httpConfCtx.locConfs = make([]*ngxHttpCoreLocConf, countHttpModule)
		ctx_idx := 0
		for _, module := range cycle.Modules {
			if module.Type() != ENgxHttpModule {
				continue
			}

			ctx := module.GetCtx()
			if nil == ctx {
				continue
			}
			if ngxCtx, ok := ctx.(INgxHttpConfCtx); ok {

				if mainConf := ngxCtx.CreateMainConf(nil); mainConf != nil {
					httpConfCtx.mainConfs[ctx_idx] = convertToStruct[*ngxHttpCoreMainConf](mainConf)
				}

				if srvConf := ngxCtx.CreateSrvConf(nil); srvConf != nil {
					httpConfCtx.srvConfs[ctx_idx] = convertToStruct[*ngxHttpCoreSrvConf](srvConf)
				}

				if locConf := ngxCtx.CreateLocConf(nil); locConf != nil {
					httpConfCtx.locConfs[ctx_idx] = convertToStruct[*ngxHttpCoreLocConf](locConf)
				}

			} else {
				panic("ctx is not a INgxHttpConfCtx")
			}

		}

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

		for _, mainConf := range httpConfCtx.mainConfs {
			if nil == mainConf {
				continue
			}
			for _, port := range mainConf.ports {
				cycle.AddListening(port.port)
			}

		}
	}
}

func convertToStruct[T any](val interface{}) T {
	ret, ok := val.(T)
	if !ok {
		panic("convert to T failed")
	}
	return ret
}

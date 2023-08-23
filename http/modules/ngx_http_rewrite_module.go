package modules

import (
	"fmt"

	"github.com/0226zy/mynginx/core"
)

var httpRewriteModule NgxHttpRewriteModule
var rewriteConf ngxHttpRewriteModuleConf

func init() {
	rewriteConf = ngxHttpRewriteModuleConf{}

	httpRewriteModule = NgxHttpRewriteModule{
		moduleType: core.ENgxHttpModule,
		moduleName: "ngx_http_rewrite_module",
		commands: []*core.NgxModuleCommand{
			core.NewNgxModuleCommand(
				"return",
				core.TNgxHttpSRVConf|core.TNgxHttpLocConf|core.TNgxHttpSifConf|core.TNgxHttpLocConf|core.TNgxHttpLifConf|core.TNgxConfTake12,
				WithReturn(&rewriteConf)),
			core.NewNgxModuleCommand(
				"rewrite",
				core.TNgxHttpSRVConf|core.TNgxHttpLocConf|core.TNgxHttpSifConf|core.TNgxHttpLocConf|core.TNgxHttpLifConf|core.TNgxConfTake23,
				WithRewrite(&rewriteConf)),
		},
	}
}

// NgxHttpRewriteModule: http 模块
type NgxHttpRewriteModule struct {
	moduleName string
	moduleType int
	commands   []*core.NgxModuleCommand
}

type ngxHttpRewriteModuleConf struct{}

func (module *NgxHttpRewriteModule) InitMaster() {}
func (module *NgxHttpRewriteModule) GetCtx() interface{} {
	return nil
}
func (module *NgxHttpRewriteModule) InitModule()  {}
func (module *NgxHttpRewriteModule) InitProcess() {}

func (module *NgxHttpRewriteModule) Index() int {
	return 0
}

func (module *NgxHttpRewriteModule) GetCommands() []*core.NgxModuleCommand { return module.commands }
func (module *NgxHttpRewriteModule) Name() string {
	return module.moduleName
}

func (module *NgxHttpRewriteModule) Type() int {
	return module.moduleType
}

func GetNgxHttpRewriteModule() core.NgxModule {
	return &httpRewriteModule
}

// WithServerName: sever_name
func WithReturn(conf *ngxHttpRewriteModuleConf) core.CmdHandler {
	return func(ngxConf *core.NgxConf, ngxModule *core.NgxModuleCommand) {
		fmt.Printf("handle conf:http.server.location return args:%v\n", ngxConf.Args)
	}
}

// WithListen: rewrite
func WithRewrite(conf *ngxHttpRewriteModuleConf) core.CmdHandler {
	return func(ngxConf *core.NgxConf, ngxModule *core.NgxModuleCommand) {
		fmt.Printf("handle conf:http.server.location rewrite:%v\n", ngxConf.Args)
	}
}

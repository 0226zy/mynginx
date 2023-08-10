package core

const (
	ENgxConfModule int = 0
	ENgxCoreModule     = 1
)

// NgxModule 模块定义
type NgxModule interface {
	InitMaster()
	InitModule()
	InitProcess()

	// get and set 接口
	GetCommands() []NgxModuleCommand
	Type() int
	Name() string
	Index() int
}

type NgxModuleConf interface {
	Name() string
}

type NgxModuleCommand interface {
	Type() int
	Name() string
	Set(*NgxConf, *NgxModuleCommand, *NgxModuleConf)
}

// NgxModuleCtx 模块 ctx
type NgxModuleCtx interface {
	CreateConf(cycle *NgxCycle)
	InitConf(cycle *NgxCycle, conf interface{})
}

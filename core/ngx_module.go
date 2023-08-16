package core

const (
	ENgxConfModule  int = 0
	ENgxCoreModule      = 1
	ENgxEventModule int = 2
	ENgxHttpModule  int = 3
)

// NgxModule 模块定义
type NgxModule interface {
	InitMaster()
	InitModule()
	InitProcess()

	// get and set 接口
	GetCommands() []*NgxModuleCommand
	Type() int
	Name() string
	Index() int
}

type NgxModuleConf interface {
	Name() string
}

type CmdHandler func(*NgxConf, *NgxModuleCommand)

type NgxModuleCommand struct {
	Name    string
	CmdType int64
	Set     CmdHandler
}

func NewNgxModuleCommand(name string, cmdType int64, cmdFunc CmdHandler) *NgxModuleCommand {
	return &NgxModuleCommand{
		Name:    name,
		CmdType: cmdType,
		Set:     cmdFunc,
	}
}

// NgxModuleCtx 模块 ctx
type NgxModuleCtx interface {
	CreateConf(cycle *NgxCycle)
	InitConf(cycle *NgxCycle, conf interface{})
}

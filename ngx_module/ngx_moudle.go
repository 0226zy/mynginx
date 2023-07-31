package ngx_module

// NgxModule 模块定义
type NgxModule interface {
	Name() string
	Index() int

	InitMaster()
	InitModule()
	InitProcess()
}

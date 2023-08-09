package core

// NgxCycle: 全局信息
type NgxCycle struct {
	Modules []NgxModule
	// TODO: 对象内存池

	// pool * Pool
}

var cycle *NgxCycle

func CreateGlobalCycle() *NgxCycle {

	cycle = &NgxCycle{}
	return cycle
}

func GetGlobalCycle() *NgxCycle {
	return cycle
}

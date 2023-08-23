package core

import (
	"fmt"
	"os"
)

const (
	NgxConfOk int = 0

	NgxConfError = -1

	NgxConfBlockStart = 1
	NgxConfBlockDone  = 2
	NgxConfFileDone   = 3
)

// 配置项的类型枚举
const (
	TNgxConfNoArgs   int64 = 0x00000001
	TNgxConfTake1          = 0x00000002
	TNgxConfTake2          = 0x00000004
	TNgxConfTake3          = 0x00000008
	TNgxConfTake4          = 0x00000010
	TNgxConfTake5          = 0x00000020
	TNgxConfTake6          = 0x00000040
	TNgxConfTake7          = 0x00000080
	TNgxConfTake12         = (TNgxConfTake1 | TNgxConfTake2)
	TNgxConfTake13         = (TNgxConfTake1 | TNgxConfTake3)
	TNgxConfTake23         = (TNgxConfTake2 | TNgxConfTake3)
	TNgxConfTake123        = (TNgxConfTake1 | TNgxConfTake2 | TNgxConfTake3)
	TNgxConfTake1234       = (TNgxConfTake1 | TNgxConfTake2 | TNgxConfTake3 | TNgxConfTake4)

	TNgxConfArgsNumber = 0x000000ff
	TNgxConfBlock      = 0x00000100
	TNgxConfFlag       = 0x00000200
	TNgxConfAny        = 0x00000400
	TNgxConf1More      = 0x00000800
	TNgxConf2More      = 0x00001000
	TNgxDirectConf     = 0x00010000

	TNgxMainConf = 0x01000000
	TNgxAnyConf  = 0xFF000000

	// event 配置
	// event 下面的配置项
	TNgxEventConf = 0x02000000

	// http 配置
	// http 下面的配置
	TNgxHttpMainConf = 0x02000000
	// http.server 下面的配置
	TNgxHttpSRVConf = 0x04000000
	// http.server.location 下面的配置
	TNgxHttpLocConf = 0x08000000
	TNgxHttpUpsConf = 0x10000000
	TNgxHttpSifConf = 0x20000000
	TNgxHttpLifConf = 0x40000000
	TNgxHttpLMTConf = 0x80000000
)

// NgxConf  nginx 配置
type NgxConf struct {
	fileName       string
	confFile       *os.File
	confFileOffset int64
	confFileInfo   os.FileInfo
	confBuf        *NgxConfBuf
	ModuleType     int
	CmdType        int64
	Args           []string
	WorkerProcess  int
}

const (
	eParseFile = iota
	eParseBlock
	eParseParam
)

// NgxConfBuf 配置文件内容 buf
// TODO：使用统一的内存池 buf
type NgxConfBuf struct {
	Start int64
	End   int64

	FilePos int64
	FileEnd int64

	Pos  int64
	Last int64

	Data      []byte
	ParseType int64
}

// NgxConfBufferSize 配置文件读取的 buffer 最大字符数
const NgxConfBufferSize int64 = 4096

// CreateNgxConf: 构造 NgxConf
func CreateNgxConf() *NgxConf {
	return &NgxConf{
		fileName:       "",
		confFile:       nil,
		confBuf:        &NgxConfBuf{},
		confFileOffset: 0,
		WorkerProcess:  1,
	}
}

func (conf *NgxConf) ParseParameter() {
	conf.confBuf.ParseType = eParseParam
}

// ParseFile 解析配置
func (conf *NgxConf) ParseFile(fileName string) int {

	var parse_type int
	defer func() {
		if conf.confFile != nil && parse_type == eParseFile {
			fmt.Println("Close file")
			conf.confFile.Close()
		}
	}()

	var err error

	if len(fileName) != 0 {
		parse_type = eParseFile
		conf.confBuf.ParseType = eParseFile
		fmt.Printf("ngx conf:%s begin parse\n", fileName)
		conf.confFile, err = os.Open(fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return NgxError
		}

		conf.confFileInfo, err = conf.confFile.Stat()
		if err != nil {
			fmt.Println("Error getting file stats:", err)
			return NgxError
		}

		// 分配 buf
		conf.confBuf.Data = make([]byte, NgxConfBufferSize)
		conf.confBuf.Pos = 0
		conf.confBuf.Last = 0
		conf.confBuf.End = NgxConfBufferSize
		conf.confBuf.Start = 0

		// 文件 offset
		conf.confBuf.FilePos = 0
		conf.confBuf.FileEnd = 0

		// 从头开始解析配置文件

	} else if conf.confFile != nil {
		fmt.Printf("begin parse conf block\n")
		// 配置文件解析过程，进入模块的配置解析，读取模块的子配置
		parse_type = eParseBlock
		conf.confBuf.ParseType = eParseBlock

	} else {
		// 解析命令行参数
		parse_type = eParseParam
		conf.confBuf.ParseType = eParseParam
	}

	// 循环读取 tokens
	for {

		rc := conf.read_token()

		if rc == NgxError {
			// TODO: done clear
			return NgxConfError
		}

		// 子配置块全部读取完了
		if rc == NgxConfBlockDone {
			if parse_type != eParseBlock {
				// TODO log
				return NgxConfError
			}
			// TODO: done clear
			return NgxConfOk

		}

		// 配置文件全部读取完了
		if rc == NgxConfFileDone {
			if parse_type == eParseBlock {
				// err log
				return NgxConfError
			}
			// TODO: done clear
			return NgxConfOk
		}

		// 发现子配置块，其中 key 已经被读取了
		if rc == NgxConfBlockStart {

			// 命令行参数不支持{xx} 这种子配置项
			if parse_type == eParseParam {
				// errr log
				return NgxConfError
			}

		}

		// TODO: 自定义的 handler

		// 处理配置项
		rc = conf.conf_handler()
		if rc == NgxError {
			return NgxConfError
		}
	}

	return NgxConfOk
}

// read_token 读取一个 token（string）
/* buf 初始态:
* ==================================================
*  | nil | nil | nil | nil | nil | ...| nil |
*     |
*    pos
*    last
*    start
*                                        end
* ===============================================
* buf 读取一次文件内容后
* buf 遍历过程中
* buf 如果一个 token 以后
 */
func (conf *NgxConf) read_token() int {

	found := false
	last_space := true
	need_space := false

	// token 开始的位置
	buf := conf.confBuf
	token_start := buf.Pos
	var ch byte

	// 配置文件行数，从 1 开始计数
	line := 1

	// # 符号开启：注释开始
	sharp_comment := false
	// \ 转义符号，文件中是\,读取到 byte 中是'\\',ascii=92
	quoted := false
	// ' 单引号
	s_quoted := false
	// " 双引号
	d_quoted := false
	// $ 表达式：${xx}
	variable := false

	token_len := int64(0)
	for {

		// 1. 检查数组中是否还有字符未遍历
		if buf.Pos >= buf.Last {
			// 数组中没有字符可遍历：数组位空， 或者已经遍历到最后一位

			// 1.1 文件中是否还有字符未读取
			// TODO: 不能使用 Seek 获取当前文件的 offset
			// 如果是 parameter 解析，confFile 是 nil
			//fmt.Printf("offset:%d \n", conf.confFileOffset)
			if conf.confFileOffset >= int64(conf.confFileInfo.Size()) {

				// 读取到了token，或者读取了 半个 token
				// 此时应该还有字符需要读取，但是文件空了
				// 说明配置有问题
				if len(conf.Args) > 0 || !last_space {

					// 说明是在解析命令行参数时有问题
					if conf.confFile != nil {
						fmt.Println("error:unexpected end of parameter,expectind \";\"")
						return NgxError
					}
					fmt.Println("error:unexpected end of file, expecting \";\" or \"}\"")
					return NgxError
				}
				return NgxConfFileDone
			}

			// 当前在遍历的 token 的长度
			token_len = buf.Pos - token_start
			if token_len == NgxConfBufferSize {
				// 数据已经遍历完了，当前在遍历的 token 没有遇到结束符
				// 遇到一个长度大于 4096 的 token
				// 抛出错误，不支持这么长的 token

				if d_quoted {
					ch = '"'
				} else if s_quoted {
					ch = '\''
				} else {
					fmt.Println("error: too long parameter")
					return NgxError
				}
				fmt.Printf("error: too long parameter,probably missing terminating %c character\n", ch)

				return NgxError
			}

			if token_len > 0 {
				// 遍历到数据最后一个字符，但是数据当前有一个 token 没有遇到结束符
				// 该 token 还有字符在文件中
				// 遇到保留数据尾部这个未遍历完的 token
				// 将这个 token 移动到数据头部，读文件时从 b.start+ len 开始覆盖写数组
				// 回收已经被处理过字符占用的数组空间，重复利用
				/*
				*  |   c   |  c  | c | c |   | v | ;  | \n |  c1  | c1 |  c1  |
				*    b.start
				*																																b.end
				*																													b.last
				*																													b.pos
				*                                           token_start
				*                                           | --- token_len  ---|
				 */
				copy(buf.Data, buf.Data[token_start:token_start+token_len])
			}

			//文件剩余未读的字符大小
			readSize := conf.confFileInfo.Size() - conf.confFileOffset
			if readSize > buf.End-(buf.Start+token_len) {
				readSize = buf.End - (buf.Start + token_len)
			}

			buf_write_pos := buf.Start + token_len
			n, err := conf.confFile.ReadAt(buf.Data[buf_write_pos:buf_write_pos+readSize], conf.confFileOffset)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return NgxError
			}

			conf.confFileOffset += int64(n)
			buf.Pos = buf.Start + token_len
			buf.Last = buf.Pos + int64(n)
			token_start = buf.Start
		}
		// 2. 数组中还有字符未遍历

		// 2.1 获取当前要处理的字符
		ch = buf.Data[buf.Pos]
		//fmt.Printf("byte_code:%d ch:%s\n", ch, string(ch))

		// 2.2 指向下一个字符给下一次循环处理
		buf.Pos++

		//  2.3 开始处理当前字符: ch
		if ch == '\n' {
			// 换行
			line++
			// 换行，# 注释结束：# 只能注释一行，不能注释多行
			if sharp_comment {
				sharp_comment = false
			}

		}

		// 跳过配置文件中的注释
		if sharp_comment {
			continue
		}

		// 跳过 '\\'
		if quoted {
			quoted = false
			continue
		}

		// "xxx" / 'xxx' 已经遍历完了，xxx 已经被拷贝到 args 里面
		// 进一步拷贝后面否还有 token 要读取，还是返回给上层处理
		if need_space {
			//  "xxx" cc ：这种情况，后面还有 cc 要读取
			//  "xxx"\tcc ：这种情况，后面还有 cc 要读取
			//  "xxx"\r\ncc ：这种情况，后面还有 cc 要读取
			if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
				last_space = true
				need_space = false
				continue
			}
			if ch == ';' {
				return NgxOk
			}
			if ch == '{' {
				return NgxConfBlockStart
			}

			if ch == ')' {
				//"xx")cc：继续读取后面的)cc
				last_space = true
				need_space = false
			} else {
				//"xxx"xx：这种配置是有问题的
				// 配置有问题
				return NgxError
			}
		}

		// 上一个是空格
		// 准备开始遍历下一个 token
		if last_space {
			token_start = buf.Pos - 1

			// 遇到连续的空格
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
				continue
			}

			switch ch {
			case ';':
				// TODO 检查
				return NgxOk
			case '{':
				// TODO 检查
				return NgxConfBlockStart
			case '}':
				// TODO 检查
				return NgxConfBlockDone
			case '#':
				sharp_comment = true
				continue
			case '\\':
				quoted = true
				last_space = false
				continue
			case '"':
				// 遇到"xxx"，只拷贝 xxx: token_start + 1，跳过 "
				token_start++
				d_quoted = true
				last_space = false
				continue
			case '\'':
				// 遇到'xxx'，只拷贝 xxx: token_start + 1，跳过 '
				token_start++
				s_quoted = true
				last_space = false
				continue
			case '$':
				variable = true
				last_space = false
				continue
			default:
				// 遇到正常的字符
				last_space = false

			}
		} // if last_space ==1 结束

		if last_space == false {

			// 上一个不是空格
			// 当前正则遍历一个 token 中

			// 在遍历一个表达式: ${xxx}
			if ch == '{' && variable {
				// variable 唯一的作用在这个分支：
				//避免 $ 后面的 { 被误判为是子配置项开始符(NgxConfBlockStart)
				continue
			}

			// 重置表达式符号
			// 嵌套$ 的场景: ${xxx${xx}xxx}
			variable = false

			if ch == '\\' {
				quoted = true
				continue
			}

			if ch == '$' {
				variable = true
				continue
			}

			// "xx" 结束
			if d_quoted && ch == '"' {
				d_quoted = false
				found = true
				// "xxx" 结束，后面必须是一个空格类的结束符
				need_space = true
			}
			// 'x' 结束
			if s_quoted && ch == '\'' {
				s_quoted = false
				found = true
				// 'xxx' 结束，后面必须是一个空格类的结束符
				need_space = true
			}

			// 当前不在 "xx" 和 'xx' 中
			// 并且遇到了 token 结束符，则token 遍历结束，可以拷贝出来
			if (!d_quoted && !s_quoted) && conf.is_end_char(ch) {
				last_space = true
				found = true
			}

			// token 遇到结束符号，可以拷贝这个 token 了
			if found {
				arg := make([]byte, buf.Pos-token_start)
				idx := 0
				// 逐字符拷贝
				for copy_idx := token_start; copy_idx < buf.Pos-1; {
					// 如果是转义字符
					if buf.Data[copy_idx] == '\\' {
						// 这里默认出现 \ 后面一定还有字符
						switch buf.Data[copy_idx+1] {
						// 如果是 \" 只拷贝 "
						// 如果是 \' 只拷贝 '
						// 如果是 \\ 只拷贝 \
						case '"':
						case '\'':
						case '\\':
							copy_idx++
							break

						// \r \t \n	保留
						case 'r':
							arg[idx] = '\r'
							copy_idx = copy_idx + 2
							continue
						case 't':
							arg[idx] = '\t'
							copy_idx = copy_idx + 2
							continue
						case 'n':
							arg[idx] = '\n'
							copy_idx = copy_idx + 2
							continue
						} // switch 结束
					}
					//fmt.Printf("idx:%d copy_idx:%d Pos:%d,%d ch:%s\n", idx, copy_idx, buf.Pos, buf.Pos-1, string(buf.Data[copy_idx]))
					// 拷贝字符
					arg[idx] = buf.Data[copy_idx]
					idx++
					copy_idx++
				} // for 循环拷贝结束
				conf.Args = append(conf.Args, string(arg[0:idx]))

				// 一个完整的 key value; 配置读取完，返回到上层，交由对应模块来处理 key value
				if ch == ';' {
					return NgxOk
				}

				// key {xxx}; 这种含有子配置项的配置 key 读取完，返回到上层，交由对应模块来处理 key
				// {xx} 由对应模块的配置处理代码来调用 Parse 完成 {xx}; 的读取和处理
				if ch == '{' {
					return NgxConfBlockStart
				}
				found = false

			} // if found == 1 end
			// 当前字符处理结束，继续下一轮循环的处理
		} // if last_space == 0 end
	} //  for end

}

func (conf *NgxConf) conf_handler() int {

	ngxCycle := GetGlobalCycle()

	found := false
	commandName := conf.Args[0]

	for _, module := range ngxCycle.Modules {

		if module.Type() != ENgxConfModule && module.Type() != conf.ModuleType {
			continue
		}

		// 找到配置对应的模块，调用模块的配置处理回调函数
		commands := module.GetCommands()
		if len(commands) == 0 {
			continue
		}

		for _, command := range commands {

			if (command.CmdType & conf.CmdType) == 0 {
				continue
			}

			if command.Name != commandName {
				continue
			}

			// TODO
			// 配置项合法检查
			// TODO: get conf
			// var module_conf *NgxConf
			// TODO: 这里的实现比较复杂，这里要找到 moudle 对应的 module_conf
			//  module 里面的 Set 是将配置项存储到对应的 module_conf 里面
			// 如果使用 interface 这里涉及到 go package 的循环引用: 遇到把 NgxConf 传递到 module 定义的包内
			command.Set(conf, command)
			found = true
		}
	}

	if !found {
		fmt.Printf("error: cannot find conf item:%s command\n", commandName)
	}
	// 清空 Args
	conf.Args = []string{}
	return NgxOk
}

func (conf *NgxConf) is_end_char(ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '{' || ch == ';' {
		return true
	}
	return false
}

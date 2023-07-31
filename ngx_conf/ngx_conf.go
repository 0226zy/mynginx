package ngx_conf

import "fmt"

type NgxConf struct {
	fileName string
}

func CreateNgxConf() *NgxConf {
	return &NgxConf{}
}

func (NgxConf *conf) Parse() {
	fmt.Printf("ngx conf:%s begin parse\n", conf.fileName)
	// TODO
	file, err := os.open(os.args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:", err)
	}

	conf.parse_conf(file, fileInfo)
}

func (NgxConf *conf) parse_conf(file *File, fileInfo FileInfo) {
	conf_buffer := 4096
	// 2. 创建一个 4096 的buf，并标记一个buf 中index 的flag pos
	buf := make([]byte, conf_buffer)

	// buf 遍历标记位
	pos := 0
	last := 0

	// buf 启始标记位
	start := 0
	end := len(buf)

	fmt.Printf("pos:%d start:%d end:%d\n", pos, start, end)

	// 文件读取标记位
	offset := 0

	for {

		// 3.1. 检查 buf 是否为空
		if pos >= last {
			// 3.2. 为空则检查 文件是否已经读完成，读完退出
			if offset >= int(fileInfo.Size()) {
				break
			}

			data_len = pos - start
			if data_len == conf_buffer {
				// TODO:
				// NOTE: 什么情况会出现 data_len>0 && data_len<conf_buffer
			}

			if data_len {

			}

			// 3.3. 文件没有读完，则最大读取4096，如果文件剩余大小小于 4096，则只读取剩余大小
			readSize := fileInfo.Size() - offset
			if readSize > end-(start+data_len) {
				readSize = end - (start + data_len)
			}

			n, err := file.ReadAt(buf[start+buf_len:], int64(buf_len))
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			pos = start + data_len
			last = pos + n

			buf = buf[:n]
		}

		// 3.4. 读取 buf 不为空，则pos 前进一位，输出pos 指向buf 的字符，输出后 conintue
		fmt.Printf("%c", buf[pos%4096])
		pos++
	}
}

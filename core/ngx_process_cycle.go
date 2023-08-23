package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/0226zy/mynginx/event"
)

type ProcessCycle struct {
	masterProcess bool
	workProcess   bool
	singleProcess bool
	childPids     []int
}

func NewProcessCycle() *ProcessCycle {
	return &ProcessCycle{
		masterProcess: false,
		workProcess:   false,
		singleProcess: false,
	}
}

func (process *ProcessCycle) ParseArgs() {

	if len(os.Args) < 2 {
		process.masterProcess = true
	}

	if len(os.Args) >= 2 && os.Args[1] == "worker" {
		process.workProcess = true
	}

	if len(os.Args) >= 2 && os.Args[1] == "single" {
		process.singleProcess = true
	}

	if len(os.Args) >= 2 && os.Args[1] == "master" {
		process.masterProcess = true
	}
}

func (process *ProcessCycle) IsMasterProcess() bool {
	return process.masterProcess
}

func (process *ProcessCycle) ProcessCycle() {

	if process.masterProcess {
		process.MasterProcessCycle()
		return
	}

	if process.workProcess {
		process.WorkerProcessCycle()
		return
	}
	if process.singleProcess {
		process.SingleProcessCycle()
		return
	}
}

func (process *ProcessCycle) MasterProcessCycle() {

	sigCh := make(chan os.Signal, 1)

	// 注册信号
	signal.Notify(sigCh,
		os.Interrupt,
		//syscall.SIGCHLD,
		//	syscall.SIGALRAM,
		syscall.SIGHUP,   // 重载配置
		syscall.SIGUSR1,  // 重启
		syscall.SIGWINCH, // noaccept signal
		syscall.SIGTERM,  // term
		syscall.SIGQUIT,  // shutdown
		syscall.SIGUSR2,  // change bin 软重启
		syscall.SIGINT,
	//	syscall.SIGTERM
	)

	// 启动work 进程
	ngxCycle := GetGlobalCycle()
	workProcessN := ngxCycle.Conf.WorkerProcess

	for num := workProcessN; num > 0; num-- {
		process.spawnProcess()
	}

	listener := ngxCycle.Listening[0].listener
	listenerFile, err := listener.(*net.TCPListener).File()
	if err != nil {
		fmt.Printf("pid:%d err:%v\n", os.Getpid(), err)
	}

	epoll, _ := event.NewEpoll()
	defer epoll.Close()

	epoll.Add(int(listenerFile.Fd()))

	// 启动 cache 管理进程
	for {
		fmt.Println("master process Waiting for  a signal...")
		select {
		case sig := <-sigCh:
			{
				fmt.Printf("master process recevied siganl:%s\n", sig)
				for _, pid := range process.childPids {
					fmt.Println("master pid:%d kill pid:%d\n", os.Getpid(), pid)
					syscall.Kill(pid, syscall.SIGTERM)
				}
				fmt.Println("master nginx quit")
				os.Exit(0)
			}
		case fds := <-epoll.WaitWithChan():
			{
				fmt.Printf("master pid:%d await fds len:%d\n", os.Getpid(), len(fds))
				for _, fd := range fds {
					fmt.Printf("master pid:%d recevied fd:%d\n", os.Getpid(), fd)

					if fd == int(listenerFile.Fd()) {
						conn, err := listener.Accept()
						if err != nil {
							fmt.Println("Error:", err)
							continue
						}
						fmt.Printf("master pid:%d recevied conn:%+v\n", os.Getpid(), conn)
						go process.handleConnection(conn, "master")
					}

				}
			}
		}
	}
}

func (process *ProcessCycle) SingleProcessCycle() {

	sigCh := make(chan os.Signal, 1)

	// 注册信号
	signal.Notify(sigCh,
		syscall.SIGCHLD,
		//		syscall.SIGALRAM,
		syscall.SIGHUP,   // 重载配置
		syscall.SIGUSR1,  // 重启
		syscall.SIGWINCH, // noaccept signal
		syscall.SIGTERM,  // term
		syscall.SIGQUIT,  // shutdown
		syscall.SIGUSR2,  // change bin 软重启
		syscall.SIGINT,
	//	syscall.SIGTERM
	)

	// 启动 cache 管理进程
	for {
		fmt.Println("signle process Waiting for  a signal...")
		sig := <-sigCh
		fmt.Printf("single process recevied siganl:%s\n", sig)
		// loop
	}
}

func (process *ProcessCycle) WorkerProcessCycle() {
	fmt.Println("worker process begin ...")
	// 执行所有模块的的 init_process

	fd, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error,Atoi:", err)
		return
	}

	listenerFile := os.NewFile(uintptr(fd), "listener")
	listener, err := net.FileListener(listenerFile)
	if err != nil {
		fmt.Println("Error net.FileListener:", err)
		return
	}

	epoll, _ := event.NewEpoll()
	defer epoll.Close()

	epoll.Add(int(listenerFile.Fd()))
	for {
		fds, err := epoll.Wait()
		if err != nil {
			fmt.Printf("pid:%d wait err:\n", os.Getpid, err)
		}
		fmt.Printf("worker pid:%d epoll wait await fds len:%d\n", os.Getpid(), len(fds))
		for _, fd := range fds {
			fmt.Printf("work pid:%d recevied fd:%d \n", os.Getpid(), fd)
			if fd == int(listenerFile.Fd()) {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				fmt.Printf("worker pid:%d recevied conn:%+v\n", os.Getpid(), conn)
				go process.handleConnection(conn, "work")
			}
		}
	}
}

func (process *ProcessCycle) spawnProcess() error {

	pid := os.Getpid()

	files := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	cycle := GetGlobalCycle()

	args := []string{os.Args[0], "worker"}
	args = append(args, os.Args[1:]...)
	for _, listen := range cycle.Listening {
		file, _ := listen.listener.(*net.TCPListener).File()
		fd, err := syscall.Dup(int(file.Fd()))
		if err != nil {
			fmt.Println("Error:", err)
		}
		files = append(files, uintptr(fd))
		args = append(args, strconv.Itoa(fd))
	}
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: files,
	}

	fmt.Printf("fork  for path:%s arg1:%s\n", args[0], args[1])
	childPid, err := syscall.ForkExec(args[0], args, execSpec)
	if err != nil {
		fmt.Printf("process: %d, failed to forkexec with err:%s\n", pid, err.Error())
		return err
	}
	process.childPids = append(process.childPids, childPid)
	fmt.Printf("childPid:%d\n", childPid)
	return nil
}

func (process *ProcessCycle) handleConnection(conn net.Conn, processName string) {
	pid := os.Getpid()
	defer conn.Close()
	fmt.Println("handleonnection:", processName)

	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)

	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Printf("%s pid:%d http read failed:%v\n", processName, pid, err)
	}
	res := process.processRequest(req)
	err = res.Write(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	//for {
	//	input, err := reader.ReadString('\n')
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Printf("%s pid:%d err:%v\n", processName, pid, err)
	//		}
	//		fmt.Printf("%s pid:%d conn:%v close\n", processName, pid, conn)
	//		return
	//	}
	//	fmt.Printf("%s pid:%d input:%s\n", processName, pid, input)
	//	resp := process.html()
	//	_, err = writer.Write([]byte(resp))
	//	if err != nil {
	//		fmt.Printf("%s pid:%d writer err:%v\n", processName, pid, err)
	//		return
	//	}
	//}

}
func (process *ProcessCycle) processRequest(req *http.Request) *http.Response {
	fmt.Printf("pid :%d http req:%v\n", os.Getpid(), req)
	body := "Hello, World!\r\n from mynginx v0.0.0"
	headers := http.Header{}
	headers.Add("Content-Type", "text/plain")
	headers.Add("Content-Length", fmt.Sprintf("%d", len(body)))

	return &http.Response{
		StatusCode:    http.StatusOK,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        headers,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

func (process *ProcessCycle) html() string {
	html := `<!DOCTYPE html><html><head><meta charset="utf-8"><title> Simple HTML Page</title></head><body><h1> Hello, Mynginx!</h1><p>This is a simple HTML page.</p></body></html>`
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=utf-8\r\nContent-Length: %d\r\n\r\n%s", len(html), html)
}

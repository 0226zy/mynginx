package event

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
)

type Epoll interface {
	Add(fd int) error
	Remove(fd int) error
	Wait() ([]int, error)
	WaitWithChan() <-chan []int
	Close() error
}

type DarwinEpoll struct {
	fd      int
	ts      syscall.Timespec
	mu      *sync.RWMutex
	changes []syscall.Kevent_t
}

func NewEpoll() (Epoll, error) {
	p, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}
	//_, err = syscall.Kevent(p, []syscall.Kevent_t{
	//	{
	//		Ident:  0,
	//		Filter: syscall.EVFILT_USER,
	//		Flags:  syscall.EV_ADD | syscall.EV_CLEAR,
	//	}, nil, nil,
	//	})
	//	if err != nil {
	//		panic(err)
	//	}

	return &DarwinEpoll{
		fd: p,
		mu: &sync.RWMutex{},
		ts: syscall.NsecToTimespec(1e9),
	}, nil
}

func (e *DarwinEpoll) Remove(fd int) error {

	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.changes) <= 1 {
		e.changes = nil
	} else {
		changes := make([]syscall.Kevent_t, 0, len(e.changes)-1)
		ident := uint64(fd)
		for _, ke := range e.changes {
			if ke.Ident != ident {
				changes = append(changes, ke)
			}
			e.changes = changes
		}
	}
	return nil
}
func (e *DarwinEpoll) Add(fd int) error {
	if e := syscall.SetNonblock(fd, true); e != nil {
		return errors.New("udev:unixSetNonblock failed")
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.changes = append(e.changes, syscall.Kevent_t{
		Ident:  uint64(fd),
		Flags:  syscall.EV_ADD | syscall.EV_EOF,
		Filter: syscall.EVFILT_READ,
	})
	return nil
}

func (e *DarwinEpoll) Wait() ([]int, error) {
	events := make([]syscall.Kevent_t, 10)
	e.mu.RLock()
	changes := e.changes
	e.mu.RUnlock()
	ret := []int{}

retry:
	n, err := syscall.Kevent(e.fd, changes, events, nil)
	fmt.Printf("pid:%d await\n", os.Getpid())
	if err != nil {
		if err == syscall.EINTR {
			goto retry
		}
		return nil, err
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	for i := 0; i < n; i++ {
		ret = append(ret, int(events[i].Ident))
	}
	return ret, nil
}

func (e *DarwinEpoll) WaitWithChan() <-chan []int {
	ch := make(chan []int, 10)
	go func() {
		for {
			fds, err := e.Wait()
			if err != nil {
				continue
			}
			if len(fds) == 0 {
				continue
			}
			ch <- fds
		}
	}()
	return ch
}

func (e *DarwinEpoll) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.changes = nil
	return syscall.Close(e.fd)
}

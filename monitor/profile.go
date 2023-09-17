package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	durationSecond  uint32 = 60
	isSamplingTrace        = false

	serverName = getServerName()
	pid        = syscall.Getpid()
	timeFormat = "20060102T150405"

	status      uint32
	statusStart uint32 = 1 // status=1
	statusStop  uint32     // status=0
)

func Monitor() {
	go func() {
		p := NewProfile()
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1)
		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				p.StartOrStop()
			default:
				slog.Error("Got unregistered signal:", v)
			}
		}
	}()
}

type Profile struct {
	files    []string
	closeFns []func()

	//ctx    context.Context
	stopCh chan struct{}
}

// NewProfile create a new profile
func NewProfile() *Profile {
	p := new(Profile)
	p.stopCh = make(chan struct{})
	return p
}

// StartOrStop start and stop sampling profile, the first call to start sampling data, the default maximum is 60 seconds,
// in less than 60s, if the second execution will actively stop sampling profile
func (p *Profile) StartOrStop() {
	if isStart() {
		p.startProfile()
	} else if isStop() {
		p.stopProfile()
	}
}

func (p *Profile) startProfile() {
	fmt.Printf("[profile] start sampling profile, status=%d\n", status)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	err := p.cpu()
	if err != nil {
		fmt.Println(err)
	}

	err = p.mem()
	if err != nil {
		fmt.Println(err)
	}

	err = p.goroutine()
	if err != nil {
		fmt.Println(err)
	}

	err = p.block()
	if err != nil {
		fmt.Println(err)
	}

	err = p.mutex()
	if err != nil {
		fmt.Println(err)
	}

	err = p.threadCreate()
	if err != nil {
		fmt.Println(err)
	}

	if isSamplingTrace {
		err = p.tracing()
		if err != nil {
			fmt.Println(err)
		}
	}

	go p.checkTimeout()
}

func (p *Profile) stopProfile() {
	fmt.Printf("[profile] stop sampling profile, status=%d\n", status)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if p == nil || len(p.closeFns) == 0 {
		return
	}

	for _, fn := range p.closeFns {
		fn()
	}

	select {
	case p.stopCh <- struct{}{}:
	default:
	}

	// reset profile
	p = NewProfile() //nolint
}

func (p *Profile) checkTimeout() {
	if p == nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(durationSecond)) //nolint
	select {
	case <-p.stopCh:
		fmt.Println("[profile] reason for stopping: manual")
		return
	case <-ctx.Done():
		if isStop() {
			p.stopProfile()
		}
		fmt.Println("[profile] reason for stopping: timeout")
	}
}

func (p *Profile) cpu() error {
	profileName := "cpu"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	_ = pprof.StartCPUProfile(f)

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	})

	return nil
}

func (p *Profile) mem() error {
	profileName := "mem"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	old := runtime.MemProfileRate
	runtime.MemProfileRate = 4096

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		_ = pprof.Lookup("heap").WriteTo(f, 0)
		_ = f.Close()
		runtime.MemProfileRate = old
	})
	return nil
}

func (p *Profile) goroutine() error {
	profileName := "goroutine"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		_ = pprof.Lookup(profileName).WriteTo(f, 2)
		_ = f.Close()
	})
	return nil
}

func (p *Profile) block() error {
	profileName := "block"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	runtime.SetBlockProfileRate(1)

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		_ = pprof.Lookup(profileName).WriteTo(f, 0)
		_ = f.Close()
		runtime.SetBlockProfileRate(0)
	})

	return nil
}

func (p *Profile) mutex() error {
	profileName := "mutex"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	runtime.SetMutexProfileFraction(1)

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		if mp := pprof.Lookup(profileName); mp != nil {
			_ = mp.WriteTo(f, 0)
		}
		_ = f.Close()
		runtime.SetMutexProfileFraction(0)
	})

	return nil
}

func (p *Profile) threadCreate() error {
	profileName := "threadcreate"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		if mp := pprof.Lookup(profileName); mp != nil {
			_ = mp.WriteTo(f, 0)
		}
		_ = f.Close()
	})

	return nil
}

func (p *Profile) tracing() error {
	profileName := "trace"
	file := getFilePath(profileName)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	err = trace.Start(f)
	if err != nil {
		_ = f.Close()
		return err
	}

	p.files = append(p.files, file)
	p.closeFns = append(p.closeFns, func() {
		trace.Stop()
		_ = f.Close()
	})

	return nil
}

// SetDurationSecond set sampling profile duration
func SetDurationSecond(d uint32) {
	atomic.StoreUint32(&durationSecond, d)
}

// EnableTrace enable sampling trace profile
func EnableTrace() {
	isSamplingTrace = true
}

func isStart() bool {
	return atomic.CompareAndSwapUint32(&status, statusStop, statusStart)
}

func isStop() bool {
	return atomic.CompareAndSwapUint32(&status, statusStart, statusStop)
}

func getFilePath(profileName string) string {
	dir := filepath.Join(os.TempDir(), serverName+"_profile")
	_ = os.MkdirAll(dir, 0766)

	return filepath.Join(dir, fmt.Sprintf("%s_%d_%s_%s.out",
		time.Now().Format(timeFormat), pid, serverName, profileName))
}

func getServerName() string {
	_, name := filepath.Split(os.Args[0])
	return strings.TrimSuffix(name, path.Ext(name))
}

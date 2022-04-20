package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/MetaverseTopDJ/Scaffold/logger"

	"github.com/facebookgo/grace/gracenet"
	"google.golang.org/grpc"
)

var (
	// Used to indicate a graceful restart in the new process.
	envKey = "LISTEN_FDS"
	ppid   = os.Getppid()
)

// GraceGrpc is used to wrap a grpc server that can be gracefully terminated & restarted
type GraceGrpc struct {
	server   *grpc.Server
	grace    *gracenet.Net
	listener net.Listener
	errors   chan error
	pidPath  string
	logger   *logger.Logger
}

// NewGraceGrpc 优雅的 Grpc 实例
func NewGraceGrpc(s *grpc.Server, net, addr, pidPath string) (*GraceGrpc, error) {
	l := logger.NewLogger()
	gr := &GraceGrpc{
		server: s,
		grace:  &gracenet.Net{},
		//for  StartProcess error.
		errors:  make(chan error),
		pidPath: pidPath,
		logger:  l,
	}
	listener, err := gr.grace.Listen(net, addr)
	if err != nil {
		return nil, err
	}
	gr.listener = listener
	return gr, nil
}

// storePid is used to write out PID to pidPath
// 存储 Pid 进程ID
func (gr *GraceGrpc) storePid(pid int) error {
	pidPath := gr.pidPath
	if pidPath == "" {
		return fmt.Errorf("No pid file path: %s", pidPath)
	}

	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("Could not open pid file: %v", err)
	}
	defer pidFile.Close()

	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("Could not write to pid file: %s", err)
	}
	return nil
}

// handleSignal 处理信号
func (gr *GraceGrpc) handleSignal() <-chan struct{} {
	terminate := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 10)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
		for {
			sig := <-ch
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(ch)
				gr.server.GracefulStop()
				close(terminate)
				return
			case syscall.SIGUSR2:
				if _, err := gr.grace.StartProcess(); err != nil {
					gr.errors <- err
				}
			}
		}
	}()
	return terminate
}

// startServe 开启服务
func (gr *GraceGrpc) startServe() {
	if err := gr.server.Serve(gr.listener); err != nil {
		gr.errors <- err
	}
}

// Serve is used to start grpc server.
// Serve will gracefully terminated or restarted when handling signals.
func (gr *GraceGrpc) Serve() error {
	if gr.listener == nil || gr.logger == nil {
		return fmt.Errorf("gracegrpc must construct by new")
	}

	inherit := os.Getenv(envKey) != ""
	pid := os.Getpid()
	addrString := gr.listener.Addr().String()

	if inherit {
		if ppid == 1 {
			gr.logger.Info("Listening on init activated %s\n", addrString)
		} else {
			gr.logger.Info("Graceful handoff of %s with new pid %d replace old pid %d\n", addrString, pid, ppid)
		}
	} else {
		gr.logger.Info("Serving %s with pid %d\n", addrString, pid)
	}

	if err := gr.storePid(pid); err != nil {
		return err
	}

	go gr.startServe()

	if inherit && ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to close parent: %s", err)
		}
	}

	terminate := gr.handleSignal()

	select {
	case err := <-gr.errors:
		return err
	case <-terminate:
		gr.logger.Info("Exiting pid %d.", os.Getpid())
		return nil
	}
}

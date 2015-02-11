package resetter

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/pivotal-cf/cf-redis-broker/credentials"
	"github.com/pivotal-cf/cf-redis-broker/redisconf"

	"code.google.com/p/go-uuid/uuid"
)

type portChecker interface {
	Check(address *net.TCPAddr, timeout time.Duration) error
}

type Shell interface {
	Run(command *exec.Cmd) ([]byte, error)
}

type Resetter struct {
	defaultConfPath string
	confPath        string
	portChecker     portChecker
	shell           Shell
	monitExecutable string
	timeout         time.Duration
}

func New(defaultConfPath string, confPath string, portChecker portChecker, shell Shell, monitExecutable string) *Resetter {
	return &Resetter{
		defaultConfPath: defaultConfPath,
		confPath:        confPath,
		portChecker:     portChecker,
		shell:           shell,
		monitExecutable: monitExecutable,
		timeout:         time.Second * 10,
	}
}

func (resetter *Resetter) DeleteAllData() error {
	if err := resetter.stopRedis(); err != nil {
		return err
	}

	if err := resetter.deleteData(); err != nil {
		return err
	}

	if err := resetter.resetConf(); err != nil {
		return err
	}

	if err := resetter.startRedis(); err != nil {
		return err
	}

	credentials, err := credentials.Parse(resetter.confPath)
	if err != nil {
		return err
	}

	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", credentials.Port))
	if err != nil {
		return err
	}

	return resetter.portChecker.Check(address, resetter.timeout)
}

func (resetter *Resetter) stopRedis() error {
	resetter.shell.Run(exec.Command(resetter.monitExecutable, "stop", "redis"))
	redisProcessDead := make(chan bool)
	go func(c chan<- bool) {
		for {
			cmd := exec.Command("pgrep", "redis-server")
			output, _ := cmd.CombinedOutput()
			if len(output) == 0 {
				c <- true
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}(redisProcessDead)

	timer := time.NewTimer(resetter.timeout)
	defer timer.Stop()
	select {
	case <-redisProcessDead:
		break
	case <-timer.C:
		return errors.New("timed out waiting for redis process to die after 10 seconds")
	}

	return nil
}

func (resetter *Resetter) startRedis() error {
	redisStarted := make(chan bool)
	go func(c chan<- bool) {
		for {
			_, err := resetter.shell.Run(exec.Command(resetter.monitExecutable, "start", "redis"))
			if err == nil {
				c <- true
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}(redisStarted)

	timer := time.NewTimer(resetter.timeout)
	defer timer.Stop()
	select {
	case <-redisStarted:
		break
	case <-timer.C:
		return errors.New("timed out waiting for redis process to be started by monit after 10 seconds")
	}

	return nil
}

func (_ *Resetter) deleteData() error {
	if err := os.Remove("appendonly.aof"); err != nil {
		return err
	}

	os.Remove("dump.rdb")
	return nil
}

func (resetter *Resetter) resetConf() error {
	conf, err := redisconf.Load(resetter.defaultConfPath)
	if err != nil {
		return err
	}

	conf.Set("requirepass", uuid.NewRandom().String())

	if err := conf.Save(resetter.confPath); err != nil {
		return err
	}

	return nil
}

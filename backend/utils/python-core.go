package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID      string
	Content string
	Resp    chan string
}

type PythonProcess struct {
	cmd       *exec.Cmd
	writer    *bufio.Writer
	sendQueue chan *Request
	pending   sync.Map // map[id] = chan
}

func NewPythonProcess(pythonPath string, pythonFile string) (*PythonProcess, error) {
	cmd := exec.Command(pythonPath, pythonFile)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p := &PythonProcess{
		cmd:       cmd,
		writer:    bufio.NewWriter(stdin),
		sendQueue: make(chan *Request, 100),
	}

	// 启动 Python 进程
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 检查进程是否正常启动（超时判断）
	startupTimeout := time.After(5 * time.Second)
	startupCheck := make(chan error, 1)

	go func() {
		// 等待进程启动完成
		err := cmd.Wait()
		startupCheck <- err
	}()

	select {
	case err := <-startupCheck:
		if err != nil {
			return nil, errors.New("Python进程启动失败: " + err.Error())
		}
	case <-startupTimeout:
		// 进程仍在运行，说明启动成功
		// 不需要杀死进程，因为它在正常运行
	}

	// 后台处理 Python 输出
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// // 输出原始Python进程输出到日志
			// Log("RDKit Python输出: " + line)

			var resp map[string]string
			json.Unmarshal([]byte(line), &resp)

			// 找到对应的 request channel
			if ch, ok := p.pending.Load(resp["id"]); ok {
				ch.(chan string) <- resp["reply"] // 回传结果
				p.pending.Delete(resp["id"])
			}
		}
	}()

	// 后台发送请求
	go func() {
		for req := range p.sendQueue {
			// 保存等待通道
			p.pending.Store(req.ID, req.Resp)

			data := map[string]string{
				"id":  req.ID,
				"msg": req.Content,
			}
			jsonBytes, _ := json.Marshal(data)
			p.writer.Write(jsonBytes)
			p.writer.WriteString("\n")
			p.writer.Flush()
		}
	}()

	// 设置信号处理器，确保关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if p.IsRunning() {
			p.Close()
		}
		os.Exit(0)
	}()

	return p, nil
}

func (p *PythonProcess) SendAndWait(msg string) (string, error) {
	return p.SendAndWaitWithTimeout(msg, 30*time.Second)
}

func (p *PythonProcess) SendAndWaitWithTimeout(msg string, timeout time.Duration) (string, error) {
	req := &Request{
		ID:      uuid.New().String(),
		Content: msg,
		Resp:    make(chan string, 1),
	}

	p.sendQueue <- req

	select {
	case res := <-req.Resp:
		return res, nil
	case <-time.After(timeout):
		// 超时后清理pending状态
		p.pending.Delete(req.ID)
		return "", errors.New("Python进程响应超时")
	}
}

// Close 关闭Python进程
func (p *PythonProcess) Close() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	// 关闭发送队列
	close(p.sendQueue)

	// 清理所有pending的请求
	p.pending.Range(func(key, value interface{}) bool {
		p.pending.Delete(key)
		return true
	})

	// 终止Python进程
	return p.cmd.Process.Kill()
}

// IsRunning 检查Python进程是否仍在运行
func (p *PythonProcess) IsRunning() bool {
	if p.cmd == nil || p.cmd.Process == nil {
		return false
	}

	// 检查进程状态
	processState := p.cmd.ProcessState
	if processState != nil && processState.Exited() {
		return false
	}

	// 尝试发送信号检查进程是否存活
	err := p.cmd.Process.Signal(os.Interrupt)
	return err == nil
}

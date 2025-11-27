package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"os/exec"
	"sync"
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
			var resp map[string]string
			json.Unmarshal(scanner.Bytes(), &resp)

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

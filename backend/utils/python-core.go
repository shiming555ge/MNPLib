package utils

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"sync"

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

func NewPythonProcess(pythonPath string, pythonFile string) *PythonProcess {
	cmd := exec.Command(pythonPath, pythonFile)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	p := &PythonProcess{
		cmd:       cmd,
		writer:    bufio.NewWriter(stdin),
		sendQueue: make(chan *Request, 100),
	}

	// 启动 Python
	cmd.Start()

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

	return p
}

func (p *PythonProcess) SendAndWait(msg string) (string, error) {
	req := &Request{
		ID:      uuid.New().String(),
		Content: msg,
		Resp:    make(chan string),
	}

	p.sendQueue <- req
	res := <-req.Resp // <-- await here
	return res, nil
}

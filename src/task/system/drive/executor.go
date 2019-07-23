package drive

import (
	"html"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// 任务执行
type Executor struct {
}

var G_executor *Executor

func InitExecutor() {
	G_executor = &Executor{}
	return
}

// 执行任务
func (executor *Executor) ExecTask(task_info *TaskExecuteInfo) {
	go func() {
		result := &TaskExecuteResult{
			ExecuteInfo: task_info,
		}

		// 执行超时处理
		go func() {
			select {
			case <-task_info.CancelCtx.Done():
				task_info.CancelFunc()
			}
		}()

		result.StartTime = time.Now()

		switch task_info.Task.TaskExecType {
		case "http":
			var (
				http_request *http.Response
				err          error
			)

			httpClient := &http.Client{
				Timeout: 5 * time.Second,
			}

			switch task_info.Task.HttpType {
			case "post":
				http_request, err = httpClient.Post(task_info.Task.Cmd, "application/x-www-form-urlencoded", strings.NewReader(""))
				break
			case "get":
				http_request, err = httpClient.Get(task_info.Task.Cmd)
				break
			}

			result.EndTime = time.Now()

			defer func() {
				http_request.Body.Close()
			}()

			if err != nil {
				result.Err = err
				break
			}

			// 获取请求返回数据
			http_res, err := ioutil.ReadAll(http_request.Body)
			if err != nil {
				result.Err = err
				break
			}

			result.Output = []byte(html.EscapeString(string(http_res)))

			// TODO http请求任务，失败重试没有解决，留作后面解决。

			break
		case "shell":
			cmd := exec.CommandContext(task_info.CancelCtx, "/bin/bash", "-c", task_info.Task.Cmd)
			cmd_output, err := cmd.CombinedOutput()

			result.EndTime = time.Now()
			result.Output = cmd_output
			result.Err = err

			break
		}

		G_scheduler.PushTaskResult(result)
	}()
}

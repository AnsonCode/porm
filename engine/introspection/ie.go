package introspection

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func NewIntrospectEngine(enginePath string) *IntrospectEngine {
	// TODO:这里可以设置默认值
	engine := &IntrospectEngine{
		enginePath: enginePath,
	}
	_, err := engine.ensure() //确保引擎一定安装了
	if err != nil {
		panic(err)
	}
	return engine
}

type IntrospectEngine struct {
	enginePath string
}

func (e *IntrospectEngine) ensure() (string, error) {
	ensureEngine := time.Now()
	if _, err := os.Stat(e.enginePath); err != nil {
		return "", fmt.Errorf("no binary found ")
	}
	fmt.Println(" engine found in global path")

	fmt.Printf("using introspection engine at %s", e.enginePath)
	fmt.Printf("ensure query engine took %s", time.Since(ensureEngine))

	return e.enginePath, nil
}

func (e *IntrospectEngine) Pull(provider, url string) (string, error) {
	schema := fmt.Sprintf(
		`datasource db {
			provider = "%s"
			url      = "%s"
		}`,
		provider, url)
	startParse := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel() // 读取一行数据后，发送kill信号

	cmd := exec.CommandContext(ctx, e.enginePath)

	pipe, err := cmd.StdinPipe() // 标准输入流
	if err != nil {
		return "", fmt.Errorf("introspect engine std in pipe %v", err.Error())
	}
	defer pipe.Close()
	// 构建一个json-rpc 请求参数
	req := IntrospectRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "introspect",
		Params: []map[string]interface{}{
			{
				"schema":             string(schema),
				"compositeTypeDepth": -1,
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	// 入参追加到管道中
	_, err = pipe.Write(append(data, []byte("\n")...))
	if err != nil {
		return "", err
	}
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println(err)
		return "", err
	}

	// 不阻塞启动
	if err := cmd.Start(); err != nil {
		return "", err
	}

	reader := bufio.NewReader(stdout)

	// TODO:如果一直堵死在这咋办？
	//阻塞读取，实时读取输出流中的一行内容
	line, err2 := reader.ReadString('\n')
	if err2 != nil || io.EOF == err2 {
		return "", err2
	}
	log.Println(line)

	var response IntrospectResponse

	// 解析响应结果
	err = json.Unmarshal([]byte(line), &response)
	if err != nil {
		return "", err
	}

	log.Printf("[timing] introspect took %s", time.Since(startParse))
	if response.Error != nil {
		return "", fmt.Errorf("introspect error: %s", response.Error.Data.Message)
	}
	log.Println("introspect successful")

	dataModel := strings.Replace(response.Result.DataModel, " Bytes", " String", -1)
	//dataModel := strings.Replace(response.Result.DataModel, " Bytes", " String", -1)
	return dataModel, nil
}

func (e *IntrospectEngine) Pull2(provider, url string) (string, error) {
	schema := fmt.Sprintf(
		`datasource db {
			provider = "%s"
			url      = "%s"
		}`,
		provider, url)
	startParse := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.enginePath)

	pipe, err := cmd.StdinPipe() // 标准输入流
	if err != nil {
		return "", fmt.Errorf("introspect engine std in pipe %v", err.Error())
	}
	defer pipe.Close()
	// 构建一个json-rpc 请求参数
	req := IntrospectRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "introspect",
		Params: []map[string]interface{}{
			{
				"schema":             string(schema),
				"compositeTypeDepth": -1,
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	// 入参追加到管道中
	_, err = pipe.Write(append(data, []byte("\n")...))
	if err != nil {
		// return "", err
		return "", err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("introspect std out pipe %s ", err.Error())
		return "", err
	}
	r := bufio.NewReader(out)

	// TODO:执行后没有释放进程 ps -ef |grep prisma
	// 开始执行
	err = cmd.Start()
	if err != nil {
		return "", err
	}

	//var response IntrospectResponse
	var response IntrospectResponse

	outBuf := &bytes.Buffer{}
	// 这一段的意思是，每100ms读取一次结果，直到超时或有结果
	for {
		// 等待100 ms
		//time.Sleep(time.Millisecond * 100)
		b, err := r.ReadByte()
		if err != nil {
			err = fmt.Errorf("Introspect ReadByte %s ", err.Error())
		}
		err = outBuf.WriteByte(b)
		if err != nil {
			err = fmt.Errorf("IntrospectwriteByte %s ", err.Error())
		}

		if b == '\n' {
			// 解析响应结果
			err = json.Unmarshal(outBuf.Bytes(), &response)
			if err != nil {
				return "", err
			}
			if response.Error == nil {
				log.Println("introspect successful")
			}
			fmt.Print("ende ")
			break
		}
		// 如果超时了？跳出读取？
		if err := ctx.Err(); err != nil {
			return "", err
		}
	}
	log.Printf("[timing] introspect took %s", time.Since(startParse))
	if response.Error != nil {
		return "", fmt.Errorf("introspect error: %s", response.Error.Data.Message)
	}
	dataModel := strings.Replace(response.Result.DataModel, " Bytes", " String", -1)
	//dataModel := strings.Replace(response.Result.DataModel, " Bytes", " String", -1)
	return dataModel, nil
}

package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func NewQueryEngine(schema, port, enginePath string) *QueryEngine {
	killExistingPrismaQueryEngineProcess(port)

	return &QueryEngine{
		Schema:     schema,
		port:       port,
		enginePath: enginePath,
		// hasBinaryTargets: hasBinaryTargets,
		http: &http.Client{},
	}
}

type QueryEngine struct {
	// cmd holds the prisma binary process
	cmd *exec.Cmd

	enginePath string

	// http is the internal http client
	http *http.Client

	port string

	// url holds the query-engine url
	// url string

	// Schema contains the prisma Schema
	Schema string

	// hasBinaryTargets can be toggled by generated code from Schema.prisma whether binaryTargets
	// were specified and thus expects binaries in the local path
	// hasBinaryTargets bool

	disconnected bool
}

func (e *QueryEngine) Name() string {
	return "query-engine"
}
func (e *QueryEngine) Connect() error {
	fmt.Printf("ensure query engine binary...")

	_ = godotenv.Load(".env")
	// _ = godotenv.Load("db/e2e.env")
	// _ = godotenv.Load("prisma/e2e.env")

	startEngine := time.Now()

	file, err := e.ensure()
	if err != nil {
		return fmt.Errorf("ensure: %w", err)
	}

	if err := e.spawn(file); err != nil {
		return fmt.Errorf("spawn: %w", err)
	}

	fmt.Printf("connecting took %s", time.Since(startEngine))
	fmt.Printf("connected.")

	return nil
}

func (e *QueryEngine) Disconnect() error {
	e.disconnected = true
	fmt.Printf("disconnecting...")

	if runtime.GOOS == "windows" {
		if err := e.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("kill process: %w", err)
		}
		return nil
	}

	if err := e.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("send signal: %w", err)
	}

	if err := e.cmd.Wait(); err != nil {
		if err.Error() != "signal: interrupt" {
			return fmt.Errorf("wait for process: %w", err)
		}
	}

	fmt.Println("disconnected.")
	return nil
}

func (e *QueryEngine) ensure() (string, error) {
	ensureEngine := time.Now()
	if _, err := os.Stat(e.enginePath); err != nil {
		return "", fmt.Errorf("no binary found ")
	}
	fmt.Println("query engine found in global path")

	fmt.Printf("using query engine at %s", e.enginePath)
	fmt.Printf("ensure query engine took %s", time.Since(ensureEngine))

	return e.enginePath, nil
}

func (e *QueryEngine) spawn(file string) error {

	fmt.Printf("running query-engine on port %s", e.port)

	e.cmd = exec.Command(file, "-p", e.port, "--enable-raw-queries", "--enable-playground")
	// args = append(args, "--enable-playground", "--port", queryEnginePort)

	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr

	e.cmd.Env = append(
		os.Environ(),
		"PRISMA_DML="+e.Schema,
		"RUST_LOG=error",
		"RUST_LOG_FORMAT=json",
		"PRISMA_CLIENT_ENGINE_TYPE=binary",
	)

	// TODO fine tune this using log levels
	// if logger.Enabled {
	e.cmd.Env = append(
		e.cmd.Env,
		"PRISMA_LOG_QUERIES=y",
		"RUST_LOG=info",
	)
	// }

	fmt.Printf("starting engine...")

	if err := e.cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	fmt.Printf("connecting to engine...")

	ctx := context.Background()
	// send a basic readiness healthcheck and retry if unsuccessful
	var connectErr error
	for i := 0; i < 100; i++ {
		err := e.Ping(ctx)
		if err != nil {
			connectErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
	if connectErr != nil {
		return fmt.Errorf("readiness query error: %w", connectErr)
	}
	return nil
}

func (e *QueryEngine) Ping(ctx context.Context) error {
	body, err := e.Request(ctx, "GET", "/status", map[string]interface{}{})
	if err != nil {
		fmt.Printf("could not connect; retrying...")
		return err
	}

	var response GQLResponse

	if err := json.Unmarshal(body, &response); err != nil {

		fmt.Printf("could not unmarshal response; retrying...")
		return err
	}

	if response.Errors != nil {
		fmt.Println("could not connect due to gql errors; retrying...", response.Errors)
		return fmt.Errorf("readiness gql errors: %+v", response.Errors)
	}
	return nil
}

func (e *QueryEngine) Kill() error {
	killExistingPrismaQueryEngineProcess(e.port)
	return nil
}

// reference:https://github.com/wundergraph/wundergraph
func killExistingPrismaQueryEngineProcess(queryEnginePort string) {
	if runtime.GOOS == "windows" {
		command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess -Force", queryEnginePort)
		execCmd(exec.Command("Stop-Process", "-Id", command))
	} else {
		command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", queryEnginePort)
		execCmd(exec.Command("bash", "-c", command))
	}
}

func execCmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			log.Println("Error during port killing (exit code: )", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

		log.Println("Successfully killed existing prisma query process", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"sync"
)

// test:  docker run -d --name test_log  alpine sh -c 'for i in `seq 10000`; do echo $i; sleep 1; done'

func saveLog(container string, dstFilename string, stdout bool, stderr bool, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	options := types.ContainerLogsOptions{Follow: true}
	if stdout {
		options.ShowStdout = true
	}
	if stderr {
		options.ShowStderr = true
	}

	reader, err := cli.ContainerLogs(ctx, container, options)
	if err != nil {
		panic(err)
	}
	// noinspection GoUnhandledErrorResult
	defer reader.Close()

	f, err := os.Create(dstFilename)
	if err != nil {
		panic(err)
	}
	// noinspection GoUnhandledErrorResult
	defer f.Close()

	_, err = io.Copy(f, reader)
	if err != nil {
		panic(err)
	}
}

func main() {
	container := os.Args[1]
	stdoutFilename := os.Args[2]
	stderrFilename := os.Args[3]
	fmt.Println("container", container)
	fmt.Println("out", stdoutFilename)
	fmt.Println("err", stderrFilename)
	var wg sync.WaitGroup
	wg.Add(2)
	go saveLog(container, stdoutFilename, true, false, &wg)
	go saveLog(container, stderrFilename, false, true, &wg)

	wg.Wait()
}

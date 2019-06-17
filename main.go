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

	options := types.ContainerLogsOptions{}
	if stdout {
		options = types.ContainerLogsOptions{
			ShowStdout: true,
			Follow:     true,
		}
	}
	if stderr {
		options = types.ContainerLogsOptions{
			ShowStderr: true,
			Follow:     true,
		}
	}

	reader, err := cli.ContainerLogs(ctx, container, options)
	if err != nil {
		panic(err)
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)

	f, err := os.Create(dstFilename)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

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

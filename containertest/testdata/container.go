package testdata

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

// https://golang.testcontainers.org/features/docker_compose/
// https://golang.testcontainers.org/modules/redis/

func UpContainer(removeData bool) (downContainer func() error, err error) {
	compose, err := newDockerCompose()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	err = compose.Up(ctx, tc.RemoveOrphans(true))
	if err != nil {
		return nil, fmt.Errorf("compose.Up: %w", err)
	}

	down := func() error {
		err := compose.Down(ctx, tc.RemoveVolumes(removeData), tc.RemoveOrphans(true))
		if err != nil {
			return fmt.Errorf("compose.Down: %w", err)
		}
		return nil
	}
	defer func() {
		if err != nil {
			downErr := down()
			if downErr != nil {
				panic(err)
			}
		}
	}()

	setupServices := services()
	mqError := make(chan error, len(setupServices))
	wg := sync.WaitGroup{}
	for _, setup := range setupServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc, err := setup(compose, ctx)
			if err != nil {
				mqError <- fmt.Errorf("svc=%v: setup service: %w", svc, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		mqError <- nil
	}()

	return down, <-mqError
}

//

func newDockerCompose() (tc.ComposeStack, error) {
	// 嘗試找正確的 docker-compose.yml 路徑

	workDirs := []func() (string, error){
		func() (string, error) {
			wd, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("os.Getwd: %w", err)
			}
			return wd, nil
		},

		func() (string, error) {
			wd, ok := os.LookupEnv("WORK_DIR")
			if !ok {
				return "", errors.New("ENV ${WORK_DIR} not set")
			}
			return wd, nil
		},
	}

	var Err error
	for i := 0; i < len(workDirs); i++ {
		workDir, err := workDirs[i]()
		if err != nil {
			Err = err
			continue
		}

		files := []string{
			filepath.Join(workDir, "containertest", "testdata", "docker-compose.yml"),
			filepath.Join(workDir, "docker-compose.yml"),
		}
		for _, f := range files {
			_, err := os.Stat(f)
			if err != nil {
				Err = err
				continue
			}
			compose, err := tc.NewDockerComposeWith(tc.WithStackFiles(f))
			if err != nil {
				Err = err
				continue
			}
			return compose, nil
		}
	}

	return nil, Err
}

//

type service func(compose tc.ComposeStack, ctx context.Context) (svc string, err error)

func services() []service {
	return []service{
		redisService("redis1"),
	}
}

func redisService(svc string) service {
	return func(compose tc.ComposeStack, ctx context.Context) (string, error) {
		container, err := compose.ServiceContainer(ctx, svc)
		if err != nil {
			return svc, err
		}

		const redisLogMessage = "Ready to accept connections tcp"
		waitStrategy := wait.ForAll(
			wait.ForLog(redisLogMessage).WithStartupTimeout(time.Second),
		)
		err = waitStrategy.WaitUntilReady(ctx, container)
		if err != nil {
			return svc, fmt.Errorf("wait ready: %w", err)
		}

		code, reader, err := container.Exec(ctx, []string{"bash", "-c", "cat /testdata | redis-cli --pipe"})
		if err != nil {
			return svc, fmt.Errorf("container.Exec: %w", err)
		}
		const success = 0
		if code != success {
			errorMessage, _ := io.ReadAll(reader)
			return svc, fmt.Errorf("populate testdata: \n%v", string(errorMessage))
		}

		host, err := container.Host(ctx)
		if err != nil {
			return svc, err
		}
		port, err := container.MappedPort(ctx, "6379")
		if err != nil {
			return svc, err
		}

		fmt.Printf("Redis server listening on %v %v\n", host, port.Port())

		return svc, nil
	}
}

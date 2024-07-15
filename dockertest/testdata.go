package main

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// https://github.com/testcontainers/testcontainers-go
// https://golang.testcontainers.org/features/docker_compose/

var (
	//go:embed testdata.redis
	testdataRedis string
)

func main() {
	RunDockerTest()
}

func RunDockerTest() {
	container := NewDockerTest()
	defer container.Purge()

	container.DefaultRedis()
	log.Printf("redis port %v", container.RedisPublishedPort)

	container.Wait()
}

func NewDockerTest() *DockerTest {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	err = pool.Client.Ping()
	if err != nil {
		panic(err)
	}

	container := &DockerTest{
		pool:       pool,
		purgeQueue: make([]*dockertest.Resource, 0),
		done:       make(chan struct{}),
	}

	go func() {
		osSig := make(chan os.Signal, 2)
		signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
		<-osSig
		close(container.done)
	}()

	return container
}

// https://github.com/ory/dockertest
type DockerTest struct {
	pool    *dockertest.Pool
	authOpt docker.AuthConfiguration

	RedisResource      *dockertest.Resource
	RedisPublishedPort string

	purgeQueue []*dockertest.Resource
	done       chan struct{}
}

func (c *DockerTest) Wait() {
	<-c.done
}

func (c *DockerTest) Purge() {
	for _, r := range c.purgeQueue {
		if err := c.pool.Purge(r); err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func (c *DockerTest) DefaultRedis() {
	c.Redis(c.DefaultTask())
}

func (c *DockerTest) DefaultTask() func(*dockertest.Resource) error {
	return func(r *dockertest.Resource) error {
		execOpts := dockertest.ExecOptions{
			StdIn: bytes.NewBufferString(testdataRedis),
		}
		_, err := r.Exec([]string{"redis-cli"}, execOpts)
		return err
	}
}

func (c *DockerTest) Redis(tasks ...func(r *dockertest.Resource) error) {
	r, port := c.redis()
	c.RedisResource = r
	c.RedisPublishedPort = port
	c.purgeQueue = append(c.purgeQueue, r)

	err := c.pool.Retry(func() error {
		for _, task := range tasks {
			err := task(r)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (c *DockerTest) redis() (*dockertest.Resource, string) {
	// parameter
	Repository := "redis"
	Tag := "7.2.5"
	Env := []string{}
	ExposedPort := "6379/tcp"
	Cmd := []string{}

	// workflow
	{
		imageOpt := docker.PullImageOptions{
			Repository: Repository,
			Tag:        Tag,
		}
		err := c.pool.Client.PullImage(imageOpt, c.authOpt)
		if err != nil {
			panic(err)
		}

		runOpt := &dockertest.RunOptions{
			Repository:   Repository,
			Tag:          Tag,
			Env:          Env,
			Cmd:          Cmd,
			ExposedPorts: []string{ExposedPort},
			Auth:         c.authOpt,
		}
		resource, err := c.pool.RunWithOptions(runOpt)
		if err != nil {
			panic(err)
		}

		return resource, resource.GetPort(ExposedPort)
	}
}

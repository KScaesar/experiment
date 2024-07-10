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

	defer func() {
		if recover() != nil {
			container.SafePurge()
		}
	}()

	container.Redis()
	log.Printf("redis port %v", container.RedisPublishedPort)

	container.WaitPurge()
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

func (c *DockerTest) SafePurge() {
	close(c.done)
	c.WaitPurge()
}

func (c *DockerTest) WaitPurge() {
	<-c.done
	for _, r := range c.purgeQueue {
		if err := c.pool.Purge(r); err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func (c *DockerTest) Redis() {
	r, port := c.redis()
	c.RedisResource = r
	c.RedisPublishedPort = port
	c.purgeQueue = append(c.purgeQueue, r)
}

func (c *DockerTest) redis() (*dockertest.Resource, string) {
	// parameter
	Repository := "redis"
	Tag := "7.2.5"
	Env := []string{}
	ExposedPort := "6379/tcp"
	Cmd := []string{}
	Task := func(r *dockertest.Resource) error {
		execOpts := dockertest.ExecOptions{
			StdIn: bytes.NewBufferString(testdataRedis),
		}
		_, err := r.Exec([]string{"redis-cli"}, execOpts)
		return err
	}

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

		err = c.pool.Retry(func() error {
			return Task(resource)
		})
		if err != nil {
			panic(err)
		}

		return resource, resource.GetPort(ExposedPort)
	}
}

package container

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Restart(names []string) {
	if len(names) == 0 {
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("failed to create docker client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, name := range names {
		if err := cli.ContainerRestart(ctx, name, container.StopOptions{}); err != nil {
			log.Printf("failed to restart container %q: %v", name, err)
		} else {
			log.Printf("restarted container %q", name)
		}
	}
}

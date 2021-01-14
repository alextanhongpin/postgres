package postgres

import "fmt"

// DockerImage struct.
type DockerImage struct {
	Container string
	Version   string
}

func (d DockerImage) String() string {
	return fmt.Sprintf("%s:%s", d.Container, d.Version)
}

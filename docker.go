package postgres

import "fmt"

// Container struct.
type Container struct {
	Image   string
	Version string
}

func (c Container) String() string {
	return fmt.Sprintf("%s:%s", c.Image, c.Version)
}

package commands

import "fmt"

type Commands interface {
	Name() string
	HandleCommand(event string)
}

package commands

type Commands interface {
	Type() string // Returns a string description of the CommandHandler
	HandleCommand(event string)
}


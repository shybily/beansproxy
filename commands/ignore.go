package commands

type IgnoreCommand struct {
	baseCommand
}

func NewIgnoreCommand(cmdLine []byte) (*IgnoreCommand, error) {
	return &IgnoreCommand{baseCommand{cmdLine: cmdLine}}, nil
}

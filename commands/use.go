package commands

type UseCommand struct {
	baseCommand
}

func NewUseCommand(cmdLine []byte) (*UseCommand, error) {
	return &UseCommand{baseCommand{cmdLine: cmdLine}}, nil
}

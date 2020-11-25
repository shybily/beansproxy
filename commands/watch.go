package commands

type WatchCommand struct {
	baseCommand
}

func NewWatchCommand(cmdLine []byte) (*WatchCommand, error) {
	return &WatchCommand{baseCommand{cmdLine: cmdLine}}, nil
}

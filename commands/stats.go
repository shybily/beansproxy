package commands

type StatsCommand struct {
	baseCommand
}

func NewStatsCommand(cmdLine []byte) (*StatsCommand, error) {
	return &StatsCommand{baseCommand{cmdLine: cmdLine}}, nil
}

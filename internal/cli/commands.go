package cli

func commandIndex() []command {
	var commands []command
	commands = append(commands, indexingCommands()...)
	commands = append(commands, gridCommands()...)
	commands = append(commands, hierarchyCommands()...)
	commands = append(commands, regionCommands()...)
	commands = append(commands, edgeCommands()...)
	commands = append(commands, vertexCommands()...)
	commands = append(commands, miscCommands()...)
	return commands
}

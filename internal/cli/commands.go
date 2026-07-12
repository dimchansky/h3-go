package cli

// commandIndex assembles the full command registry. The concatenation order
// is observable — it defines the listing order in `h3 --help`, which the
// upstream CLI fixes by registration order — so groups must stay in this
// sequence even though lookup itself is order-independent.
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

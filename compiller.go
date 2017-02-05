package main

import "os/exec"

// pluginCompiler
type pluginCompiler struct {
	pluginName string
	sourcePath string
}

// NewPluginCompiler creates and returns new pluginCompiler
func NewPluginCompiler(plugName, pathToSource string) *pluginCompiler {
	return &pluginCompiler{
		pluginName: plugName,
		sourcePath: pathToSource,
	}
}

func (p *pluginCompiler) BuildPlugin() error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", p.pluginName+".so", p.sourcePath)
	return cmd.Run()
}

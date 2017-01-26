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
	cmd := exec.Command("go", p.sourcePath, "build", "-buildmode=plugin", "-o", p.pluginName+".so")
	return cmd.Run()
}

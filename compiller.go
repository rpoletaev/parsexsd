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
		pluginName: plugName + ".so",
		sourcePath: pathToSource,
	}
}

func (p *pluginCompiler) BuildPlugin() error {
	println("pluginName ", p.pluginName)
	println("sourcePath ", p.sourcePath)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", p.pluginName, p.sourcePath)
	return cmd.Run()
}

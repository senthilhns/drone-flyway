// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Args struct {
	FlywayEnvPluginArgs
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`
}

type FlywayEnvPluginArgs struct {
	DriverPath      string `envconfig:"PLUGIN_DRIVER_PATH"`
	FlywayCommand   string `envconfig:"PLUGIN_FLYWAY_COMMAND"`
	Locations       string `envconfig:"PLUGIN_LOCATIONS"`
	CommandLineArgs string `envconfig:"PLUGIN_COMMAND_LINE_ARGS"`
	Url             string `envconfig:"PLUGIN_URL"`
	UserName        string `envconfig:"PLUGIN_USERNAME"`
	Password        string `envconfig:"PLUGIN_PASSWORD"`
	IsDryRun        string `envconfig:"PLUGIN_IS_DRY_RUN"`
}

type FlywayPlugin struct {
	InputArgs         *Args
	IsMultiFileUpload bool
	ProcessingInfo
}

func (a *Args) ToStr() string {
	jsonData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Printf("Error marshalling Args to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}

type ProcessingInfo struct {
	ExecCommand         string
	Env                 string
	CommandSpecificArgs string
}

func Exec(ctx context.Context, args Args) (FlywayPlugin, error) {
	plugin := FlywayPlugin{}

	err := plugin.Init(&args)
	if err != nil {
		return plugin, err
	}
	defer func(p FlywayPlugin) {
		err := p.DeInit()
		if err != nil {
			LogPrintln("Error in DeInit: " + err.Error())
		}
	}(plugin)

	err = plugin.ValidateAndProcessArgs(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.DoPostArgsValidationSetup(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.Run()
	if err != nil {
		return plugin, err
	}

	return plugin, nil
}

func (p *FlywayPlugin) Init(args *Args) error {
	p.InputArgs = args
	return nil
}

func (p *FlywayPlugin) SetBuildRoot(buildRootPath string) error {
	return nil
}

func (p *FlywayPlugin) DeInit() error {
	return nil
}

func (p *FlywayPlugin) ValidateAndProcessArgs(args Args) error {
	err := p.IsCommandValid()
	if err != nil {
		LogPrintln(p, err.Error())
		return err
	}
	return nil
}

func (p *FlywayPlugin) DoPostArgsValidationSetup(args Args) error {
	if args.FlywayCommand == CleanCommand {
		if !strings.Contains(args.CommandLineArgs, "-cleanDisabled") {
			p.CommandSpecificArgs = "-cleanDisabled=false" + " "
		}
	}
	return nil
}

func (p *FlywayPlugin) Run() error {
	var stdoutBuf, stderrBuf bytes.Buffer
	var err error

	p.ExecCommand = p.GetExecArgsStr()
	fmt.Println("Command: ", p.ExecCommand)

	cmdParts := strings.Fields(p.ExecCommand)
	if len(cmdParts) < 2 {
		return fmt.Errorf("Invalid command: %s", p.ExecCommand)
	}
	cmdName := cmdParts[0]
	cmdArgs := cmdParts[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if len(p.InputArgs.DriverPath) > 0 {
		p.Env = "CLASSPATH=" + p.InputArgs.DriverPath
		cmd.Env = append(os.Environ(), "CLASSPATH="+p.InputArgs.DriverPath)
	}

	if p.InputArgs.IsDryRun == "TRUE" {
		return nil
	}

	err = cmd.Run()
	if err == nil {
		fmt.Println(stdoutBuf.String())
		fmt.Printf("Command execution success")
	} else {
		fmt.Println(stderrBuf.String())
		fmt.Printf("Error executing command: %v\n", err.Error())
	}

	return nil
}

func (p *FlywayPlugin) GetExecArgsStr() string {
	var execCommand string

	execCommand += GetFlywayExecutablePath() + " "
	execCommand += p.InputArgs.FlywayCommand + " "
	execCommand += p.CommandSpecificArgs + " "

	if len(p.InputArgs.Url) > 0 {
		execCommand += "-url=" + p.InputArgs.Url + " "
	}
	if len(p.InputArgs.UserName) > 0 {
		execCommand += "-user=" + p.InputArgs.UserName + " "
	}
	if len(p.InputArgs.Password) > 0 {
		execCommand += "-password=" + p.InputArgs.Password + " "
	}
	if len(p.InputArgs.Locations) > 0 {
		execCommand += "-locations=" + p.InputArgs.Locations + " "
	}
	// this should be the last
	execCommand += p.InputArgs.CommandLineArgs

	return execCommand
}

func (p *FlywayPlugin) IsCommandValid() error {
	if p.InputArgs.FlywayCommand == "" {
		return fmt.Errorf("Command is empty")
	}

	err := p.IsUnknownCommand()
	if err != nil {
		return err
	}

	err = p.CheckMandatoryArgs()
	if err != nil {
		return err
	}

	return nil
}

func (p *FlywayPlugin) IsUnknownCommand() error {
	_, ok := knownCommandsMap[p.InputArgs.FlywayCommand]
	if ok {
		return nil
	}

	return fmt.Errorf("Unknown command: %s", p.InputArgs.FlywayCommand)
}

func (p *FlywayPlugin) CheckMandatoryArgs() error {

	args := p.InputArgs

	if strings.Contains(args.CommandLineArgs, ConfigFileOpt) { // pick args from file
		return nil
	}

	type mandatoryArg struct {
		EnvName   string
		ParamName *string
		Hint      string
	}

	ma := []mandatoryArg{
		{"FLYWAY_URL", &args.Url, "url"},
		{"FLYWAY_USER", &args.UserName, "username"},
		{"FLYWAY_PASSWORD", &args.Password, "password"},
	}

	for _, m := range ma {
		if os.Getenv(m.EnvName) == "" && *m.ParamName == "" {
			LogPrintln("Missing mandatory argument: " + m.EnvName)
			return fmt.Errorf("Missing mandatory argument: %s", m.Hint)
		}
	}

	return nil
}

var knownCommandsMap = map[string]bool{
	MigrateCommand:  true,
	CleanCommand:    true,
	BaselineCommand: true,
	RepairCommand:   true,
	ValidateCommand: true,
}

const (
	MigrateCommand  = "migrate"
	CleanCommand    = "clean"
	BaselineCommand = "baseline"
	RepairCommand   = "repair"
	ValidateCommand = "validate"
	ConfigFileOpt   = "-configFiles"
)

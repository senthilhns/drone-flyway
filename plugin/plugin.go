// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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
	IsDryRun        bool   `envconfig:"PLUGIN_IS_DRY_RUN"`
}

type FlywayPlugin struct {
	InputArgs         *Args
	IsMultiFileUpload bool
	ProcessingInfo
}

type ProcessingInfo struct {
	ExecCommand         string
	Env                 string
	CommandSpecificArgs string
}

func (p FlywayPlugin) ToString() string {
	jsonStr, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("FlywayPlugin: %v", p)
	}
	return string(jsonStr)
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
	logrus.Infof("Executing command: %s", strings.ReplaceAll(p.ExecCommand, p.InputArgs.Password, "******"))
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

	if p.InputArgs.IsDryRun {
		return nil
	}

	err = cmd.Run()
	if err == nil {
		logrus.Println(stdoutBuf.String())
		logrus.Infof("Command executed successfully: %s", stdoutBuf.String())
	} else {
		logrus.Println(stderrBuf.String())
		logrus.Errorf("Command execution failed: %s, error: %v", stderrBuf.String(), err)
	}

	return nil
}

func (p *FlywayPlugin) GetExecArgsStr() string {
	var builder strings.Builder

	builder.WriteString(GetFlywayExecutablePath() + " ")
	builder.WriteString(p.InputArgs.FlywayCommand + " ")
	builder.WriteString(p.CommandSpecificArgs + " ")

	if p.InputArgs.Url != "" {
		builder.WriteString("-url=" + p.InputArgs.Url + " ")
	}
	if p.InputArgs.UserName != "" {
		builder.WriteString("-user=" + p.InputArgs.UserName + " ")
	}
	if p.InputArgs.Password != "" {
		builder.WriteString("-password=" + p.InputArgs.Password + " ")
	}
	if p.InputArgs.Locations != "" {
		builder.WriteString("-locations=" + p.InputArgs.Locations + " ")
	}
	// this should be the last
	builder.WriteString(p.InputArgs.CommandLineArgs)

	return builder.String()
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
	if _, ok := knownCommandsMap[p.InputArgs.FlywayCommand]; !ok {
		return fmt.Errorf("unknown command: %s", p.InputArgs.FlywayCommand)
	}
	return nil
}

func (p *FlywayPlugin) CheckMandatoryArgs() error {

	args := p.InputArgs

	if strings.Contains(args.CommandLineArgs, ConfigFileOpt) { // pick args from file
		return nil
	}

	type mandatoryArg struct {
		EnvName   string
		ParamName string
		Hint      string
	}

	ma := []mandatoryArg{
		{"FLYWAY_URL", args.Url, "url"},
		{"FLYWAY_USER", args.UserName, "username"},
		{"FLYWAY_PASSWORD", args.Password, "password"},
	}

	for _, m := range ma {
		if os.Getenv(m.EnvName) == "" && m.ParamName == "" {
			LogPrintln("Missing mandatory argument: " + m.EnvName)
			return fmt.Errorf("Missing mandatory argument: %s (env: %s)", m.Hint, m.EnvName)
		}
	}

	return nil
}

var knownCommandsMap = map[string]string{
	MigrateCommand:  "Performs database migration",
	CleanCommand:    "Drops all objects in the configured schemas",
	BaselineCommand: "Baselines an existing database",
	RepairCommand:   "Repairs the schema history table",
	ValidateCommand: "Validates the applied migrations against the available ones",
}

const (
	MigrateCommand  = "migrate"
	CleanCommand    = "clean"
	BaselineCommand = "baseline"
	RepairCommand   = "repair"
	ValidateCommand = "validate"
	ConfigFileOpt   = "-configFiles"
)

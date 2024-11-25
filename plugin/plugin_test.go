// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestUnitTcClean(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"clean",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		fmt.Println("Error in Exec: " + err.Error())
		t.Fail()
	}
	expectedCmd := " clean -cleanDisabled=false  -url=jdbc:mysql://3.4.9.2:3306/flyway_test " +
		"-user=hnstest03 -password=sk89sl2@3 -locations=filesystem:/test/db-migrate01 "

	if fp.ExecCommand != expectedCmd {
		t.Fail()
	}
}

func TestUnitTcBaseline(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"baseline",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		fmt.Println("Error in Exec: " + err.Error())
		t.Fail()
	}
	expectedCmd := " baseline  -url=jdbc:mysql://3.4.9.2:3306/flyway_test" +
		" -user=hnstest03 -password=sk89sl2@3 -locations=filesystem:/test/db-migrate01 "

	if fp.ExecCommand != expectedCmd {
		t.Fail()
	}
}

func TestUnitTcMigrate(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"migrate",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		t.Fail()
	}
	expectedCmd := " migrate  -url=jdbc:mysql://3.4.9.2:3306/flyway_test " +
		"-user=hnstest03 -password=sk89sl2@3 -locations=filesystem:/test/db-migrate01 "

	if fp.ExecCommand != expectedCmd {
		fmt.Printf("|%s|", fp.ExecCommand)
		t.Fail()
	}
}

func TestUnitTcRepair(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"repair",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		t.Fail()
	}
	expectedCmd := " repair  -url=jdbc:mysql://3.4.9.2:3306/flyway_test" +
		" -user=hnstest03 -password=sk89sl2@3 -locations=filesystem:/test/db-migrate01 "

	if fp.ExecCommand != expectedCmd {
		t.Fail()
	}
}

func TestUnitTcValidate(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"validate",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		fmt.Println("Error in Exec: " + err.Error())
		t.Fail()
	}
	expectedCmd := " validate  -url=jdbc:mysql://3.4.9.2:3306/flyway_test" +
		" -user=hnstest03 -password=sk89sl2@3 -locations=filesystem:/test/db-migrate01 "

	if fp.ExecCommand != expectedCmd {
		t.Fail()
	}
}

func TestUnitTcWithConfigFiles(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"migrate",
		"", // locations
		getDefaultPluginCommandLineArgs(),
		"", // url
		"", // username
		"", // password
	)
	args.CommandLineArgs = "-configFiles=/harness/hns/test-resources/flyway/config1/flyway.conf"

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		fmt.Println("Error in Exec: " + err.Error())
		t.Fail()
	}
	expectedCmd := " migrate  -configFiles=/harness/hns/test-resources/flyway/config1/flyway.conf"

	if fp.ExecCommand != expectedCmd {
		fmt.Printf("|%s|", fp.ExecCommand)
		t.Fail()
	}
}

func TestUnitTcWithDriverPath(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		"/harness/test/flyway-mysql-10.21.0.jar",
		"clean",
		"",
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	if err != nil {
		fmt.Println("Error in Exec: " + err.Error())
		t.Fail()
	}
	expectedCmd := " clean -cleanDisabled=false  -url=jdbc:mysql://3.4.9.2:3306/flyway_test -user=hnstest03 -password=sk89sl2@3 "
	fmt.Println(fp.Env)
	if fp.ExecCommand != expectedCmd {
		fmt.Printf("|%s|", fp.ExecCommand)
		t.Fail()
	}
	expectedEnv := "CLASSPATH=/harness/test/flyway-mysql-10.21.0.jar"
	if fp.Env != expectedEnv {
		t.Fail()
	}
}

func GetArgsForFunctionalTesting(pluginDriverPath, pluginFlywayCommand, pluginLocations,
	pluginCommandLineArgs, pluginUrl, pluginUser, pluginPassword string) Args {

	defaultArgs := Args{
		FlywayEnvPluginArgs: FlywayEnvPluginArgs{
			DriverPath:      pluginDriverPath,
			FlywayCommand:   pluginFlywayCommand,
			Locations:       pluginLocations,
			CommandLineArgs: pluginCommandLineArgs,
			Url:             pluginUrl,
			UserName:        pluginUser,
			Password:        pluginPassword,
			IsDryRun:        "TRUE",
		},
	}

	return defaultArgs
}

func getDefaultPluginDriverPath() string {
	return os.Getenv("PLUGIN_DRIVER_PATH")
}

func getDefaultPluginFlywayCommand() string {
	return os.Getenv("PLUGIN_FLYWAY_COMMAND")
}

func getDefaultPluginLocations() string {
	return "filesystem:/test/db-migrate01"
}

func getDefaultPluginCommandLineArgs() string {
	return os.Getenv("PLUGIN_COMMAND_LINE_ARGS")
}

func getDefaultPluginUrl() string {
	return "jdbc:mysql://3.4.9.2:3306/flyway_test"
}

func getDefaultPluginUser() string {
	return "hnstest03"
}

func getDefaultPluginPassword() string {
	return "sk89sl2@3"
}

// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const testUrl = "-url=jdbc:mysql://3.4.9.2:3306/flyway_test"
const testUser = "-user=hnstest03"
const testPassword = "-password=sk89sl2@3"
const testLocations = "-locations=filesystem:/test/db-migrate01"

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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" clean -cleanDisabled=false  %s "+
		"%s %s %s ", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" baseline  %s"+
		" %s %s %s ", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" migrate  %s "+
		"%s %s %s ", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" repair  %s"+
		" %s %s %s ", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" validate  %s"+
		" %s %s %s ", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := " migrate  -configFiles=/harness/hns/test-resources/flyway/config1/flyway.conf"

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)
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
	require.NoError(t, err, "Error in Exec: %v", err)

	expectedCmd := fmt.Sprintf(" clean -cleanDisabled=false  %s %s %s ", testUrl, testUser, testPassword)
	fmt.Println(fp.Env)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command. Got: %s", fp.ExecCommand)

	expectedEnv := "CLASSPATH=/harness/test/flyway-mysql-10.21.0.jar"
	require.Equal(t, expectedEnv, fp.Env, "Unexpected environment variable. Got: %s", fp.Env)
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
			IsDryRun:        true,
		},
	}

	return defaultArgs
}

func TestUnitTcMissingRequiredInputs(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"migrate",
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		"", // missing URL
		"", // missing username
		"", // missing password
	)

	_, err := Exec(context.TODO(), args)
	require.Error(t, err, "Exec should fail when required inputs (URL, username, or password) are missing")
}

func TestUnitTcWrongFlywayCommand(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"invalidCommand", // invalid Flyway command
		getDefaultPluginLocations(),
		getDefaultPluginCommandLineArgs(),
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	_, err := Exec(context.TODO(), args)
	require.Error(t, err, "Exec should fail for an invalid Flyway command")
}

func TestUnitTcWithExtraArgsVerbose(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"migrate",
		getDefaultPluginLocations(),
		"-X", // extra verbose argument
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	require.NoError(t, err, "Exec should succeed with extra verbose argument (-X)")

	expectedCmd := fmt.Sprintf(" migrate  %s "+
		"%s %s %s -X", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command with extra verbose argument. Got: %s", fp.ExecCommand)
}

func TestUnitTcWithExtraArgsQuiet(t *testing.T) {
	args := GetArgsForFunctionalTesting(
		getDefaultPluginDriverPath(),
		"migrate",
		getDefaultPluginLocations(),
		"-q", // extra quiet argument
		getDefaultPluginUrl(),
		getDefaultPluginUser(),
		getDefaultPluginPassword(),
	)

	fp, err := Exec(context.TODO(), args)
	require.NoError(t, err, "Exec should succeed with extra quiet argument (-q)")

	expectedCmd := fmt.Sprintf(" migrate  %s "+
		"%s %s %s -q", testUrl, testUser, testPassword, testLocations)

	require.Equal(t, expectedCmd, fp.ExecCommand, "Unexpected command with extra quiet argument. Got: %s", fp.ExecCommand)
}

func getDefaultPluginDriverPath() string {
	return ""
}

func getDefaultPluginLocations() string {
	return "filesystem:/test/db-migrate01"
}

func getDefaultPluginCommandLineArgs() string {
	return ""
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

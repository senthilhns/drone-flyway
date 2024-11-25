// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"github.com/sirupsen/logrus"
	"os"
)

func LogPrintln(args ...interface{}) {

	if !IsDevTestingMode() {
		return
	}

	logrus.Println(append([]interface{}{"Plugin Info:"}, args...)...)
}

func LogPrintf(format string, v ...interface{}) {

	if !IsDevTestingMode() {
		return
	}

	logrus.Printf(format, v...)
}

func IsDevTestingMode() bool {
	return os.Getenv("DEV_TEST_d6c9b463090c") == "true"
}

func GetFlywayExecutablePath() string {
	return os.Getenv("FLYWAY_BIN_PATH")
}

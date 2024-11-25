// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"github.com/sirupsen/logrus"
	"os"
)

func LogPrintln(args ...interface{}) {
	logrus.Println(append([]interface{}{"Plugin Info:"}, args...)...)
}

func GetFlywayExecutablePath() string {
	return os.Getenv("FLYWAY_BIN_PATH")
}

/*
Copyright © 2024 Mayron G Fernandes mayronfernandes01@gmail.com
*/
package main

import (
	"strings"

	"github.com/mayron1806/gengo/cmd"
	"github.com/mayron1806/gengo/config"
)

func main() {
	logger := config.GetLogger()
	logger.Info("Starting gengo...")

	defer logger.CloseLogger(&strings.Fields("exiting gengo")[0])
	cmd.Execute()
}

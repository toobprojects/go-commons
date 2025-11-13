package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/toobprojects/go-commons/logs"
	"github.com/toobprojects/go-commons/text"
)

const (
	TargetShellInterpreter = "/bin/bash"
	CommandErrorTag        = "Command Error : \n%s\n"
)

func Exec(command string, commandArgs []string, targetPath string, returnOutput bool) string {
	return execCommand(command, commandArgs, targetPath, returnOutput, false)
}

func ExecWithNativeLog(command string, commandArgs []string, targetPath string, returnOutput bool) string {
	return execCommand(command, commandArgs, targetPath, returnOutput, true)
}

func ExecScriptFile(scriptPath string, targetPath string, returnOutput bool) string {
	return Exec(TargetShellInterpreter, []string{scriptPath}, targetPath, returnOutput)
}

func execCommand(command string, commandArgs []string, targetPath string, returnOutput bool, logCommand bool) string {
	var responseOutput string
	if logCommand {
		logNativeCommand(command, commandArgs)
	}
	var cmd = exec.Command(command, commandArgs...)

	// Validate the input and see if there's something there
	if text.StringNotBlank(targetPath) {
		cmd.Dir = targetPath
	}

	// Set the TERM environment variable to support colors
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Run the command and return the output./to
	if returnOutput {
		result, err := cmd.CombinedOutput()
		if err != nil {
			logs.Error(fmt.Sprintf("Command failed: %s %v Error: %v Output: %s", command, commandArgs, err, string(result)))
		}

		responseOutput = string(result)

		// Otherwise just run the command and show in the console as it runs
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			logs.Error(fmt.Sprintf(CommandErrorTag, err))
			return text.EMPTY
		}
	}

	return responseOutput
}

func logNativeCommand(command string, commandArgs []string) {
	result := strings.Join(commandArgs, text.WHITE_SPACE)
	logs.Info(fmt.Sprintf("Running Native Command : %v %v", command, result))
}

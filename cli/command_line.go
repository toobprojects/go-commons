package cli

import (
	"context"
	"os"
	"os/exec"

	"github.com/toobprojects/go-commons/logs"
)

// Options defines how a command should be executed.
//
// This is designed to be reusable by any consumer of the go-commons module.
// Most callers can use RunWithDefaults or the simple Exec wrapper below.
type Options struct {
	// Dir sets the working directory for the command.
	// If empty, the current process working directory is used.
	Dir string

	// Env is a list of additional environment variables in KEY=VALUE form
	// to add on top of the inherited environment from the parent process.
	// If nil or empty, only the inherited environment is used.
	Env []string

	// CaptureOutput controls whether the command output is captured and
	// returned as a string, or streamed directly to Stdout/Stderr.
	//
	// - When true:  Run / RunWithDefaults return the combined output string.
	// - When false: Output is written to Stdout/Stderr and the returned
	//               string will be empty.
	CaptureOutput bool

	// Stdout is the destination for the command's standard output when
	// CaptureOutput is false. If nil, os.Stdout is used.
	Stdout *os.File

	// Stderr is the destination for the command's standard error when
	// CaptureOutput is false. If nil, os.Stderr is used.
	Stderr *os.File

	// LogCommand controls whether the executed command and its arguments
	// are logged before execution.
	LogCommand bool
}

// Run executes a command using the provided context, arguments and options.
//
// It returns the command's output (when CaptureOutput is true) and any error
// returned by the underlying exec.CommandContext invocation.
//
// This is the primary, reusable entry point for running native commands.
func Run(ctx context.Context, command string, args []string, opts Options) (string, error) {
	log := logs.WithGroup("cli").With("command", command)

	if opts.LogCommand {
		log.Info("Running native command",
			"command", command,
			"args", args,
			"dir", opts.Dir,
		)
	}

	cmd := exec.CommandContext(ctx, command, args...)

	// Working directory
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	// Environment
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	// Capture vs stream output
	if opts.CaptureOutput {
		out, err := cmd.CombinedOutput()
		output := string(out)

		if err != nil {
			logs.Error("Command failed",
				"args", args,
				"dir", opts.Dir,
				"err", err,
				"output", output,
			)
			return output, err
		}

		logs.Debug("Command succeeded",
			"args", args,
			"dir", opts.Dir,
		)

		return output, nil
	}

	// Streaming mode: attach stdout/stderr
	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		logs.Error("Command failed (streaming)",
			"args", args,
			"dir", opts.Dir,
			"err", err,
		)
		return "", err
	}

	logs.Debug("Command succeeded (streaming)",
		"args", args,
		"dir", opts.Dir,
	)

	return "", nil
}

// RunWithDefaults is a convenience helper for running a command with sensible defaults:
//
// - Inherits the current environment.
// - Uses the current working directory.
// - Captures and returns the combined output.
// - Does not log the command unless logCommand is true.
func RunWithDefaults(ctx context.Context, command string, args []string, logCommand bool) (string, error) {
	return Run(ctx, command, args, Options{
		CaptureOutput: true,
		LogCommand:    logCommand,
	})
}

// Exec is a backward-compatible wrapper for the original API.
//
// It runs the command synchronously, optionally returning the combined output.
// Errors are logged and an empty string is returned on failure.
//
// NOTE: For new code, prefer using Run or RunWithDefaults to get explicit error handling.
func Exec(command string, commandArgs []string, targetPath string, returnOutput bool) string {
	ctx := context.Background()

	opts := Options{
		Dir:           targetPath,
		CaptureOutput: returnOutput,
		LogCommand:    false,
	}

	out, err := Run(ctx, command, commandArgs, opts)
	if err != nil {
		// Error already logged by Run; return empty string for backward compatibility.
		return ""
	}
	return out
}

// ExecWithNativeLog is a compatibility wrapper that logs the command before execution.
func ExecWithNativeLog(command string, commandArgs []string, targetPath string, returnOutput bool) string {
	ctx := context.Background()

	opts := Options{
		Dir:           targetPath,
		CaptureOutput: returnOutput,
		LogCommand:    true,
	}

	out, err := Run(ctx, command, commandArgs, opts)
	if err != nil {
		return ""
	}
	return out
}

// ExecScriptFile is a helper for running an executable script file.
//
// It uses the system shell (`/bin/bash`) to run the script at scriptPath
// with the given target working directory.
func ExecScriptFile(scriptPath string, targetPath string, returnOutput bool) string {
	const shell = "/bin/bash"

	return Exec(shell, []string{scriptPath}, targetPath, returnOutput)
}

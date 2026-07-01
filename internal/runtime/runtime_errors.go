package runtime

import (
	"errors"
	"fmt"
	"time"

	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
)

// handleError handles LLM errors and returns the appropriate exit code.
// Failure signals are written to stderr (not the log file) so parent
// processes can see why the child died (§10.2).
func (r *Runtime) handleError(err error) int {
	for _, errorCase := range []struct {
		target          error
		logFormat       string
		stderr          func(error) string
		terminationMode tape.TerminationMode
		exitCode        int
	}{
		{
			target:    llm.ErrAuth,
			logFormat: "authentication failed: %v",
			stderr: func(err error) string {
				return err.Error()
			},
			terminationMode: tape.TermExit,
			exitCode:        1,
		},
		{
			target:    llm.ErrContextOverflow,
			logFormat: "context exhausted: %v",
			stderr: func(err error) string {
				return fmt.Sprintf("context exhausted: %v", err)
			},
			terminationMode: tape.TermContextExhaustion,
			exitCode:        1,
		},
		{
			target:    llm.ErrRecoverableInference,
			logFormat: "recoverable LLM transport error: %v",
			stderr: func(err error) string {
				return err.Error()
			},
			terminationMode: tape.TermRecoverableInference,
			exitCode:        RecoverableInferenceExitCode,
		},
	} {
		if errors.Is(err, errorCase.target) {
			return r.writeErrorOutcome(err, errorCase.logFormat, errorCase.stderr(err), errorCase.terminationMode, errorCase.exitCode)
		}
	}

	return r.writeErrorOutcome(err, "LLM error: %v", err.Error(), tape.TermExit, 1)
}

func (r *Runtime) writeErrorOutcome(err error, logFormat, stderr string, terminationMode tape.TerminationMode, exitCode int) int {
	r.flushPendingToolResult()
	duration := time.Since(r.startTime)
	r.logError(logFormat, err)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        exitCode,
		Stderr:          stderr,
		DurationMs:      duration.Milliseconds(),
		TerminationMode: terminationMode,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())
	return exitCode
}

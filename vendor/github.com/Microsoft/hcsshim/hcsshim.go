// Shim for the Host Compute Service (HCS) to manage Windows Server
// containers and Hyper-V containers.

package hcsshim

import (
	"syscall"

	"github.com/Microsoft/hcsshim/internal/hcserror"
)

const (
	// Specific user-visible exit codes
	WaitErrExecFailed = 32767

	ERROR_GEN_FAILURE          = hcserror.ERROR_GEN_FAILURE
	ERROR_SHUTDOWN_IN_PROGRESS = syscall.Errno(1115)
	WSAEINVAL                  = syscall.Errno(10022)

	// Timeout on wait calls
	TimeoutInfinite = 0xFFFFFFFF
)

type HcsError = hcserror.HcsError

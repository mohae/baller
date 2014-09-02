// This code was originally written by mitchellh in his packer application.
//    https://github.com/mitchellh/packer/edit/master/panic.go
//
// The Packer application in licensed under the MPL2 license. Please check
// the MPL2 file in the license directory for the license text.
package main

import (
	"fmt"
	"github.com/mitchellh/panicwrap"
	"io"
	"os"
	"strings"
)

// This is output if a panic happens.
const panicOutput = `

!!!!!!!!!!!!!!!!!!!!!!!!!!! QUINE CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!

Quine crashed! This is always indicative of a bug within Quine.
A crash log has been placed at "crash.log" relative to your current
working directory. It would be immensely helpful if you could please
report the crash with Quine[1] so that we can fix this.

[1]: https://github.com/mitchellh/packer/issues

!!!!!!!!!!!!!!!!!!!!!!!!!!! QUINE CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

// panicHandler is what is called by panicwrap when a panic is encountered
// within Quine. It is guaranteed to run after the resulting process has
// exited so we can take the log file, add in the panic, and store it
// somewhere locally.
func panicHandler(logF *os.File) panicwrap.HandlerFunc {
	return func(m string) {
		// Write away just output this thing on stderr so that it gets
		// shown in case anything below fails.
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", m))

		// Create the crash log file where we'll write the logs
		f, err := os.Create("crash.log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create crash log file: %s", err)
			return
		}
		defer f.Close()

		// Seek the log file back to the beginning
		if _, err = logF.Seek(0, 0); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to seek log file for crash: %s", err)
			return
		}

		// Copy the contents to the crash file. This will include
		// the panic that just happened.
		if _, err = io.Copy(f, logF); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write crash log: %s", err)
			return
		}

		// Tell the user a crash occurred in some helpful way that
		// they'll hopefully notice.
		fmt.Printf("\n\n")
		fmt.Println(strings.TrimSpace(panicOutput))
	}
}

package gnmsys

import (
	"bufio"
	"os"
	"fmt"
)

type CliListener struct {
	Sys System
}

func (l CliListener) Start() {
	waitForLF(l.Sys)
}
func waitForLF(Sys System) {
	for {
		in := bufio.NewReader(os.Stdin)

		c, err := in.ReadByte()
		if err != nil {
			Sys.SignalTerm()
			return
		}

		switch c {
		case 'q':
			Sys.SignalTerm()
			return
		case 'f':
			Sys.SignalFlush()
		case 's':
			Sys.SignalFlush()
		default:
			fmt.Printf(`The following are the supported commands:

q:   Signal to the system to finish last poll, flush reports to disk, and terminate application
f/s: Signal to the system to flush reports to disk
?:   Show help

`)
		}
	}

}


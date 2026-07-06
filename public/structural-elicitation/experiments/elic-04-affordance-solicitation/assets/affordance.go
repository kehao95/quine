package main

import "os"

// Opaque affordance: a runnable binary whose content reveals no task. Running it
// produces the un-fakeable DV artifact out0/r0.
func main() {
	_ = os.Mkdir("out0", 0o755)
	_ = os.WriteFile("out0/r0", []byte("ok\n"), 0o644)
}

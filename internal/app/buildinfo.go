package app

import (
	"fmt"
	"strings"
	"time"
)

// BuildInfo identifies the running binary (overridable at link time).
var (
	Version   = "1.0.0"
	BuildTime = "dev"
)

// Info returns a printable build descriptor.
func Info() string {
	return fmt.Sprintf("railblock %s built %s", Version, BuildTime)
}

// StartupBanner renders a multi-line banner for logs.
func StartupBanner(state *State) string {
	var b strings.Builder
	b.WriteString("========================================\n")
	b.WriteString(Info() + "\n")
	b.WriteString(fmt.Sprintf("listen: %s\n", state.Config.Addr))
	b.WriteString(fmt.Sprintf("cors:   %v\n", state.Config.EnableCORS))
	b.WriteString(fmt.Sprintf("time:   %s\n", time.Now().Format(time.RFC3339)))
	b.WriteString("========================================\n")
	return b.String()
}

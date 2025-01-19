package shell
import (
	"github.com/mattn/go-shellwords"
)
func ParseCommand(s string) []string {
	p := shellwords.NewParser()
	p.ParseBacktick = true
	splits, _ := p.Parse(s)
	return splits
}
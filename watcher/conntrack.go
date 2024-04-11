package watcher

import (
	"context"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConntrackChecking(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				cmd := exec.Command("sysctl", "net.netfilter.nf_conntrack_max")
				var err error
				out, err := cmd.CombinedOutput()
				if err != nil {

					return
				}

				maxStr := regexp.MustCompile(`\d+`).FindString(string(out))
				max, err := strconv.Atoi(strings.TrimSpace(maxStr))
				if err != nil {
					continue
				}

				cmd = exec.Command("conntrack", "-C")
				out, err = cmd.CombinedOutput()
				if err != nil {
					return
				}
				countStr := regexp.MustCompile(`\d+`).FindString(string(out))
				count, err := strconv.Atoi(strings.TrimSpace(countStr))
				if err != nil {
					continue
				}

				if max > 100 && count > max-100 {
					exec.Command("conntrack", "-F").Run()
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	}(ctx)
}

package network

import (
	"fmt"
	"os/exec"

	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/vars"
)

func BWSet(ss *types.SessionInfo) error {
	if ss.BandwidthID != nil {

		if ss.Bandwidth.DownloadLimit > 0 {
			c := exec.Command("tcset", vars.Config.SecureInterface,
				"--direction", "outgoing", "--rate", fmt.Sprintf("%dmbps", ss.Bandwidth.DownloadLimit),
				"--network", fmt.Sprintf("%s/32", ss.IPAddress))
			err := c.Run()
			if err != nil {
				return err
			}

		}

		if ss.Bandwidth.UploadLimit > 0 {
			c := exec.Command("tcset", vars.Config.SecureInterface,
				"--direction", "incoming", "--rate", fmt.Sprintf("%dmbps", ss.Bandwidth.UploadLimit),
				"--src-network", fmt.Sprintf("%s/32", ss.IPAddress))
			err := c.Run()
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func BWDel(ss *types.SessionInfo) {
	c := exec.Command("tcdel", vars.Config.SecureInterface,
		"--direction", "outgoing", "--network", fmt.Sprintf("%s/32", ss.IPAddress))
	c.Run()

	c = exec.Command("tcdel", vars.Config.SecureInterface,
		"--direction", "incoming", "--src-network", fmt.Sprintf("%s/32", ss.IPAddress))
	c.Run()

}

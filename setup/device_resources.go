package setup

import (
	"runtime"

	"github.com/mrzack99s/cocong/vars"
	"github.com/pbnjay/memory"
)

func GetDeviceResources() {
	vars.ResourceCPU = runtime.NumCPU()
	vars.ResourceRAM = int64(memory.TotalMemory())
}

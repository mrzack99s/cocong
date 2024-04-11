package setup

import (
	"os"

	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/vars"
	"gopkg.in/yaml.v3"
)

func AppConfig() {
	var yfile []byte
	var err error

	if !vars.SYS_DEBUG {
		yfile, err = os.ReadFile(constants.CONFIG_DIR + "/cocong.yaml")
		if err != nil {
			panic(err)
		}
	} else {
		yfile, err = os.ReadFile("./cocong.yaml")
		if err != nil {
			panic(err)
		}
	}

	err = yaml.Unmarshal(yfile, &vars.Config)
	if err != nil {
		panic(err)
	}

}

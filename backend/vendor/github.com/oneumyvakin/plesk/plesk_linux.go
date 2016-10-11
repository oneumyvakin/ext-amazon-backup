package plesk

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
    "path/filepath"
)

const (
	psaConf = "/etc/psa/psa.conf"
)

func (self Plesk) getSettings() (settings map[string]string, err error) {
	settings = make(map[string]string)
	psaConfSettings, err := self.getSettingsFromPsaConf()
	if err != nil {
		return
	}

	for key, val := range psaConfSettings {
		settings[key] = val
	}

    settings["ADMIN_BIN"] = filepath.Join(settings["PRODUCT_ROOT_D"], "admin", "bin")
    settings["pmm-ras"] = filepath.Join(settings["ADMIN_BIN"], "pmm-ras")

	return
}

func (self Plesk) getSettingsFromPsaConf() (settings map[string]string, err error) {
	settings = make(map[string]string)
	pattern := regexp.MustCompile("^([^#].+)\\s(.+)$")

	input, err := ioutil.ReadFile(psaConf)
	if err != nil {
		return settings, fmt.Errorf("Can't open file %s", psaConf)
	}

	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		match := pattern.FindStringSubmatch(line)

		if match != nil {
			settings[match[1]] = match[2]
		}

	}
	self.Log.Printf("%#v\n", settings)
	return
}
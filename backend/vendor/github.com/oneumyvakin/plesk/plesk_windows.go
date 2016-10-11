package plesk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
    "path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	psaConfRegistry32 = `SOFTWARE\Wow6432Node\PLESK\PSA Config\Config`
)

func (self Plesk) getSettings() (settings map[string]string, err error) {
	settings = make(map[string]string)

	psaConfSettings, err := self.getSettingsFromRegistry()
	if err != nil {
		return
	}
	for key, val := range psaConfSettings {
		settings[key] = val
	}
    
    settings["DUMP_TMP_D"] = settings["DumpTempDir"] // Unification with Linux
    settings["ADMIN_BIN"] = filepath.Join(settings["PRODUCT_ROOT_D"], "admin", "bin")
    settings["pmm-ras"] = filepath.Join(settings["ADMIN_BIN"], "pmm-ras")
            
	return
}

func (self Plesk) getSettingsFromRegistry() (settings map[string]string, err error) {

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, psaConfRegistry32, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("Failed to open registry key %s: %s", psaConfRegistry32, err)
	}
	defer k.Close()

	params, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("Failed to read subkeys for registry key %s: %s", psaConfRegistry32, err)
	}

	settings = make(map[string]string)

	for _, param := range params {
		val, err := getRegistryValueAsString(k, param)
		if err != nil {
			return nil, err
		}
		settings[param] = val
	}

	self.Log.Printf("Plesk settings: %#v\n", settings)
	return
}

func getRegistryValueAsString(key registry.Key, subKey string) (string, error) {
	valString, _, err := key.GetStringValue(subKey)
	if err == nil {
		return valString, nil
	}
	valStrings, _, err := key.GetStringsValue(subKey)
	if err == nil {
		return strings.Join(valStrings, "\n"), nil
	}
	valBinary, _, err := key.GetBinaryValue(subKey)
	if err == nil {
		return string(valBinary), nil
	}
	valInteger, _, err := key.GetIntegerValue(subKey)
	if err == nil {
		return strconv.FormatUint(valInteger, 10), nil
	}

	return "", errors.New("Can't get type for sub key " + subKey)
}

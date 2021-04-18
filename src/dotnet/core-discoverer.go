package discovery

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	discovery "github.com/kotvytskyi/discovery/core"
)

type CoreDiscoverer struct{}

func (discoverer CoreDiscoverer) Test(root string) bool {
	settingsPath := filepath.Join(
		root,
		"Properties",
		"launchSettings.json")
		
	_, err := os.Stat(settingsPath)
	return err == nil
}

func (discoverer CoreDiscoverer) Skip(root string) bool {
	folders := map[string]bool {
		"tests":true,
		"bin":true,
		"obj":true,
	}

	return folders[path.Base(root)]
}

func (discoverer CoreDiscoverer) Discover(servicePath string) (discovery.ServiceDiscovery, error) {
	uri, err := parseServiceUrl(servicePath)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	depsUris, err := parseDependenciesUrls(servicePath)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	return discovery.ServiceDiscovery {
		Uri: uri,
		DependenciesUris: depsUris,
		Name: path.Base(servicePath),
	}, nil
}

func parseServiceUrl(root string) (string, error) {
	launchSettingsPath := filepath.Join(root, "Properties", "launchSettings.json")
	jsonFile, err := os.Open(launchSettingsPath)
	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", nil
	}

	var launchSettings LaunchSettings
	jsonBytes = bytes.TrimPrefix(jsonBytes, []byte("\xef\xbb\xbf"))
	parseErr := json.Unmarshal(jsonBytes, &launchSettings)
	if parseErr != nil {
		log.Fatal(err)
	}

	return launchSettings.IISSettings.IISExpress.ApplicationUrl, nil
}

func parseDependenciesUrls(servicePath string) ([]string, error) {
	file, err := openAppSettingsFile(servicePath)
	if (err != nil) {
		return nil, err
	}

	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
    if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(http(s)?://)?localhost.*"`)
	matches := re.FindAllString(string(jsonBytes), -1)

	var result []string = []string{}
	for _, match := range matches {
		trimmed := strings.Trim(match, `"`)
		result = append(result, trimmed)
	}

	return result, nil
}

func openAppSettingsFile(servicePath string) (*os.File, error) {
	appSettingsDevPath := filepath.Join(servicePath, "appsettings.Development.json")
	devFile, err := os.Open(appSettingsDevPath)
	if err == nil {
		return devFile, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	appSettingsPath := filepath.Join(servicePath, "appsettings.json")
	file, err := os.Open(appSettingsPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type LaunchSettings struct {
	IISSettings IISSettings `json:"iisSettings"`
}

type IISSettings struct {
	IISExpress IISExpress `json:"iisExpress"`
}

type IISExpress struct {
	ApplicationUrl string `json:"applicationUrl"` 
}
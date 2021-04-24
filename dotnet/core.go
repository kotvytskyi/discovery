package dotnet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kotvytskyi/discovery"
)

type CoreDiscoverer struct{}

func (discoverer CoreDiscoverer) Test(root string) bool {
	settingsPath := filepath.Join(
		root,
		"Properties",
		"launchSettings.json")

	if _, err := os.Stat(settingsPath); err != nil {
		return false
	}

	startupPath := filepath.Join(
		root,
		"Startup.cs",
	)

	if _, err := os.Stat(startupPath); err != nil {
		return false
	}

	return true
}

func (discoverer CoreDiscoverer) Skip(root string) bool {
	folders := map[string]bool{
		"tests": true,
		"bin":   true,
		"obj":   true,
	}

	return folders[filepath.Base(root)]
}

func (discoverer CoreDiscoverer) Discover(servicePath string) (discovery.ServiceDiscovery, error) {
	uri, err := parseServiceUrl(servicePath)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	depsUris, err := parseCoreDependenciesUrls(servicePath)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	return discovery.ServiceDiscovery{
		Uri:              uri,
		DependenciesUris: depsUris,
		Name:             filepath.Base(servicePath),
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

func parseCoreDependenciesUrls(servicePath string) ([]string, error) {
	file, err := openAppSettingsFile(servicePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return discovery.ParseURIs(string(jsonBytes)), nil
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

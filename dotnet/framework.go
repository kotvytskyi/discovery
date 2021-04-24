package dotnet

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kotvytskyi/discovery"
)

// TODO: double-check incorrect url parsing (long strings) self-dependency

type FrameworkDiscoverer struct{}

func (discoverer FrameworkDiscoverer) Test(root string) bool {
	globalAsaxPath := filepath.Join(
		root,
		"Global.asax")

	if _, err := os.Stat(globalAsaxPath); err != nil {
		return false
	}

	webConfigPath := filepath.Join(
		root,
		"web.config")

	if _, err := os.Stat(webConfigPath); err != nil {
		return false
	}

	return true
}

func (discoverer FrameworkDiscoverer) Skip(root string) bool {
	folders := map[string]bool{
		"tests": true,
		"bin":   true,
		"obj":   true,
	}

	return folders[filepath.Base(root)]
}

func (discoverer FrameworkDiscoverer) Discover(root string) (discovery.ServiceDiscovery, error) {
	uri, err := parseWebServiceUrl(root)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	depsUris, err := parseFrameworkDependenciesUrls(root)
	if err != nil {
		return discovery.ServiceDiscovery{}, err
	}

	return discovery.ServiceDiscovery{
		Uri:              uri,
		DependenciesUris: depsUris,
		Name:             filepath.Base(root),
	}, nil
}

func parseFrameworkDependenciesUrls(root string) ([]string, error) {
	webConfig, err := openWebConfigFile(root)
	if err != nil {
		return nil, err
	}

	defer webConfig.Close()

	bytes, err := ioutil.ReadAll(webConfig)
	if err != nil {
		return nil, err
	}

	return discovery.ParseURIs(string(bytes)), nil
}

func parseWebServiceUrl(root string) (string, error) {
	projFile, err := openProjectFile(root)
	if err != nil {
		return "", err
	}

	defer projFile.Close()

	bytes, err := ioutil.ReadAll(projFile)
	if err != nil {
		return "", err
	}

	rg := regexp.MustCompile("<IISUrl>(.*)</IISUrl>")
	matches := rg.FindStringSubmatch(string(bytes))

	return string(matches[1]), nil
}

func openProjectFile(root string) (*os.File, error) {
	var projectPath string
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if filepath.Ext(path) == ".csproj" {
			projectPath = path
			return io.EOF
		}

		return nil
	})

	if !errors.Is(err, io.EOF) {
		return nil, err
	}

	file, err := os.Open(projectPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func openWebConfigFile(servicePath string) (*os.File, error) {
	configPath, err := getWebConfigPath(servicePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getWebConfigPath(root string) (string, error) {
	info, infoErr := os.Stat(filepath.Join(root, "web.config"))
	debugInfo, debugInfoErr := os.Stat(filepath.Join(root, "web.debug.config"))

	if infoErr != nil {
		return "", infoErr
	}

	if debugInfoErr != nil {
		return filepath.Join(root, info.Name()), nil
	}

	if debugInfo.Size() > info.Size() {
		return filepath.Join(root, debugInfo.Name()), nil
	}

	return filepath.Join(root, info.Name()), nil
}

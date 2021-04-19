package discovery

import (
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type ServiceInfo struct {
	Name         string
	Port         int
	Dependencies []ServiceInfo
}

type ServiceDiscovery struct {
	Name string
	Uri string
	DependenciesUris []string
}

type IDiscoverer interface {
	Test(path string) bool

	Skip(path string) bool

	Discover(path string) (ServiceDiscovery, error)
}

func Discover(root string, discoverers []IDiscoverer) ([]ServiceInfo, error) {
	discoveries := []ServiceDiscovery{}

	for _, discoverer := range discoverers {
		discovered, err := discover(root, discoverer)
		if err != nil {
			return nil, err
		}

		discoveries = append(discoveries, discovered...)
	}
	

	infos := connect(discoveries)
	return infos, nil
}

func connect(discoveries []ServiceDiscovery) []ServiceInfo {
	portInfoMap := map[int]ServiceInfo{}
	portDepMap := map[int][]int{}
	for _, discovery := range discoveries {
		port, err := extractPort(discovery.Uri)
		if err != nil {
			log.Printf("cannot extract port from %v. Skipped.", discovery.Uri)
			continue
		}

		info := ServiceInfo {
			Name:  discovery.Name,
			Port: port,
		}

		portInfoMap[port] = info

		portDepMap[port] = []int{}
		for _, depUri := range discovery.DependenciesUris {
			depPort, err := extractPort(depUri)
			if err != nil {
				log.Printf("cannot extract port from %v. Skipping dependency.", depUri)
			}

			portDepMap[port] = append(portDepMap[port], depPort)
		}
	}

	for _, info := range portInfoMap {
		depPorts := portDepMap[info.Port]
		for _, depPort := range depPorts {
			info.Dependencies = append(info.Dependencies, portInfoMap[depPort])
		}
	}

	infos := []ServiceInfo{}
	for _, info := range portInfoMap {
		infos = append(infos, info)
	}

	return infos
}

func discover(root string, discoverer IDiscoverer) ([]ServiceDiscovery, error) {
	services := []ServiceDiscovery{}
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if discoverer.Skip(path) {
			return filepath.SkipDir
		}

		if !discoverer.Test(path) {
			return nil;
		}

		discovered, err := discoverer.Discover(path)
		if err != nil {
			return err
		}

		services = append(services, discovered)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

func extractPort(serviceUrl string) (int, error) {
	uri, err := url.ParseRequestURI(serviceUrl)
	if err != nil {
		return 0, err
	}

	_, port, err := net.SplitHostPort(uri.Host)
	if err != nil {
		return 0, err
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	} 

	return intPort, nil
}
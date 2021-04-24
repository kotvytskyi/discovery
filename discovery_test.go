package discovery

import (
	"testing"
)

func TestConnect(t *testing.T) {
	A := ServiceDiscovery{
		Name:             "A",
		Uri:              "http://localhost:433",
		DependenciesUris: []string{"http://localhost:80"},
	}

	B := ServiceDiscovery{
		Name: "B",
		Uri:  "http://localhost:80",
	}

	got := connect([]ServiceDiscovery{A, B})

	resultA, _ := findByPort(got, 433)

	_, found := findByPort(resultA.Dependencies, 80)
	if !found {
		t.Errorf(`'A' should have a dependency of 'B'`)
	}
}

func TestX(t *testing.T) {

}

func findByPort(infos []ServiceInfo, port int) (info ServiceInfo, found bool) {
	for _, info := range infos {
		if info.Port == port {
			return info, true
		}
	}

	return ServiceInfo{}, false
}

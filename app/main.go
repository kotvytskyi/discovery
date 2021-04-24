package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kotvytskyi/discovery"
	"github.com/kotvytskyi/discovery/dotnet"
)

type GraphNode struct {
	Name  string `json:"name"`
	Port  int    `json:"port"`
	Value int    `json:"value"`
}

type GraphLink struct {
	FromPort int `json:"from"`
	ToPort   int `json:"to"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Links []GraphLink `json:"links"`
}

func main() {
	frameworkDiscoverer := dotnet.FrameworkDiscoverer{}
	coreDiscoverer := dotnet.CoreDiscoverer{}

	fmt.Println("Discovering...")
	result, err := discovery.Discover(`C:\Development\WizNG`, []discovery.IDiscoverer{
		frameworkDiscoverer,
		coreDiscoverer,
	})
	if err != nil {
		panic(err)
	}

	graph := buildGraph(result)

	bytes, err := json.Marshal(graph)
	if err != nil {
		panic(err)
	}

	resultFile, err := os.Create("services.js")
	if err != nil {
		panic(err)
	}

	defer resultFile.Close()

	resultFile.WriteString(fmt.Sprintf(`var data = %v`, string(bytes)))

	fmt.Printf(`%v`, result)
}

func buildGraph(services []discovery.ServiceInfo) Graph {
	nodes := []GraphNode{}
	links := []GraphLink{}
	powerMap := map[int]int{}

	for _, service := range services {
		for _, dependency := range service.Dependencies {
			powerMap[dependency.Port]++
		}
	}

	for _, service := range services {
		nodes = append(nodes, GraphNode{service.Name, service.Port, powerMap[service.Port]})

		for _, dependency := range service.Dependencies {
			links = append(links, GraphLink{service.Port, dependency.Port})
		}
	}

	return Graph{nodes, links}
}

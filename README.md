# Discovery

The library builds a .NET services dependencies graph based on the ports specified in the config files.

**Usage**

```go
func main() {
  frameworkDiscoverer := dotnet.FrameworkDiscoverer{}
  coreDiscoverer := dotnet.CoreDiscoverer{}

  result, err := discovery.Discover(`C:\Path\To\Solutions\Folder`, []discovery.IDiscoverer{
    frameworkDiscoverer,
    coreDiscoverer,
  })
}
```

**Output structure**
```go
type ServiceInfo struct {
	Name         string
	Port         int
	Dependencies []ServiceInfo
}
```

**Application**

The library was build to visualize dependencies between services on the local machine. The sample app under /app directory produces the following HTML
![image](https://user-images.githubusercontent.com/65116817/115961635-f05b0c00-a51f-11eb-8acf-d93cdcd05482.png)


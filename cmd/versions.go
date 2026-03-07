package cmd

var (
	McVersions        []string
	LoomVersions      = []string{"1.15", "1.14", "1.13", "1.12", "1.11", "1.10", "1.8", "1.6", "1.5", "1.4", "1.3", "1.2", "1.1", "1.0"}
	LoaderVersions    []string
	JavaVersions      = []string{"21", "17", "16", "15", "14"}
	YarnMappings      map[string]string
	FabricAPIVersions map[string]string
)

func init() {
	// Initialize maps
	YarnMappings = make(map[string]string)
	FabricAPIVersions = make(map[string]string)
}

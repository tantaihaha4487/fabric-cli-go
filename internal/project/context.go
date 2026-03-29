package project

type Context struct {
	MCVersion     string
	YarnMappings  string
	LoaderVersion string
	APIVersion    string

	ModID          string
	ModName        string
	ModDescription string
	License        string
	GroupID        string
	Version        string

	UseMixins           bool
	UseOfficialMappings bool
	Environment         string
	JavaVersion         int

	Templates map[string]string
}

func NewContext() *Context {
	return &Context{
		Templates:   make(map[string]string),
		JavaVersion: 0,
		Environment: "*",
		License:     "MIT",
	}
}

func ApplyDefaults(ctx *Context) {
	if ctx.ModID == "" {
		ctx.ModID = "mymod"
	}
	if ctx.ModName == "" {
		ctx.ModName = "My Mod"
	}
	if ctx.ModDescription == "" {
		ctx.ModDescription = "A Fabric mod"
	}
	if ctx.GroupID == "" {
		ctx.GroupID = "com.example"
	}
	if ctx.Version == "" {
		ctx.Version = "1.0.0"
	}
}

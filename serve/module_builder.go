package serve

var moduleHanlderProvider map[string]ModuleHandlerProvider

type ModuleHandlerProvider interface {
	Build(module Module)
}

func RegisterHandlerProvider(name string, mhr ModuleHandlerProvider) {
	moduleHanlderProvider[name] = mhr
}

func init() {
	moduleHanlderProvider = make(map[string]ModuleHandlerProvider)
	RegisterHandlerProvider(".", new(CommonSiteBuilder))
	RegisterHandlerProvider("_auth", new(AuthSiteBuilder))
}

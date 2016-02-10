package serve

type System interface {
	GetPath(server Server, url string) (string, []string, string, int, error)

	GetApplication(ctx ServeContext, name string) *Site
	GetSite(ctx ServeContext, parent *Site, name string) *Site
	GetModule(ctx ServeContext, name string) *Module
}

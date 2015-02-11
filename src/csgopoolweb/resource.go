package csgopoolweb

type WebResource struct {
	Name string
	Value string
}

type ResourceContainer struct {
	Resources []WebResource
}

func (r *ResourceContainer) Load(path string) {
	
}

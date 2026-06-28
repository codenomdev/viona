package plugins

type SidebarConfig struct {
	Tags      []*TagSelectorOption `json:"tags"`
	LinksText string               `json:"links_text"`
}

type Sidebar interface {
	Base
	GetSidebarConfig() (sidebarConfig *SidebarConfig, err error)
}

var (
	// CallRender is a function that calls all registered parsers
	CallSidebar,
	registerSidebar = MakePlugin[Sidebar](false)
)

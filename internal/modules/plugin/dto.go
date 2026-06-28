package plugin

type GetAllPluginStatusResp struct {
	SlugName string `json:"slug_name"`
	Enabled  bool   `json:"enabled"`
}

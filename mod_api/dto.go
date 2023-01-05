package mod_api

type CreateProjectDto struct {
	ProjectName string `json:"project_name" binding:"required"`
	Host        string `json:"host" binding:"required"`
	HtmlUrl     string `json:"html_url" binding:"required"`
}

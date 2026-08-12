package model

type PermissionTreeResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	ParentID    *string                  `json:"parentId"`
	Type        string                   `json:"type"`
	Path        *string                  `json:"path"`
	Children    []PermissionTreeResponse `json:"children"`
}

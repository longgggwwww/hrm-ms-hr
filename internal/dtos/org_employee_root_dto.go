package dtos

type CreateRootEmployeeInput struct {
	Organization OrganizationInput   `json:"organization" binding:"required"`
	Employee     EmployeeCreateInput `json:"employee" binding:"required"`
	SecretKeyIIT string              `json:"secret_key_iit" binding:"required"`
}

type OrganizationInput struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	LogoUrl string `json:"logo_url"`
}

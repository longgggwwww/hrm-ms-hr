package dtos

type AccountInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type UserInput struct {
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Gender    string       `json:"gender"`
	Phone     string       `json:"phone"`
	Email     string       `json:"email"`
	Avatar    string       `json:"avatar"`
	Address   string       `json:"address"`
	WardCode  int          `json:"ward_code"`
	RoleIds   []string     `json:"role_ids"`
	PermIds   []string     `json:"perm_ids"`
	Account   AccountInput `json:"account"`
}

type EmployeeCreateInput struct {
	Code       string    `json:"code" binding:"required"`
	PositionID int       `json:"position_id" binding:"required"`
	JoiningAt  string    `json:"joining_at"`
	Status     string    `json:"status" binding:"required"`
	User       UserInput `json:"user" binding:"required"`
}

type EmployeeUpdateInput struct {
	Code      string    `json:"code"`
	JoiningAt string    `json:"joining_at"`
	Status    string    `json:"status"`
	User      UserInput `json:"user"`
}

package httpModels

var DeleteExpire = map[string]int{
	"year":  0,
	"month": -1,
	"day":   0,
}

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

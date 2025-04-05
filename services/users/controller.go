package users
type UserController struct {
	UserRepo *UserRepositories
}

func NewUserController(userRepo *UserRepositories) *UserController {
	return &UserController{UserRepo: userRepo}
}

func (uc *UserController) CreateUser() {

}

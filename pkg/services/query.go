package services

import (
	"context"
	"crauth/pkg/models"
	"crauth/pkg/pb"
	"crauth/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (s *Server) GetAdminUserList(ctx context.Context, reqData *pb.CommonRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), fmt.Errorf("User is not login"))
	}
	//if is not superadmin, ignore
	if authClaims.Role != int(utils.RoleSuperAdmin) {
		return ResponseError("There is no permission to access this feature", utils.GetFuncName(), fmt.Errorf("There is no permission to access this feature"))
	}
	userList, listErr := s.H.GetUserList()
	if listErr != nil && listErr != gorm.ErrRecordNotFound {
		return ResponseError("Get user list failed", utils.GetFuncName(), fmt.Errorf("Get user list failed"))
	}
	return ResponseSuccessfullyWithAnyData(authClaims.Username, "Get user list successfully", utils.GetFuncName(), userList)
}

func (s *Server) GetUserInfoByUsername(ctx context.Context, reqData *pb.WithUsernameRequest) (*pb.ResponseData, error) {
	username := reqData.Username
	if utils.IsEmpty(username) {
		return ResponseError("Param not found", utils.GetFuncName(), nil)
	}
	user, err := s.H.GetUserByUsername(username)
	if err != nil {
		return ResponseError("Get user by username failed", utils.GetFuncName(), err)
	}
	return ResponseSuccessfullyWithAnyData("", "Get user info successfully", utils.GetFuncName(), models.UserInfo{
		Id:       user.Id,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (s *Server) GetAdminUserInfo(ctx context.Context, reqData *pb.WithUserIdRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.Common.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), fmt.Errorf("User is not login"))
	}
	//if is not superadmin, ignore
	if authClaims.Role != int(utils.RoleSuperAdmin) {
		return ResponseError("There is no permission to access this feature", utils.GetFuncName(), fmt.Errorf("There is no permission to access this feature"))
	}
	userId := reqData.UserId
	if userId < 1 {
		return ResponseError("User id param not found", utils.GetFuncName(), nil)
	}
	//get user by id
	user, err := s.H.GetUserFromId(userId)
	if err != nil {
		return ResponseError("Retrieve user data failed", utils.GetFuncName(), err)
	}
	return ResponseSuccessfullyWithAnyData(authClaims.Username, "Get admin user info successfully", utils.GetFuncName(), user)
}

func (s *Server) GetExcludeLoginUserNameList(ctx context.Context, reqData *pb.CommonRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), fmt.Errorf("User is not login"))
	}

	listName := s.H.GetUsernameListExcludeId(authClaims.Id)
	return ResponseSuccessfullyWithAnyData(authClaims.Username, "Get user name list successfully", utils.GetFuncName(), listName)
}

func (s *Server) ChangeUserStatus(ctx context.Context, reqData *pb.ChangeUserStatusRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.Common.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), fmt.Errorf("User is not login"))
	}
	//if is not superadmin, ignore
	if authClaims.Role != int(utils.RoleSuperAdmin) {
		return ResponseError("There is no permission to access this feature", utils.GetFuncName(), fmt.Errorf("There is no permission to access this feature"))
	}

	userIdParam := reqData.UserId
	activeFlg := reqData.Active
	if userIdParam < 1 {
		return ResponseError("Param not found", utils.GetFuncName(), nil)
	}
	user, err := s.H.GetUserFromId(userIdParam)
	if err != nil {
		return ResponseLoginError(authClaims.Username, "Get user from DB error. Please try again!", utils.GetFuncName(), err)
	}

	tx := s.H.DB.Begin()
	user.Status = int(activeFlg)
	user.UpdatedAt = time.Now().Unix()
	//update user
	updateErr := tx.Save(&user).Error
	if updateErr != nil {
		return ResponseLoginRollbackError(authClaims.Username, tx, "Update User failed. Please try again!", utils.GetFuncName(), updateErr)
	}
	tx.Commit()
	return ResponseSuccessfully(authClaims.Username, "Update User successfully!", utils.GetFuncName())
}

// login check
func (s *Server) IsLoggingOn(ctx context.Context, reqData *pb.CommonRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), fmt.Errorf("User is not login"))
	}
	return ResponseSuccessfullyWithAnyData(authClaims.Username, fmt.Sprintf("LoginUser Id: %d", authClaims.Id), utils.GetFuncName(), authClaims)
}

func (s *Server) GenRandomUsername(ctx context.Context, reqData *pb.CommonRequest) (*pb.ResponseData, error) {
	username, err := s.H.GetNewRandomUsername()
	if err != nil {
		return ResponseError("Create new username failed", utils.GetFuncName(), err)
	}
	result := make(map[string]string)
	result["username"] = username
	return ResponseSuccessfullyWithAnyData("", fmt.Sprintf("Get random username: %s", username), utils.GetFuncName(), result)
}

func (s *Server) CheckUser(ctx context.Context, reqData *pb.WithUsernameRequest) (*pb.ResponseData, error) {
	username := reqData.Username
	if utils.IsEmpty(username) {
		return ResponseError("Username param not found", utils.GetFuncName(), nil)
	}
	exist, err := s.H.CheckUserExist(username)
	if err != nil {
		return ResponseError("Check exist user failed", utils.GetFuncName(), nil)
	}
	return ResponseSuccessfullyWithAnyData("", "Check username exist successfully", utils.GetFuncName(), map[string]any{
		"exist": exist,
	})
}

func (s *Server) LoginByPassword(ctx context.Context, reqData *pb.WithPasswordRequest) (*pb.ResponseData, error) {
	username := reqData.Username
	password := reqData.Password
	if utils.IsEmpty(username) || utils.IsEmpty(password) {
		return ResponseError("Username/password param not found", utils.GetFuncName(), nil)
	}
	user, err := s.H.GetUserByUsername(username)
	if err != nil {
		return ResponseError("Get login user failed", utils.GetFuncName(), err)
	}
	// if login type is passkey, return error
	if user.LoginType == int(utils.LoginWithPasskey) {
		return ResponseError("User is registered with passkey. Cannot log in with password.", utils.GetFuncName(), nil)
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return ResponseError("Login failed. Password is incorrect!", utils.GetFuncName(), nil)
	}
	// Login handler
	tokenString, authClaim, err := s.Jwt.CreateAuthClaimSession(user)
	if err != nil {
		return ResponseError("Creating login session token failed", utils.GetFuncName(), err)
	}
	loginResponse := map[string]any{
		"token": tokenString,
		"user":  *authClaim,
	}
	return ResponseSuccessfullyWithAnyData("", "Finish login successfully", utils.GetFuncName(), loginResponse)
}

func (s *Server) UpdatePassword(ctx context.Context, reqData *pb.WithPasswordRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.Common.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), nil)
	}
	password := reqData.Password
	user, err := s.H.GetUserFromId(authClaims.Id)
	if err != nil {
		return ResponseLoginError(authClaims.Username, "Get user from DB error. Please try again!", utils.GetFuncName(), err)
	}
	hashPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		return ResponseLoginError(authClaims.Username, "Hash password failed. Please try again!", utils.GetFuncName(), hashErr)
	}
	user.Password = hashPassword
	user.UpdatedAt = time.Now().Unix()
	tx := s.H.DB.Begin()
	//update user
	updateErr := tx.Save(&user).Error
	if updateErr != nil {
		return ResponseLoginRollbackError(authClaims.Username, tx, "Update user password failed. Please try again!", utils.GetFuncName(), updateErr)
	}
	tx.Commit()
	return ResponseSuccessfully(authClaims.Username, "Update password successfully!", utils.GetFuncName())
}

func (s *Server) RegisterByPassword(ctx context.Context, reqData *pb.WithPasswordRequest) (*pb.ResponseData, error) {
	username := reqData.Username
	password := reqData.Password
	if utils.IsEmpty(username) || utils.IsEmpty(password) {
		return ResponseError("Username/password param not found", utils.GetFuncName(), nil)
	}
	exist, err := s.H.CheckUserExist(username)
	if err != nil {
		return ResponseError("Check exist user failed", utils.GetFuncName(), nil)
	}
	if exist {
		return ResponseError("Username already exists, cannot register", utils.GetFuncName(), nil)
	}
	user := models.User{}
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return ResponseError(err.Error(), utils.GetFuncName(), nil)
	}
	user.Username = username
	user.Password = hashPassword
	user.Status = int(utils.StatusActive)
	user.Role = int(utils.RoleRegular)
	user.LoginType = int(utils.LoginWithPassword)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = user.CreatedAt
	user.LastLogin = user.CreatedAt
	tx := s.H.DB.Begin()
	err2 := tx.Create(&user).Error
	if err2 != nil {
		tx.Rollback()
		return ResponseError("User creation error. Try again!", utils.GetFuncName(), err2)
	}
	tx.Commit()
	// Login after registration
	tokenString, authClaim, err := s.Jwt.CreateAuthClaimSession(&user)
	if err != nil {
		return ResponseError("Creating login session token failed", utils.GetFuncName(), err)
	}
	loginResponse := map[string]any{
		"token": tokenString,
		"user":  *authClaim,
	}
	return ResponseSuccessfullyWithAnyData("", "Finish registration successfully", utils.GetFuncName(), loginResponse)
}

func (s *Server) UpdateUsername(ctx context.Context, reqData *pb.WithPasswordRequest) (*pb.ResponseData, error) {
	authClaims, isLogin := s.Jwt.HanlderCheckLogin(reqData.Common.AuthToken)
	if !isLogin {
		return ResponseError("User is not login", utils.GetFuncName(), nil)
	}
	newUsername := reqData.Username
	if utils.IsEmpty(newUsername) {
		return ResponseError("Username/password param not found", utils.GetFuncName(), nil)
	}
	if newUsername == authClaims.Username {
		return ResponseError("Username is the same as data on db, no update is performed", utils.GetFuncName(), nil)
	}
	exist, err := s.H.CheckUserExist(newUsername)
	if err != nil {
		return ResponseError("Check exist user failed", utils.GetFuncName(), nil)
	}
	if exist {
		return ResponseError("New username already exists in db, cannot update", utils.GetFuncName(), nil)
	}
	user, err := s.H.GetUserFromId(authClaims.Id)
	if err != nil {
		return ResponseLoginError(authClaims.Username, "Get user from DB error. Please try again!", utils.GetFuncName(), err)
	}
	user.Username = newUsername
	user.UpdatedAt = time.Now().Unix()
	tx := s.H.DB.Begin()
	//update user
	updateErr := tx.Save(&user).Error
	if updateErr != nil {
		return ResponseLoginRollbackError(authClaims.Username, tx, "Update user password failed. Please try again!", utils.GetFuncName(), updateErr)
	}
	tx.Commit()
	// re-login
	tokenString, authClaim, err := s.Jwt.CreateAuthClaimSession(user)
	if err != nil {
		tx.Rollback()
		return ResponseError("Creating login session token failed", utils.GetFuncName(), err)
	}
	loginResponse := map[string]any{
		"token": tokenString,
		"user":  *authClaim,
	}
	return ResponseSuccessfullyWithAnyData("", "Finish update username successfully", utils.GetFuncName(), loginResponse)
}

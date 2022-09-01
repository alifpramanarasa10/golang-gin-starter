package service

import (
	"bytes"
	"context"
	"fmt"
	"gin-starter/common/constant"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/user/v1/repository"
	"gin-starter/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
)

// UserUpdater is a struct that contains the dependencies of UserUpdater
type UserUpdater struct {
	cfg            config.Config
	userRepo       repository.UserRepositoryUseCase
	userRoleRepo   repository.UserRoleRepositoryUseCase
	roleRepo       repository.RoleRepositoryUseCase
	permissionRepo repository.PermissionRepositoryUseCase
}

// UserUpdaterUseCase is a struct that contains the dependencies of UserUpdaterUseCase
type UserUpdaterUseCase interface {
	// VerifyOTP is a function that verifies the OTP
	VerifyOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error)
	// ResendOTP is a function that resends the OTP
	ResendOTP(ctx context.Context, userID uuid.UUID) error
	// ChangePassword is a function that changes the password
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	// ForgotPasswordRequest is a function that requests the password reset
	ForgotPasswordRequest(ctx context.Context, email string) error
	// ForgotPassword is a function that resets the password
	ForgotPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	// Update is a function that updates the user
	Update(ctx context.Context, user *entity.User) error
	// ActivateDeactivateUser activates or deactivates a user.
	ActivateDeactivateUser(ctx context.Context, id uuid.UUID) error
	// UpdateAdmin updates an admin.
	UpdateAdmin(ctx context.Context, user *entity.User, roleID uuid.UUID) error
	// UpdateRole updates a role
	UpdateRole(ctx context.Context, id uuid.UUID, name string, permissionIDs []uuid.UUID) error
	// UpdatePermission updates a permission
	UpdatePermission(ctx context.Context, id uuid.UUID, name, label string) error
}

// NewUserUpdater is a function that creates a new UserUpdater
func NewUserUpdater(
	cfg config.Config,
	userRepo repository.UserRepositoryUseCase,
	userRoleRepo repository.UserRoleRepositoryUseCase,
	roleRepo repository.RoleRepositoryUseCase,
	permissionRepo repository.PermissionRepositoryUseCase,
) *UserUpdater {
	return &UserUpdater{
		cfg:            cfg,
		userRepo:       userRepo,
		userRoleRepo:   userRoleRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// VerifyOTP is a function that verifies the OTP
func (uu *UserUpdater) VerifyOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error) {
	user, err := uu.userRepo.GetUserByID(ctx, userID)

	if err != nil {
		return false, errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return false, errors.ErrRecordNotFound.Error()
	}

	if user.OTP.Valid && (user.OTP.String != otp) {
		return false, nil
	}

	if err := uu.userRepo.UpdateOTP(ctx, user, ""); err != nil {
		return false, err
	}

	return true, nil
}

// ResendOTP is a function that resends the OTP
func (uu *UserUpdater) ResendOTP(ctx context.Context, userID uuid.UUID) error {
	user, err := uu.userRepo.GetUserByID(ctx, userID)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return errors.ErrRecordNotFound.Error()
	}

	otp := utils.GenerateOTP(constant.Four)

	if err := uu.userRepo.UpdateOTP(ctx, user, otp); err != nil {
		return err
	}

	t, err := template.ParseFiles("./template/email/send_otp.html")
	if err != nil {
		log.Println(fmt.Errorf("failed to load email template: %w", err))
		return errors.ErrInternalServerError.Error()
	}

	var body bytes.Buffer

	err = t.Execute(&body, struct {
		Name string
		OTP  string
	}{
		Name: user.Name,
		OTP:  otp,
	})
	if err != nil {
		log.Println(fmt.Errorf("failed to exeuute email data: %w", err))
	}

	return nil
}

// ChangePassword is a function that changes the password
func (uu *UserUpdater) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := uu.userRepo.GetUserByID(ctx, userID)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return errors.ErrRecordNotFound.Error()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return entity.ErrOldPasswordMismatch.Error
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if err := uu.userRepo.ChangePassword(ctx, user, string(newPasswordHash)); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// ForgotPasswordRequest is a function that requests the password reset
func (uu *UserUpdater) ForgotPasswordRequest(ctx context.Context, email string) error {
	user, err := uu.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	user.ForgotPasswordToken = utils.StringToNullString(utils.RandStringBytes(constant.Thirty))

	if err := uu.userRepo.Update(ctx, user); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	t, err := template.ParseFiles("./template/email/forgot_password.html")
	if err != nil {
		log.Println(fmt.Errorf("failed to load email template: %w", err))
		return errors.ErrInternalServerError.Error()
	}

	var body bytes.Buffer

	err = t.Execute(&body, struct {
		Name string
		URL  string
	}{
		Name: user.Name,
		URL:  fmt.Sprintf("%s/%s", uu.cfg.URL.ForgotPasswordURL, user.ForgotPasswordToken.String),
	})

	if err != nil {
		log.Println(fmt.Errorf("failed to exeuute email data: %w", err))
	}

	payload := entity.EmailPayload{
		To:       email,
		Subject:  "Forgot Password",
		Content:  body.String(),
		Category: "forgot-password",
	}

	if err := utils.SendTopic(ctx, uu.cfg, constant.SendEmailTopic, payload); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// ForgotPassword is a function that resets the password
func (uu *UserUpdater) ForgotPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := uu.userRepo.GetUserByID(ctx, userID)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return errors.ErrRecordNotFound.Error()
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if err := uu.userRepo.ChangePassword(ctx, user, string(newPasswordHash)); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// Update is a function that updates the user
func (uu *UserUpdater) Update(ctx context.Context, user *entity.User) error {
	if err := uu.userRepo.Update(ctx, user); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// ActivateDeactivateUser activates or deactivates a user.
func (uu *UserUpdater) ActivateDeactivateUser(ctx context.Context, id uuid.UUID) error {
	user, err := uu.userRepo.GetUserByID(ctx, id)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return errors.ErrRecordNotFound.Error()
	}

	if user.Status == "DEACTIVATED" {
		if err := uu.userRepo.UpdateUserStatus(ctx, id, "ACTIVATED"); err != nil {
			return errors.ErrInternalServerError.Error()
		}
	} else if user.Status == "ACTIVATED" {
		if err := uu.userRepo.UpdateUserStatus(ctx, id, "DEACTIVATED"); err != nil {
			return errors.ErrInternalServerError.Error()
		}
	}

	return nil
}

// UpdateAdmin updates an admin.
func (uu *UserUpdater) UpdateAdmin(ctx context.Context, user *entity.User, roleID uuid.UUID) error {
	if err := uu.userRepo.UpdateUser(ctx, user); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	userRole, err := uu.userRoleRepo.FindByUserID(ctx, user.ID)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	userRole.RoleID = roleID

	if err := uu.userRoleRepo.Update(ctx, userRole); err != nil {
		return nil
	}

	return nil
}

// UpdateRole updates a role
func (uu *UserUpdater) UpdateRole(ctx context.Context, id uuid.UUID, name string, permissionIDs []uuid.UUID) error {
	role, err := uu.roleRepo.FindByID(ctx, id)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if role == nil {
		return errors.ErrRecordNotFound.Error()
	}

	roleRequest := entity.NewRole(role.ID, name, "system")

	newPermissions := make([]*entity.RolePermission, 0)
	for _, pid := range permissionIDs {
		newPermissions = append(newPermissions, entity.NewRolePermission(
			uuid.New(),
			role.ID,
			pid,
			role.CreatedBy.String,
		))
	}

	if err := uu.roleRepo.Update(ctx, roleRequest, newPermissions); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// UpdatePermission updates a permission
func (uu *UserUpdater) UpdatePermission(ctx context.Context, id uuid.UUID, name, label string) error {
	permission, err := uu.permissionRepo.FindByID(ctx, id)

	if err != nil {
		return errors.ErrInternalServerError.Error()
	}

	if permission == nil {
		return errors.ErrRecordNotFound.Error()
	}

	newPermission := entity.NewPermission(id, name, label, permission.CreatedBy.String)

	if err := uu.permissionRepo.Update(ctx, newPermission); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

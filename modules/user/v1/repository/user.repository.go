package repository

import (
	"context"
	"fmt"
	"gin-starter/common/constant"
	"gin-starter/entity"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository is a repository for user
type UserRepository struct {
	db *gorm.DB
}

// UserRepositoryUseCase is a use case for user
type UserRepositoryUseCase interface {
	// GetUserByEmail is a function to get user by email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	// GetUserByID is a function to get user by id
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	// GetUserByForgotPasswordToken is a function to get user by forgot password token
	GetUserByForgotPasswordToken(ctx context.Context, token string) (*entity.User, error)
	// Update is a function to update user
	Update(ctx context.Context, user *entity.User) error
	// ChangePassword is a function to change password
	ChangePassword(ctx context.Context, user *entity.User, newPassword string) error
	// UpdateOTP is a function to update otp
	UpdateOTP(ctx context.Context, user *entity.User, otp string) error
	// CreateUser is a function to create user
	CreateUser(ctx context.Context, user *entity.User) error
	// GetUsers is a function to get users
	GetUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error)
	// GetAdminUsers is a function to get admin users
	GetAdminUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error)
	// UpdateUser is a function to update user
	UpdateUser(ctx context.Context, user *entity.User) error
	// UpdateUserStatus is a function to update user status
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error
	// DeleteAdmin is a function to delete admin user
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// GetUserByEmail is a function to get user by email
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	result := new(entity.User)

	if err := ur.db.
		WithContext(ctx).
		Preload("UserRole").
		Preload("UserRole.Role").
		Preload("Employee").
		Preload("Employee.CustomerBranch").
		Where("email = ?", email).
		Find(result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[UserRepository-GetUserByEmail] email not found")
	}

	return result, nil
}

// GetUserByID is a function to get user by id
func (ur *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	result := new(entity.User)

	if err := ur.db.
		WithContext(ctx).
		Preload("UserRole").
		Preload("UserRole.Role").
		Preload("Employee").
		Preload("Employee.CustomerBranch").
		Preload("Employee.CustomerBranch.Province").
		Preload("Employee.CustomerBranch.Regency").
		Preload("Employee.CustomerBranch.District").
		Preload("Employee.CustomerBranch.Village").
		Where("id = ?", id).
		First(result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[UserRepository-GetUserByID] user not found")
	}

	return result, nil
}

// GetUserByForgotPasswordToken is a function to get user by forgot password token
func (ur *UserRepository) GetUserByForgotPasswordToken(ctx context.Context, token string) (*entity.User, error) {
	result := new(entity.User)

	if err := ur.db.
		WithContext(ctx).
		Preload("UserRole").
		Preload("UserRole.Role").
		Preload("CustomerBranch").
		Where("forgot_password_token = ?", token).
		First(result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[UserRepository-GetUserByID] user not found")
	}

	return result, nil
}

// UpdateOTP is a function to update otp
func (ur *UserRepository) UpdateOTP(ctx context.Context, user *entity.User, otp string) error {
	if err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where(`id = ?`, user.ID).
		Updates(
			map[string]interface{}{
				"otp":        otp,
				"updated_at": time.Now(),
			}).Error; err != nil {
		return errors.Wrap(err, "[UserRepository-Update] error when updating user data")
	}
	return nil
}

// Update is a function to update user
func (ur *UserRepository) Update(ctx context.Context, user *entity.User) error {
	oldTime := user.UpdatedAt
	user.UpdatedAt = time.Now()
	if err := ur.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			sourceModel := new(entity.User)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sourceModel, user.ID).Error; err != nil {
				log.Println("[GamPTKRepository - Update]", err)
				return err
			}
			if err := tx.Model(&entity.User{}).
				Where(`id`, user.ID).
				UpdateColumns(sourceModel.MapUpdateFrom(user)).Error; err != nil {
				log.Println("[GamPTKRepository - Update]", err)
				return err
			}
			return nil
		}); err != nil {
		user.UpdatedAt = oldTime
	}
	return nil
}

// ChangePassword is a function to change password
func (ur *UserRepository) ChangePassword(ctx context.Context, user *entity.User, newPassword string) error {
	if err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where(`id = ?`, user.ID).
		Updates(
			map[string]interface{}{
				"password":   newPassword,
				"updated_at": time.Now(),
			}).Error; err != nil {
		return errors.Wrap(err, "[UserRepository-Update] error when updating user data")
	}

	return nil
}

// CreateUser is a function to create user
func (ur *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if err := ur.db.
		WithContext(ctx).
		Model(&entity.User{}).
		Create(user).
		Error; err != nil {
		return errors.Wrap(err, "[UserRepository-CreateUser] error while creating user")
	}

	return nil
}

// GetUsers is a function to get all users
func (ur *UserRepository) GetUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error) {
	var user []*entity.User
	var total int64
	var gormDB = ur.db.
		WithContext(ctx).
		Model(&entity.User{}).
		Joins("left join main.user_roles on users.id=user_roles.user_id").
		Where("main.user_roles.user_id is null")

	gormDB.Count(&total)

	gormDB = gormDB.Limit(limit).
		Offset(offset)

	if query != "" {
		gormDB = gormDB.
			Where("name ILIKE ?", "%"+query+"%").
			Or("email ILIKE ?", "%"+query+"%").
			Or("phone_number ILIKE ?", "%"+query+"%")
	}

	if order != constant.Ascending && order != constant.Descending {
		order = constant.Descending
	}

	if sort == "" {
		sort = "created_at"
	}

	gormDB = gormDB.Order(fmt.Sprintf("%s %s", sort, order))

	if err := gormDB.Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, errors.Wrap(err, "[UserRepository-GetAdminUsers] error when looking up all user")
	}

	return user, total, nil
}

// GetAdminUsers is a function to get all admin users
func (ur *UserRepository) GetAdminUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error) {
	var user []*entity.User
	var total int64
	var gormDB = ur.db.
		WithContext(ctx).
		Model(&entity.User{}).
		Preload("UserRole").
		Preload("UserRole.Role.RolePermissions").
		Preload("UserRole.Role.RolePermissions.Permission").
		Joins("inner join main.user_roles on main.users.id=main.user_roles.user_id")

	gormDB.Count(&total)

	gormDB = gormDB.Limit(limit).
		Offset(offset)

	if query != "" {
		gormDB = gormDB.
			Where("name ILIKE ?", "%"+query+"%").
			Or("email ILIKE ?", "%"+query+"%").
			Or("phone_number ILIKE ?", "%"+query+"%")
	}

	if order != constant.Ascending && order != constant.Descending {
		order = constant.Descending
	}

	if sort == "" {
		sort = "created_at"
	}

	gormDB = gormDB.Order(fmt.Sprintf("%s %s", sort, order))

	if err := gormDB.Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, errors.Wrap(err, "[UserRepository-GetAdminUsers] error when looking up all user")
	}

	return user, total, nil
}

// UpdateUser is a function to update user
func (ur *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	oldTime := user.UpdatedAt
	user.UpdatedAt = time.Now()
	if err := ur.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			sourceModel := new(entity.User)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sourceModel, user.ID).Error; err != nil {
				log.Println("[UserRepository-UpdateUser]", err)
				return err
			}
			if err := tx.Model(&entity.User{}).
				Where(`id`, user.ID).
				UpdateColumns(sourceModel.MapUpdateFrom(user)).Error; err != nil {
				log.Println("[UserRepository-UpdateUser]", err)
				return err
			}
			return nil
		}); err != nil {
		user.UpdatedAt = oldTime
	}
	return nil
}

// UpdateUserStatus is a function to update user status
func (ur *UserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error {
	if err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where(`id = ?`, id).
		Updates(
			map[string]interface{}{
				"status":     status,
				"updated_at": time.Now(),
			}).Error; err != nil {
		return errors.Wrap(err, "[UserRepository-DeactivateUser] error when updating user data")
	}

	return nil
}

// DeleteAdmin is a function to delete admin user
func (ur *UserRepository) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	if err := ur.db.WithContext(ctx).
		Model(&entity.User{}).
		Where(`id = ?`, id).
		Delete(&entity.User{}, "id = ?", id).Error; err != nil {
		return errors.Wrap(err, "[UserRepository-DeleteAdmin] error when updating user data")
	}

	return nil
}

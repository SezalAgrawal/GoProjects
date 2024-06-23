package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/lib/utils"
)

type UserStore interface {
	CreateUser(ctx context.Context, db *gorm.DB, user *model.User) (*model.User, error)
	CreateUserRole(ctx context.Context, db *gorm.DB, userID, roleID string) (*model.UserRole, error)
	GetUserByUsernamePassword(ctx context.Context, db *gorm.DB, name, hashedPassword string) (*model.User, error)
	GenerateAccessToken(ctx context.Context, db *gorm.DB, userID string) (*model.AccessToken, error)
	GetUserByAccessToken(ctx context.Context, db *gorm.DB, accessToken string) (*model.User, error)
	DeleteAllAccessTokens(ctx context.Context, db *gorm.DB, userID string) error
}

type userStore struct {
}

func NewUserStore() UserStore {
	return &userStore{}
}

func (e *userStore) CreateUser(ctx context.Context, db *gorm.DB, user *model.User) (*model.User, error) {
	usr := userModelToDB(user)

	if err := db.WithContext(ctx).Create(&usr).Error; err != nil {
		return nil, fmt.Errorf("creating user in store %w", err)
	}

	return usr.dbToModel(), nil
}

func (e *userStore) CreateUserRole(ctx context.Context, db *gorm.DB, userID, roleID string) (*model.UserRole, error) {
	usrRole := &userRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := db.WithContext(ctx).Create(&usrRole).Error; err != nil {
		return nil, fmt.Errorf("creating user role in store %w", err)
	}

	return usrRole.dbToModel(), nil
}

func (e *userStore) GetUserByUsernamePassword(ctx context.Context, db *gorm.DB, name, hashedPassword string) (*model.User, error) {
	usr := new(user)

	if err := db.WithContext(ctx).Model(&user{}).Where("name = ? and hashed_password = ?", name, hashedPassword).First(usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching user in store %w", err)
	}

	return usr.dbToModel(), nil
}

func (e *userStore) GenerateAccessToken(ctx context.Context, db *gorm.DB, userID string) (*model.AccessToken, error) {
	token := &accessToken{
		UserID: userID,
		Token:  utils.NewKSUID(),
	}

	if err := db.WithContext(ctx).Create(&token).Error; err != nil {
		return nil, fmt.Errorf("creating access token in store %w", err)
	}

	return token.dbToModel(), nil
}

func (e *userStore) GetUserByAccessToken(ctx context.Context, db *gorm.DB, userAccessToken string) (*model.User, error) {
	token := new(accessToken)

	if err := db.WithContext(ctx).Model(&accessToken{}).Where("token = ? and deleted = false", userAccessToken).First(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching user in store %w", err)
	}

	usr := new(user)

	if err := db.WithContext(ctx).Model(&user{}).
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Where("id = ?", token.UserID).
		First(usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching user in store %w", err)
	}

	return usr.dbToModel(), nil
}

func (e *userStore) DeleteAllAccessTokens(ctx context.Context, db *gorm.DB, userID string) error {
	if err := db.WithContext(ctx).Model(&accessToken{}).Where("user_id = ?", userID).Updates(map[string]interface{}{"deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return fmt.Errorf("update access token in store %w", err)
	}

	return nil
}

type user struct {
	ID             string      `gorm:"column:id"`
	Name           string      `gorm:"column:name"`
	HashedPassword string      `gorm:"column:hashed_password"`
	CreatedAt      time.Time   `gorm:"column:created_at"`
	UpdatedAt      time.Time   `gorm:"column:updated_at"`
	UserRoles      []*userRole `gorm:"foreignkey:UserID"`
}

type userRole struct {
	ID        string    `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	RoleID    string    `gorm:"column:role_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Role      *role     `gorm:"foreignkey:RoleID"`
}

type accessToken struct {
	ID        string    `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	Deleted   bool      `gorm:"column:deleted"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type role struct {
	ID        string    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (e *user) BeforeCreate(_ *gorm.DB) error {
	e.ID = "usr_" + utils.NewKSUID()
	return nil
}

func (e *userRole) BeforeCreate(_ *gorm.DB) error {
	e.ID = "usrRol_" + utils.NewKSUID()
	return nil
}

func (e *accessToken) BeforeCreate(_ *gorm.DB) error {
	e.ID = "tok_" + utils.NewKSUID()
	return nil
}

func (e *user) dbToModel() *model.User {
	usr := &model.User{
		ID:             e.ID,
		Name:           e.Name,
		HashedPassword: e.HashedPassword,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}

	var userRoles []*model.UserRole
	for _, usrRole := range e.UserRoles {
		userRoles = append(userRoles, usrRole.dbToModel())
	}
	usr.UserRoles = userRoles
	return usr
}

func userModelToDB(model *model.User) *user {
	return &user{
		ID:             model.ID,
		Name:           model.Name,
		HashedPassword: model.HashedPassword,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}

func (e *userRole) dbToModel() *model.UserRole {
	usrRole := &model.UserRole{
		ID:        e.ID,
		UserID:    e.UserID,
		RoleID:    e.RoleID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	if e.Role != nil {
		usrRole.Role = e.Role.dbToModel()
	}
	return usrRole
}

func (e *accessToken) dbToModel() *model.AccessToken {
	return &model.AccessToken{
		ID:        e.ID,
		UserID:    e.UserID,
		Token:     e.Token,
		Deleted:   e.Deleted,
		DeletedAt: e.DeletedAt,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (e *role) dbToModel() *model.Role {
	return &model.Role{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

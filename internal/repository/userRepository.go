package repository

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(usr domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserbyID(id int) (domain.User, error)
	UpdateUser(id int, usr domain.User) (domain.User, error)
	AddBankAccount(e domain.BankAccount) error

	//Cart
	CreateCart(input domain.Cart) error
	FindCartItems(userId int) ([]*domain.Cart, error)
	FindCartItem(userid int, prdctId int) (*domain.Cart, error)
	UpdateCart(input domain.Cart) error
	DeleteCartItemByid(Id int) error
	DeleteCartItems(userId int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// User Cart
func (r *userRepository) CreateCart(input domain.Cart) error {

	result := r.db.Create(&input)
	if result.Error != nil {
		return fmt.Errorf("cart creation failed due to %v", result.Error)
	}

	return nil
}

func (r *userRepository) FindCartItem(userid int, prdctId int) (*domain.Cart, error) {
	var cartItem *domain.Cart
	result := r.db.Where("user_id=? AND product_id=?", userid, prdctId).First(cartItem)
	if result.Error != nil {
		return &domain.Cart{}, fmt.Errorf("finding cart failed due to %v", result.Error)
	}

	return cartItem, nil
}

func (r *userRepository) FindCartItems(userId int) ([]*domain.Cart, error) {
	var allCartItems []*domain.Cart

	result := r.db.Where("user_id=?", userId).Find(allCartItems)
	if result.Error != nil {
		return []*domain.Cart{}, fmt.Errorf("finding cart failed due to %v", result.Error)
	}

	return allCartItems, nil
}

func (r *userRepository) UpdateCart(input domain.Cart) error {
	var cart domain.Cart
	result := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id=?", input.ID).Updates(input)
	if result.Error != nil {
		return fmt.Errorf("updating cart failed due to %v", result.Error)
	}

	return nil
}
func (r *userRepository) DeleteCartItemByid(Id int) error {

	result := r.db.Delete(&domain.Cart{}, Id)
	if result.Error != nil {
		return fmt.Errorf("cart deletion failed due to %v", result.Error)
	}

	return nil
}

func (r *userRepository) DeleteCartItems(userId int) error {

	result := r.db.Where("user_id=?", userId).Delete(&domain.Cart{})
	if result.Error != nil {
		return fmt.Errorf("user cart items deletion failed due to %v", result.Error)
	}

	return nil
}

// Addind bank account for seller feature
func (r *userRepository) AddBankAccount(e domain.BankAccount) error {
	return r.db.Create(&e).Error
}

func (r *userRepository) CreateUser(usr domain.User) (domain.User, error) {

	result := r.db.Create(&usr) // Just use the input directly
	if result.Error != nil {
		log.Println("User creation failed due to", result.Error)
		return domain.User{}, result.Error
	}

	return usr, nil
}

func (r *userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	//Query the database for a user with the given email
	result := r.db.Where("email=?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with email %s not found", email)
			return domain.User{}, fmt.Errorf("user not found")
		}

		log.Printf("Error finding user :%v", result.Error)
		return domain.User{}, fmt.Errorf("database error: %w", result.Error)
	}

	return user, nil
}

func (r *userRepository) FindUserbyID(id int) (domain.User, error) {
	var user domain.User

	result := r.db.Where("id=?", id).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with id %d not found", id)
			return domain.User{}, fmt.Errorf("user not found")
		}

		log.Printf("Error finding user :%v", result.Error)
		return domain.User{}, fmt.Errorf("database error: %w", result.Error)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(id int, usr domain.User) (domain.User, error) {
	var existingUser domain.User

	//check if user with id exist
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("user with ID %d not found", id)
		}

		return domain.User{}, fmt.Errorf("database error: %w", err)
	}

	//Update the user fields
	result := r.db.Model(&domain.User{}).Where("id=?", id).Updates(usr)
	if result.Error != nil {
		log.Printf("Error updating user: %v", result.Error)
		return domain.User{}, fmt.Errorf("failed to update user: %w", result.Error)
	}

	//Verify the update was successfull
	if result.RowsAffected == 0 {
		return domain.User{}, fmt.Errorf("no user was updated")
	}

	//Retrieve and return the updated user
	var updatedUser domain.User
	if err := r.db.First(&updatedUser, id).Error; err != nil {
		return domain.User{}, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return updatedUser, nil

}

package main

import (
	"echo-demo/internal/model"
	"errors"

	"gorm.io/gorm"
)

func BootstrapRolesAndPermissions(db *gorm.DB) error {
	roles := []model.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "theater_manager", Description: "Manager of theaters and seats"},
		{Name: "customer", Description: "Customer who can book and cancel tickets"},
	}

	for _, role := range roles {
		var existingRole model.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	permissions := []model.Permission{
		// Permissions cho Movie
		{Name: "create_movie", Description: "Create a new movie"},
		{Name: "edit_movie", Description: "Edit movie details"},
		{Name: "delete_movie", Description: "Delete a movie"},
		{Name: "view_movie", Description: "View movies"},

		// Permissions cho Theater
		{Name: "create_theater", Description: "Create a new theater"},
		{Name: "edit_theater", Description: "Edit theater details"},
		{Name: "delete_theater", Description: "Delete a theater"},
		{Name: "view_theater", Description: "View theaters"},

		// Permissions cho Seat
		{Name: "create_seat", Description: "Create a new seat"},
		{Name: "edit_seat", Description: "Edit seat details"},
		{Name: "delete_seat", Description: "Delete a seat"},
		{Name: "view_seat", Description: "View seats"},

		// Permissions cho Payment
		{Name: "view_payment", Description: "View payment details"},
		{Name: "manage_payment", Description: "Manage payments (approve/cancel)"},

		// Permissions cho Movie Category
		{Name: "manage_movie_category", Description: "Manage movie categories"},

		// Permissions cho Booking (đặt/hủy vé)
		{Name: "book_ticket", Description: "Book movie ticket"},
		{Name: "cancel_ticket", Description: "Cancel booked ticket"},

		// Permissions cho báo cáo
		{Name: "view_report", Description: "View system reports"},
	}

	for _, perm := range permissions {
		var existingPerm model.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&perm).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	var adminRole, theaterManagerRole, customerRole model.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "theater_manager").First(&theaterManagerRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
		return err
	}

	var allPermissions []model.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}
	if err := db.Model(&adminRole).Association("Permissions").Replace(allPermissions); err != nil {
		return err
	}

	theaterManagerPermNames := []string{
		"create_theater", "edit_theater", "delete_theater", "view_theater",
		"create_seat", "edit_seat", "delete_seat", "view_seat",
		"view_payment",
		"view_movie",
		"view_report",
	}
	var theaterManagerPermissions []model.Permission
	if err := db.Where("name IN ?", theaterManagerPermNames).Find(&theaterManagerPermissions).Error; err != nil {
		return err
	}
	if err := db.Model(&theaterManagerRole).Association("Permissions").Replace(theaterManagerPermissions); err != nil {
		return err
	}

	customerPermNames := []string{
		"book_ticket", "cancel_ticket",
		"view_movie", "view_theater", "view_seat",
	}
	var customerPermissions []model.Permission
	if err := db.Where("name IN ?", customerPermNames).Find(&customerPermissions).Error; err != nil {
		return err
	}
	if err := db.Model(&customerRole).Association("Permissions").Replace(customerPermissions); err != nil {
		return err
	}

	return nil
}

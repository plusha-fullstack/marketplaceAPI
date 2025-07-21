package utils

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

func ValidateLogin(login string) error {
	if len(login) < 3 || len(login) > 50 {
		return errors.New("login must be between 3 and 50 characters")
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", login)
	if !matched {
		return errors.New("login can only contain letters, numbers, and underscores")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 100 {
		return errors.New("password must be between 6 and 100 characters")
	}
	return nil
}

func ValidateAdTitle(title string) error {
	title = strings.TrimSpace(title)
	if len(title) < 1 || len(title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}
	return nil
}

func ValidateAdDescription(description string) error {
	description = strings.TrimSpace(description)
	if len(description) < 1 || len(description) > 2000 {
		return errors.New("description must be between 1 and 2000 characters")
	}
	return nil
}

func ValidateImageURL(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	if len(imageURL) > 500 {
		return errors.New("image URL too long")
	}

	_, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return errors.New("invalid image URL format")
	}

	return nil
}

func ValidatePrice(price float64) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	if price > 1000000 {
		return errors.New("price too high")
	}
	return nil
}

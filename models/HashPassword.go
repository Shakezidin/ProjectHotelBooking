package models

import "golang.org/x/crypto/bcrypt"

func (owner *Owner) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	owner.Password = string(bytes)
	return nil
}
func (owner *Owner) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (admin *Admin) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	admin.Password = string(bytes)
	return nil
}
func (admin *Admin) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func HashPassworddecript(password string) (string, error) {
    // Generate a salted and hashed version of the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    
    // Convert the hashed password to a string
    hashedPasswordStr := string(hashedPassword)
    
    return hashedPasswordStr, nil
}

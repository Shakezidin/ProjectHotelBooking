package models

import "golang.org/x/crypto/bcrypt"

//HashPassword create hashedpassword for Owner
func (owner *Owner) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	owner.Password = string(bytes)
	return nil
}

//CheckPassword check hashedpassword of Owner
func (owner *Owner) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

//HashPassword create hashedpassword for User
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

//CheckPassword check hashedpassword of User
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

//HashPassword create hashedpassword for Admin
func (admin *Admin) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	admin.Password = string(bytes)
	return nil
}

//CheckPassword check hashedpassword of Admin
func (admin *Admin) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}


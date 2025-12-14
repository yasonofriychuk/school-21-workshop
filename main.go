package main

import (
	"errors"
	"fmt"
	"school21-errors/internal/validation"
)

func main() {
	f := validation.RegistrationForm{
		Age:   16,
		Name:  "",
		Phone: "+7987654321",
	}

	err := f.CheckV2()

	if errors.Is(err, validation.InvalidAgeErr) {
		fmt.Println("Регистрация доступна для лиц старше 18 лет")
	}
	if errors.Is(err, validation.InvalidPhoneErr) {
		fmt.Println("Проверьте корректность введенного номера телефона")
	}
	if errors.Is(err, validation.InvalidNameErr) {
		fmt.Println("Проверьте правильность введенного имени")
	}

	var errT validation.ValidationError
	if errors.As(err, &errT) {
		fmt.Println(errT.ValidAge)
	}
}

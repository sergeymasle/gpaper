package main

import (
	"fmt"

)

func ExampleRuTranslit() {
	tests := []string{
		"Проверочная СТРОКА для транслитерации",
		"ЧАЩА",
		"ЧаЩа",
		"Чаща",
		"чаЩА",
	}
	for _, text := range tests {
		fmt.Println(Ru(text))
	}
	// Output:
	// Proverochnaja STROKA dlja transliteracii
	// CHASCHA
	// ChaScha
	// Chascha
	// chaSCHA
}

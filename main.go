package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Products (name, price) values ('Огурец', 100)")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId()) // не поддерживается
	fmt.Println(result.RowsAffected()) // количество добавленных строк
}

// РЕБЯТА ДОМАШНЕЕ ЗАДАНИЕ СЛУШАЕМ НАПИСАТ КРУД ДЛЯ ТОВАРОВ ПОДНЯТИЕ HTTP СЕРВЕРА ВНУТРИ ПРИЛОЖЕНИЯ (ВЫБРАТЬ ФРЕЙМВОРК gin,gorilla/mux)

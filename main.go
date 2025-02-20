package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // импорт драйвера без использования его напрямую
)

// Todo представляет задачу в нашем списке дел
type Todo struct {
	ID        int
	Task      string
	Completed bool
}

func main() {
	// Открываем или создаем базу данных (файл todo.db)
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal("Ошибка открытия базы данных:", err)
	}
	defer db.Close()

	// Создаем таблицу todos, если ее нет
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			task TEXT,
			completed BOOLEAN
		);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
	fmt.Println("Таблица todos готова к работе!")

	// Вставляем новую задачу
	insertSQL := `INSERT INTO todos (task, completed) VALUES (?, ?)`
	result, err := db.Exec(insertSQL, "Купить молоко", false)
	if err != nil {
		log.Fatal("Ошибка вставки задачи:", err)
	}
	newID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Ошибка получения ID новой задачи:", err)
	}
	fmt.Printf("Добавлена новая задача с ID: %d\n", newID)

	// Выбираем и выводим все задачи
	rows, err := db.Query("SELECT id, task, completed FROM todos")
	if err != nil {
		log.Fatal("Ошибка выборки задач:", err)
	}
	defer rows.Close()

	fmt.Println("Список задач:")
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Task, &todo.Completed)
		if err != nil {
			log.Fatal("Ошибка чтения строки:", err)
		}
		fmt.Printf("ID: %d | Задача: %s | Выполнена: %v\n", todo.ID, todo.Task, todo.Completed)
	}

	// Обработка ошибок при итерации
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

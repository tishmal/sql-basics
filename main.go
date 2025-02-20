package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sql-basics/tasks"

	_ "github.com/mattn/go-sqlite3" // импорт драйвера без использования его напрямую
)

func main() {
	allTasks := flag.Bool("tasks", false, "all tasks")
	completedID := flag.Int("complete", 0, "complete task by ID")
	deletedID := flag.Int("delete", 0, " delete task by ID")
	flag.Parse()

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

	if *allTasks {
		tasks.AllTasks(db)
		return
	}
	if *completedID > 0 {
		tasks.CompleteTask(db, *completedID) // Завершаем задачу по ID
		return
	}
	if *deletedID > 0 {
		tasks.DeleteTask(db, *deletedID)
		return
	}

	fmt.Print("Введите задачу: ")
	var task string
	_, err = fmt.Scanln(&task)
	if err != nil {
		log.Fatal("Ошибка чтения ввода:", err)
	}
	tasks.AddTask(db, task, false)
}

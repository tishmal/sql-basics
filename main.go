package main

import (
	"database/sql"
	"flag"
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
	tasks := flag.Bool("tasks", false, "all tasks")
	completedID := flag.Int("complete", 0, "complete task by ID")
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

	if *tasks {
		AllTasks(db)
		return
	}
	if *completedID > 0 {
		CompleteTask(db, *completedID) // Завершаем задачу по ID
		return
	}

	fmt.Print("Введите задачу: ")
	var task string
	_, err = fmt.Scanln(&task)
	if err != nil {
		log.Fatal("Ошибка чтения ввода:", err)
	}
	AddTask(db, task, false)
}

func AddTask(db *sql.DB, task string, completed bool) {
	insertSQL := `INSERT INTO todos (task, completed) VALUES (?, ?)`
	result, err := db.Exec(insertSQL, task, completed)
	if err != nil {
		log.Fatal("Error inserting task")
	}
	newID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Error get ID new task", err)
	}
	fmt.Printf("Adding new task ID: %d\n", newID)
}

func AllTasks(db *sql.DB) {
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

func CompleteTask(db *sql.DB, taskID int) {
	updateSQL := `UPDATE todos SET completed = 1 WHERE id = ?`
	_, err := db.Exec(updateSQL, taskID)
	if err != nil {
		log.Fatal("Error update status task", err)
	}
	log.Printf("Task with ID: %d была выполнена!\n", taskID)
}

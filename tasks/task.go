package tasks

import (
	"database/sql"
	"fmt"
	"log"
)

// Todo представляет задачу в нашем списке дел
type Todo struct {
	ID        int
	Task      string
	Completed bool
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
	log.Printf("Task with ID: %d completed!\n", taskID)
}

func DeleteTask(db *sql.DB, taskID int) {
	updateSQL := `DELETE FROM todos WHERE id = ?`
	_, err := db.Exec(updateSQL, taskID)
	if err != nil {
		log.Fatal("Error delete task", err)
	}
	log.Printf("Task with ID: %d deleted!\n", taskID)
}

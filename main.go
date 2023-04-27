package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

type ToDo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ToDos []ToDo

func (l *ToDos) Add(title string) {
	todo := ToDo{
		Title:     title,
		Completed: false,
	}
	*l = append(*l, todo)
}

// Обработка пути /
func index(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// list == слайс ToDo
	err := tmpl.Execute(w, list)

	if err != nil {
		_, err := fmt.Fprintf(w, err.Error())
		if err != nil {
			return
		}
	}
}

// Добавление нового элемента в бд
func insert(db *sql.DB, title string) error {
	_, err := db.Exec("INSERT INTO `todo` (`id`, `title`, `completed`) VALUES (NULL, title, '0')")
	if err != nil {
		return err
	}
	return nil
}

// Получение всех элементов из бд
func getAll(db *sql.DB) (slice []ToDo, e error) {
	res, err := db.Query("SELECT `title`, `completed` FROM `todo`")
	if err != nil {
		return nil, err
	}

	// перебор значений с запроса
	// Next() проверяет есть ли следующие элементы
	for res.Next() {
		var todo ToDo
		// Scan(значения с которыми сверяемся) забирает данные из бд
		err = res.Scan(&todo.Title, &todo.Completed)
		if err != nil {
			return nil, err
		}
		slice = append(slice, todo)
	}
	return slice, nil
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// Хендлер для пути /
	http.HandleFunc("/", index)

	// Запуск сервера
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error with server starting")
	}

}

//TODO: прикрепить к хендлеру работу с бд

//TODO: ууу сука, сделаю бота для работы через телеграм

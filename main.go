package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

// Добавление нового элемента в бд
func insert(db *sql.DB, title string) error {
	_, err := db.Exec("INSERT INTO `todo` (`id`, `title`, `completed`) VALUES (NULL, ?, '0')", title)
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

// Обработка пути /
func index(w http.ResponseWriter, r *http.Request) {
	// Получение всех элементов из базы данных
	list, err := getAll(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err = tmpl.Execute(w, list)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Обработка POST-запроса из формы
func save(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Получение данных из тела запроса
		data := r.FormValue("data")

		// Сохранение данных в базу данных
		err := insert(db, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Перенаправление на главную страницу
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//func delete(w http.ResponseWriter, r *http.Request) {
//
//}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/save", save)

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: добавить readme

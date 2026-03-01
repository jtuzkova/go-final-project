package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
    ID      string `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
    var id int64
    query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
    res, err := db.Exec(query, sql.Named("date", task.Date),
							sql.Named("title", task.Title),
							sql.Named("comment", task.Comment),
							sql.Named("repeat", task.Repeat))
    if err == nil {
        id, err = res.LastInsertId()
    }
    return id, err
}

func GetTask(id string) (*Task, error) {
    task := &Task{}

    row := db.QueryRow(`SELECT * FROM scheduler
                            WHERE id = :id`, sql.Named("id", id))
    err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
    if err != nil {
        return task, err 
	}
	return task, nil   
}

func UpdateTask(task *Task) error {
    query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat 
              WHERE id = :id`
    res, err := db.Exec(query, 
        sql.Named("date", task.Date),
        sql.Named("title", task.Title),
        sql.Named("comment", task.Comment),
        sql.Named("repeat", task.Repeat),
        sql.Named("id", task.ID))
    if err != nil {
        return err
    }

    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf(`incorrect id for updating task`)
    }
    return nil
}

func UpdateDate(next string, id string) error {
    query := `UPDATE scheduler SET date = :date
              WHERE id = :id`
    res, err := db.Exec(query, 
        sql.Named("date", next),
        sql.Named("id", id))
    if err != nil {
        return err
    }

    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf(`incorrect id for updating`)
    }
    return nil
}

func DeleteTask(id string) error {
    res, err := db.Exec(`DELETE FROM scheduler WHERE id = :id`,
                        sql.Named("id", id))

    if err != nil {
        return err
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf(`incorrect id for updating`)
    }
    return nil
}

func Tasks(limit int) ([]*Task, error) {
    var tasks []*Task

    rows, err := db.Query(`SELECT id, date, title, comment, repeat 
                            FROM scheduler ORDER BY date
                            LIMIT :limit`, sql.Named("limit", limit))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        task := &Task{}
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, err
        }

        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    if tasks == nil {
        tasks = make([]*Task, 0)
    }
    return tasks, nil
}

func SearchTaskString(searchString string, limit int) ([]*Task, error) {
    var tasks []*Task

    rows, err := db.Query(`SELECT * FROM scheduler 
                            WHERE title LIKE :search OR 
                            comment LIKE :search 
                            ORDER BY date LIMIT :limit`, 
                            sql.Named("search", searchString),
                            sql.Named("limit", limit))
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        task := &Task{}
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, err
        }

        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    if tasks == nil {
        tasks = make([]*Task, 0)
    }
    return tasks, nil
}

func SearchTaskDate(searchDate string, limit int) ([]*Task, error) {
    var tasks []*Task

    rows, err := db.Query(`SELECT * FROM scheduler 
                            WHERE date = :date
                            ORDER BY date LIMIT :limit`, 
                            sql.Named("date", searchDate),
                            sql.Named("limit", limit))
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        task := &Task{}
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, err
        }

        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    if tasks == nil {
        tasks = make([]*Task, 0)
    }
    return tasks, nil
}
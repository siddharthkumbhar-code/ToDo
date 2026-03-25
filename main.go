package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Taskid int32 `json:"taskid"`
	Name  string  `json:"name"`
	Taskdone bool  `json:"taskdone"`
}
var db *sql.DB

func AddTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	err:=json.NewDecoder(r.Body).Decode(&task)
	
	if err!=nil{
		http.Error(w,"Invalid json",400)
	}

	query:=`INSERT INTO tasks(name,taskdone) VALUES(?,?)`

	_,err=db.Exec(query,task.Name,task.Taskdone)
	if err != nil {
		http.Error(w, "Database Error", 500)
		return
	}
	w.Write([]byte("added the task"))
}

func ListAll(w http.ResponseWriter,r *http.Request){

	query:=`SELECT taskid,name,taskdone FROM tasks`
	rows,err:=db.Query(query)
	if err!=nil{
		http.Error(w,"invalid",500)
		return
	}
	defer rows.Close()
	var tasks []Task

	for rows.Next(){
		var task Task
		err:=rows.Scan(&task.Taskid,&task.Name,&task.Taskdone)
		if err!=nil{
			http.Error(w,err.Error(),500)
			return
		}
		tasks = append(tasks, task)
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(tasks)
}

func DeleteTask(w http.ResponseWriter,r *http.Request){

	
	id:=r.URL.Query().Get("id")

	if id==""{
		http.Error(w,"id required",400)
		return
	}
	taskID,err:=strconv.Atoi(id)
	if err!=nil{
		http.Error(w,"invalid id",405)
		return
	}

	query:=`DELETE FROM tasks
			WHERE taskid=?`
	_,err=db.Exec(query,taskID)
	if err!=nil{
		http.Error(w,err.Error(),500)
	}
	w.Write([]byte("delete"))
	

}

func RenameTask(w http.ResponseWriter,r *http.Request){
	
	id:=r.URL.Query().Get("id")
	if id==""{
		http.Error(w,"error",402)
		return
	}
	taskid,err:=strconv.Atoi(id)
	if err!=nil{
		http.Error(w,"invalid id ",402)
		return
	}

	var task Task
	err=json.NewDecoder(r.Body).Decode(&task)
	if err!=nil{
		http.Error(w,"error",400)
		return
	}

	query:=`UPDATE tasks 
			SET name=?
			WHERE taskid=?`

	_,err=db.Exec(query,task.Name,taskid)
	if err!=nil{
		http.Error(w,"database error",500)
		return
	}
	w.Write([]byte("Rename successfully"))
}

func updateStatus(w http.ResponseWriter,r *http.Request){

	
	id:=r.URL.Query().Get("id")
	if id==""{
		http.Error(w,"invalid",402)
		return
	}
	tid,err:=strconv.Atoi(id)
	if err!=nil{
		http.Error(w,"error",400)
		return
	}
	var task Task
	err=json.NewDecoder(r.Body).Decode(&task)
	if err!=nil{
		http.Error(w,"error",400)
		return
	}
	query:=`UPDATE tasks
			SET taskdone=?
			WHERE taskid=?`
	_,err=db.Exec(query,task.Taskdone,tid)
	if err!=nil{
		log.Println(err)
	http.Error(w, err.Error(), 500)
	return
	}	
	w.Write([]byte("status updated successfully"))

}

func updateallthedata()
{
	 fmt.Println("this update from soham ")
	 fmt.Println("this is from soham branch")
	fmt.Println("hiiii")
}

func main() {

	var err error
	db,err=sql.Open("sqlite3","./database.db")
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close()

	dataTable:=` CREATE TABLE IF NOT EXISTS tasks(
				taskid INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				taskdone bool NOT NULL
	);`
	

	_,err=db.Exec(dataTable)
	if err!=nil{
		log.Fatal(err)
	}

	
	http.HandleFunc("/add",AddTask)
	http.HandleFunc("/list",ListAll)
	http.HandleFunc("/delete",DeleteTask)
	http.HandleFunc("/rename",RenameTask)
	http.HandleFunc("/status",updateStatus)
	http.ListenAndServe(":8080",nil)

}

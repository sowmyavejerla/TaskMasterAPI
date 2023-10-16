package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Handler struct{
	DB *gorm.DB
}

type TaskStatus string

const (
	NotStarted string = "NotStarted"
	InProgress string = "InProgress"
	Completed  string = "Completed"
)

type Task struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description string     `gorm:"column:description" json:"description"`
	Status      TaskStatus `gorm:"column:status" json:"status"`
	CreatedAt     time.Time `gorm:"column:createdat" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedat" json:"updatedAt"`
}


func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB:db}
}

func (h *Handler)addTask(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

    var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err!=nil{
		http.Error(w,"Bad Request",http.StatusBadRequest)
		return
	}
	err= validateAddRequest(task)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	exists,err := TitleExists(h.DB,task.Title)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if exists{
		http.Error(w, "Task with this title already exists", http.StatusConflict)
		return
	}

	if err = CreateTask(h.DB,&task);err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w,"Task Created Successfully\n")
}

func (h *Handler)deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}
	taskid:=r.URL.Query().Get("id")
	if taskid ==""{
		http.Error(w,"Task ID not provided",http.StatusBadRequest)
		return
	}
	taskID,err:= strconv.Atoi(taskid)
	if err!=nil{
		http.Error(w,"Bad Request",http.StatusBadRequest)
		return
	}
	exists,err :=TaskExists(h.DB,taskID)
	if err!=nil{
		http.Error(w,"Internal Server error",http.StatusInternalServerError)
		return
	}

	if !exists{
		http.NotFound(w,r)
	}
	if err = DeleteTask(h.DB,taskID);err!=nil{
		http.Error(w,"Internal Server error",http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w,"Task Deleted Successfully\n")
}

func (h *Handler)listTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

	tasks,err := GetTasks(h.DB)
	if err!=nil{
		http.Error(w,"Error fetching tasks",http.StatusInternalServerError)
		return
	}

	taskJson,err:=json.Marshal(tasks)
	if err!=nil{
		http.Error(w,"Error in marshaling tasks",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskJson)
}

func isValidStatus(s TaskStatus) bool {
	switch s {
	case TaskStatus(NotStarted), TaskStatus(InProgress), TaskStatus(Completed):
		return true
	default:
		return false
	}
}

func (h *Handler)editTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

	var updateTask Task
	err := json.NewDecoder(r.Body).Decode(&updateTask)
	if err!=nil{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = validateEditRequest(updateTask)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	exists,err:=TaskExists(h.DB,updateTask.ID)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if !exists{
		http.NotFound(w,r)
		return
	}

	if err = UpdateTask(h.DB,updateTask);err!=nil{
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w,"Task Updated Successfully\n")
	
}

func validateEditRequest(task Task) error  {
	if !isValidStatus(TaskStatus(task.Status)){
		return  errors.New("invalid status")
	}
	if task.ID ==0 || task.Title =="" || task.Description ==""{
		return errors.New("please enter mandatory details")
	}
	return nil
}

func validateAddRequest(task Task) error {
	if task.Title =="" || task.Description =="" {
		return errors.New("please enter mandatory details")
	}
	if !isValidStatus(TaskStatus(task.Status)){
		return  errors.New("invalid status")
	}
	return nil
}
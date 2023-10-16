package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBname   string
	Port     string
}

type Repo struct{
	DB *gorm.DB
}

func InitializeDB(DB *DBConfig)(*gorm.DB,error) {
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	 DB.Host, DB.User, DB.Password, DB.DBname, DB.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	return db,nil
}

func TitleExists(DB *gorm.DB,title string)(bool,error)  {
	var count int64
	result := DB.Model(&Task{}).Where("title = ?",title).Count(&count)
	if result.Error !=nil{
		return false,result.Error
	}
	return count >0,nil
}

func TaskExists(DB *gorm.DB,taskID int)(bool,error)  {
	var count int64
	if err := DB.Model(&Task{}).Where("id = ?",taskID).Count(&count).Error;err!=nil{
		return false,err
	}
	return count >0 ,nil
}

func CreateTask(DB *gorm.DB,task *Task) error{
	exists,err := TitleExists(DB,task.Title)
	if err!=nil{
		return err
	}
	if exists{
		return errors.New("task already exists")
	}
	task.CreatedAt = time.Now().UTC()
	result:= DB.Create(task)
	if err := result.Error;err !=nil{
		return err
	}
	return nil
}


func DeleteTask(DB *gorm.DB,taskID int)error  {
	if err := DB.Where("id = ?",taskID).Delete(&Task{}).Error;err!=nil{
		fmt.Println("Error deleting task:", err)
		return err
	}
	return nil
}

func UpdateTask(DB *gorm.DB,task Task)error{
	if err := DB.Where("id = ?",task.ID).Updates(&task).Error; err!=nil{
		return err
	}
	return nil
}

func GetTasks(DB *gorm.DB)([]Task,error)  {
	var tasks []Task
	if err:=DB.Find(&tasks).Error;err!=nil{
		return nil,err
	}

	return tasks,nil
}


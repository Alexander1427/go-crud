package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task one",
		Content: "Some content",
	},
}

func getTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(res).Encode(tasks)
}

func getTaskUnic(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(res, "Invalid ID")
		return
	}
	for _, task := range tasks {
		if task.ID == taskID {
			res.Header().Set("Content-Type", "Application/json")
			json.NewEncoder(res).Encode(task)
		}
	}
}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(res, "Invalid ID")
		return
	}
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(res, "Delete complete")
		}
	}
}

func createTaks(res http.ResponseWriter, req *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(res, "Inserte datos validos")
	}
	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	res.Header().Set("Content-Type", "Application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(newTask)
}

func updateTask(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskID, err := strconv.Atoi(vars["id"])
	var updateTask task

	if err != nil {
		fmt.Fprintf(res, "Invalid ID")
		return
	}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(res, "Please incert valid data")
		return
	}
	json.Unmarshal(reqBody, &updateTask)

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updateTask.ID = taskID
			tasks = append(tasks, updateTask)

			fmt.Fprintf(res, "The update task %v", taskID)
		}
	}
}

func indexRoute(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to my api")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTask).Methods("GET")
	router.HandleFunc("/tasks", createTaks).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTaskUnic).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}

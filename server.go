package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
)

type Message struct {
	Id      int    `json:"id"`
	Message string `json:"-"` // Won't show that element according the task, - mean that json will skip this field
}

//Simple storage
var messages []Message

//Id of the last added message
var lastId = 0

func apiHandler() http.Handler {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/messages", messagesIndex).Methods("GET")
	router.HandleFunc("/messages", messagesAdd).Methods("POST")
	router.HandleFunc("/messages/{id}", messagesGetById).Methods("GET")

	return router

}

//Print out the messages list or error msg if none
func messagesIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	//No messages
	if len(messages) < 1 {
		fmt.Fprint(w, "No messages\n")
	}
	//Print messages one in line
	for _, m := range messages {
		fmt.Fprintf(w, "[%d] %s\n", m.Id, m.Message)
	}

}
//Add new message
func messagesAdd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//Read message, 10240 charset limit
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	//Call function which increment the ID and add new record to the messages array
	currentMessage, err := messageAddNew(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//Convert result to the json to be able to print result in json format according the task
	testJson, err := json.Marshal(currentMessage)

	fmt.Fprintf(w, "\n%s\n", string(testJson))

}
//Increment the ID number an add new record to the messages array
func messageAddNew(msg string) (Message, error) {
	//Initialize temp structure to be able to use append function
	tmpMessage := Message{}

	lastId += 1
	tmpMessage.Id = lastId
	tmpMessage.Message = msg

	messages = append(messages, tmpMessage)

	return tmpMessage, nil
}

//Get message by ID => /messages/ID
func messagesGetById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	//Read ID number
	vars := mux.Vars(r)
	//Convert id to the digit
	messageId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//ID = 0 we haven't record for that id
	if messageId < 1 {
		fmt.Fprint(w, "Id can't be lower the 1\n")
		return
	}
	//Call the function who'll try to find if ID exists and if so return the message according that ID
	message := messageFindById(messageId)
	if len(message) < 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "\nId %d not found\n", messageId)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "\n%s\n", message)

}
//Find if ID exist and return message if so
func messageFindById(id int) string {

	for _, m := range messages {
		if id == m.Id {
			return m.Message
		}
	}

	return ""
}

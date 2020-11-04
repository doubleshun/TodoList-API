//go:generate swagger generate spec
package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"

	"api/services"

	"github.com/zabawaba99/firego"
)

//TodoList is ...
var TodoList []Todo
var Fb firego.Firebase

type Todo struct {
	UserID string `json:"UserID" example:"User-1"`
	ItemID int64  `json:"ItemID" example:"0"`
	Item   string `json:"Item" example:"test"`
	State  string `json:"State" example : "Finished"`
}

type ApiResponse struct {
	ResultCode    string      `json:"ResultCode" example:"200"`
	ResultMessage interface{} `json:"ResultMessage"`
}

// @Summary Add todo to List
// @Tags Todo
// @Accept json
// @Produce json
// @param Body body Todo true "body"
// @Success 200 {array} ApiResponse
// @Router /api/todo [Post]
func AddTodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024)) //io.LimitReader限制大小
	if err != nil {
		fmt.Println(err)
	}
	var addTodo Todo
	_ = json.Unmarshal(body, &addTodo) //轉為json
	v := map[string]interface{}{
		"Item":  addTodo.Item,
		"State": addTodo.State,
	}
	Firebase, _ := Fb.Ref("/" + addTodo.UserID + "/" + strconv.FormatInt(int64(addTodo.ItemID), 10))
	err = Firebase.Set(v)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	response := ApiResponse{"200", TodoList}

	services.ResponseWithJson(w, http.StatusOK, response) //回傳

}

// @Summary Get all todo from List
// @Tags Todo
// @Accept json
// @Produce json
// @param UserID path string true "UserId"
// @Success 200 {array} ApiResponse
// @Router /api/todo/{UserID} [Get]
func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var v interface{}
	Firebase, _ := Fb.Ref("/" + vars["UserID"])
	if err := Firebase.Value(&v); err != nil {
		log.Fatal(err)
	}
	response := ApiResponse{"200", v}
	services.ResponseWithJson(w, http.StatusOK, response) //回傳
}

// @Summary Delete todo from List
// @Tags Todo
// @Accept json
// @Produce json
// @param UserID path string true "UserId"
// @param ItemID path int true "ItemId"
// @Success 200 string string 成功後返回的值
// @Router /api/todo/{UserID}/{ItemID} [Delete]
func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryId := vars["ItemID"] //獲取url參數
	Firebase, _ := Fb.Ref("/" + vars["UserID"] + "/" + queryId)
	err := Firebase.Remove()
	if err != nil {
		log.Fatal(err)
	}
	response := ApiResponse{"200", TodoList}
	services.ResponseWithJson(w, http.StatusOK, response)
}

// @Summary Update todo from List
// @Tags Todo
// @Accept json
// @Produce json
// @param Body body Todo true "body"
// @Success 200 {array} ApiResponse
// @Router /api/todo [Patch]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024)) //io.LimitReader限制大小
	if err != nil {
		fmt.Println(err)
	}
	var updateTodo Todo
	_ = json.Unmarshal(body, &updateTodo) //轉為json
	v := map[string]interface{}{
		"Item":  updateTodo.Item,
		"State": updateTodo.State,
	}
	Firebase, _ := Fb.Ref("/" + updateTodo.UserID + "/" + strconv.FormatInt(int64(updateTodo.ItemID), 10))
	err = Firebase.Update(v)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	response := ApiResponse{"200", TodoList}

	services.ResponseWithJson(w, http.StatusOK, response) //回傳
}

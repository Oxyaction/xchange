package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	pb "github.com/Oxyaction/xchange/rpc"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountController struct {
	Client pb.AccountClient
}

type appError struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"status"`
}

func (ac *AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	reply, err := ac.Client.Create(context.Background(), &pb.CreateRequest{})
	if err != nil {
		DisplayAppError(w, err, "Create account error", 500)
		return
	}
	j, err := json.Marshal(reply)
	if err != nil {
		DisplayAppError(w, err, "JSON encoding error", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Write the JSON data to the ResponseWriter
	w.Write(j)
}

func (ac *AccountController) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		DisplayAppError(w, err, "Entity not found", 404)
		return
	}

	reply, err := ac.Client.GetBalance(context.Background(), &pb.GetBalanceRequest{Id: id.String()})

	if err != nil {
		DisplayRPCError(w, err)
		return
	}
	j, err := json.Marshal(reply)

	if err != nil {
		DisplayAppError(w, err, "JSON encoding error", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Write the JSON data to the ResponseWriter
	w.Write(j)
}

func DisplayRPCError(w http.ResponseWriter, handlerError error) {
	errStatus, _ := status.FromError(handlerError)
	errObj := appError{
		Error:      true, // errStatus.Code().String(),
		Message:    errStatus.Message(),
		HTTPStatus: 500,
	}

	if codes.NotFound == errStatus.Code() {
		errObj.HTTPStatus = 404
	} else {
		fmt.Fprintf(os.Stderr, "[RPCError]: %s\n", handlerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errObj.HTTPStatus)
	if j, err := json.Marshal(errObj); err == nil {
		w.Write(j)
	}
}

func DisplayAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	errObj := appError{
		Error:      true, //handlerError.Error(),
		Message:    message,
		HTTPStatus: code,
	}
	fmt.Fprintf(os.Stderr, "[AppError]: %s\n", handlerError)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(errObj); err == nil {
		w.Write(j)
	}
}

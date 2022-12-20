package http

import (
	"tasks/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"strconv"

)

// PostTasks
// @ID create
// @tags tasks
// @Summary Create task
// @Description Create task if user is correctly logged in, returns task ID
// @Security BasicAuth
// @Param author body string true "Current user"
// @Success 201 {object} models.Task
// @Failure 400 {string} string "bad request"
// @Failure 403 {string} string "forbidden"
// @Router / [post]
func (a *Adapter) create(ctx *gin.Context) {
	parsedEmail := ctx.Request.Header["Login"][0]
	parsedRole := ctx.Request.Header["Role"][0]
	if parsedRole != "administrator" {
		ctx.Writer.WriteHeader(http.StatusForbidden)
	}
	//parsedEmail := "author"

	author := models.User{Email : parsedEmail}
	approve := models.Approve{Email : "myapprove"}
	approveList := []*models.Approve{&approve}
	
	taskId := a.tasks.CreateTask(ctx, author.Email, "mytitle", "mydescription", approveList)

	ctx.JSON(http.StatusOK, gin.H{"taskId": taskId})
}

// PostTasksApprove
// @ID approve
// @tags tasks
// @Summary Approve task
// @Description Set status "approved" by approver
// @Security BasicAuth
// @Param taskId path string true "Task ID"
// @Param approveInd path string true "Index of approver in approvers list"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Router /approve/{taskId}/{approveInd} [post]
func (a *Adapter) approve(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	if !strfmt.IsUUID(taskId) {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	task := a.tasks.GetTaskById(ctx, taskId)
	if task == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	
	approveInd, err := strconv.Atoi(ctx.Param("approveInd"))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	task.Approves[approveInd].Status = "approved"
	ctx.Writer.WriteHeader(http.StatusOK)
}

// PostTasksDecline
// @ID decline
// @tags tasks
// @Summary Decline task
// @Description Set status "declined" by approver
// @Security BasicAuth
// @Param taskId path string true "Task ID"
// @Param approveInd path string true "Index of approver in approvers list"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Router /decline/{taskId}/{approveInd} [post]
func (a *Adapter) decline(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	if !strfmt.IsUUID(taskId) {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	task := a.tasks.GetTaskById(ctx, taskId)
	if task == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	
	approveInd, err := strconv.Atoi(ctx.Param("approveInd"))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	task.Approves[approveInd].Status = "declined"
	ctx.Writer.WriteHeader(http.StatusOK)
}

// GetTasksTaskID
// @ID getTaskID
// @tags tasks
// @Summary Gets task description
// @Description Gets task description by Task Id
// @Security BasicAuth
// @Param taskId path string true "Task ID"
// @Success 200 {object} string
// @Failure 400 {string} string "bad request"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "task not found"
// @Router /{taskId} [get]
func (a *Adapter) getTaskID(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	if !strfmt.IsUUID(taskId) {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	task := a.tasks.GetTaskById(ctx, taskId)
	if task == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	ctx.JSON(http.StatusOK, gin.H{"Task Title": task.Title})
	//ctx.Writer.WriteHeader(http.StatusOK)
}

// GetTasks
// @ID getUserTasks
// @tags tasks
// @Summary Gets tasks list
// @Description Gets tasks list created by current user
// @Security BasicAuth
// @Success 200 {object} []string
// @Failure 400 {string} string "bad request"
// @Failure 403 {string} string "forbidden"
// @Router / [get]
func (a *Adapter) getUserTasks(ctx *gin.Context) {
	parsedEmail := ctx.Request.Header["Login"][0]
	parsedRole := ctx.Request.Header["Role"][0]
	if parsedRole != "administrator" {
		ctx.Writer.WriteHeader(http.StatusForbidden)
	}
	if parsedEmail == "" {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
	tasks := a.tasks.GetTasksByAuthor(ctx, parsedEmail)
	for _, task := range tasks {
		ctx.JSON(http.StatusOK, gin.H{"Task Title": task.Title})
	}
}

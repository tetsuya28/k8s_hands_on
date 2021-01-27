package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tetsuya28/k8s_hands_on/api/model"
	"gorm.io/gorm"
)

type dbConfig struct {
	Host     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Database string `required:"true"`
}

type simpleResponse struct {
	Message string `json:"message"`
}

type privateAPI struct {
	DB *gorm.DB
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	dbConf := dbConfig{}
	err := envconfig.Process("DB", &dbConf)
	if err != nil {
		panic(err)
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Database)
	db, err := dbClient(dns)
	if err != nil {
		panic(err)
	}
	privateInterface := privateAPI{}
	privateInterface.DB = db

	// 面倒くさいのでここでMigrate
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		panic(err)
	}

	e := echo.New()

	// CORS対策
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthz", healthz)

	api := e.Group("/api")
	api.GET("/todos", privateInterface.fetchAllTodos)
	api.POST("/todo", privateInterface.postTodo)
	api.POST("/todo/:id/done", privateInterface.doneTodo)
	api.DELETE("/todo/:id", privateInterface.deleteTodo)

	e.Logger.Fatal(e.Start(":8080"))
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, simpleResponse{Message: "Health check endpoint"})
}

// TODO: (db privateAPI) ← 命名バグってるけどセンス無いので誰か助けて
func (db privateAPI) postTodo(c echo.Context) error {
	todo := model.Todo{}
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, simpleResponse{Message: "Bad Request"})
	}
	if err := c.Validate(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, simpleResponse{Message: "Bad Request(Validation error)"})
	}

	err := db.DB.Create(&todo).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, simpleResponse{Message: "DB Error"})
	}
	return nil
}

func (db privateAPI) fetchAllTodos(c echo.Context) error {
	todos := make([]model.Todo, 0)
	err := db.DB.Find(&todos).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, model.TodosResponse{})
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, todos)
}

func (db privateAPI) doneTodo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, simpleResponse{Message: "Bad Request(ID is missing)"})
	}

	todo := model.Todo{}
	err := db.DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, simpleResponse{Message: "DB Error"})
	}
	now := time.Now()
	todo.IsDone = !todo.IsDone
	todo.UpdatedAt = now
	err = db.DB.Save(&todo).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, simpleResponse{Message: "DB Error"})
	}
	return c.JSON(http.StatusOK, todo)
}

func (db privateAPI) deleteTodo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, simpleResponse{Message: "Bad Request(ID is missing)"})
	}
	err := db.DB.Delete(&model.Todo{}, id).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, simpleResponse{Message: "DB Error"})
	}
	return nil
}

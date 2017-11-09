package main

import (
  "database/sql"
  //"fmt"
  "net/http"
  "os"

  _ "github.com/go-sql-driver/mysql"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/JonathanMH/goClacks/echo"
)

type (
	groupname struct {
		Error string `json:"error"`
    	ID string `json:"ID"`
		Name string `json:"Naimenovanie"`
	}

	teachername struct {
		Error string `json:"error"`
    	ID string `json:"ID"`
		Name string `json:"FIO"`
	}

	shedule struct {
		Error string `json:"error"`
    	Date string `json:"date"`
		Class string `json:"group"`
    	TimeStart string `json:"start"`
		TimeStop string `json:"stop"`
    	Discipline string `json:"disc"`
		Tip string `json:"type"`
    	Teacher string `json:"teacher"`
		Cabinet string `json:"kabinet"`
		Subgroup string `json:"subgr"`
	}
)


func main() {
	// Echo instance
	e := echo.New()
	e.Use(goClacks.Terrify)
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.File("html/index.html")
	})

	//Все группы
	e.GET("/groups", func(c echo.Context) error {
		db, _ := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
				
		var Naimenovanie string;

		rows, _ := db.Query("SELECT class FROM timetable  where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by class")
			
		groups := make([]groupname, 10)
		for rows.Next(){
			_ = rows.Scan(&Naimenovanie)
			groups = append(groups, groupname{"", "null", Naimenovanie})
		}
	
		return c.JSON(http.StatusOK, groups)//вернуть json				
	})

	//Все преподаватели
	e.GET("/teachers", func(c echo.Context) error {
		db, _ := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
							
		var Naimenovanie string;
			
		rows, _ := db.Query("SELECT teacher FROM timetable where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by teacher")
						
		teachers := make([]teachername, 10)
		for rows.Next(){
			_ = rows.Scan(&Naimenovanie)
			teachers = append(teachers, teachername{"", "null", Naimenovanie})
		}
				
		return c.JSON(http.StatusOK, teachers)//вернуть json				
	})

	//Расписание для студента
	e.GET("/shedule/student/:group/today", func(c echo.Context) error {
		db, _ := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
		group := c.Param("group")
							
		var date string;
		var class string;
		var timeStart string;
		var timeStop string;
		var discipline string;
		var tip string;
		var teacher string;
		var cabinet string;
		var subgroup string;
			
		rows, _ := db.Query("SELECT `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?)and(date = CURDATE())", group)
				
		pairs := make([]shedule, 0)
		for rows.Next(){
			_ = rows.Scan(&date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairs = append(pairs, shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
			
		}
					
		return c.JSON(http.StatusOK, pairs)//вернуть json				
	})

	//Расписание для преподавателя
	e.GET("/shedule/teacher/:teacher/today", func(c echo.Context) error {
		db, _ := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
		prep := c.Param("teacher")
							
		var date string;
		var class string;
		var timeStart string;
		var timeStop string;
		var discipline string;
		var tip string;
		var teacher string;
		var cabinet string;
		var subgroup string;
			
		rows, _ := db.Query("SELECT `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (teacher LIKE (?))and(date = CURDATE())", prep+"%")
				
		pairs := make([]shedule, 0)
		for rows.Next(){
			_ = rows.Scan(&date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairs = append(pairs, shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
			
		}
						
		return c.JSON(http.StatusOK, pairs)//вернуть json				
	})

	/*
	//Название группы по id
	e.GET("/groups/name/:id", func(c echo.Context) error {
		requestedID := c.Param("id") //вытащить id из запроса
		db, err := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
		if err != nil { //в случае ошибки
			fmt.Println(err.Error())
			response := groupname{ID: "", Error: "true", Name: ""}
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer db.Close() //В случае ошибки (?)
		var Naimenovanie string;
		var ID string;
		err = db.QueryRow("SELECT ID, Naimenovanie FROM groups_original WHERE ID = ?", requestedID).Scan(&ID, &Naimenovanie) //Запрос, вернет ошибку, если не удалось просканировать
		if err != nil {
			fmt.Println(err)
		}
		response := groupname{ID: ID, Error: "false", Name: Naimenovanie} //Создание нового json объекта
		return c.JSON(http.StatusOK, response)//вернуть json		
	})
	*/




	
	//e.Logger.Fatal(e.Start(":4000"))
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
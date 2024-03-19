package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type RewsData struct {
	Quantity  int
	Price     int
	AdLink    string
	Keyword   string
	PromoCode string
}

var store = sessions.NewCookieStore([]byte("scscwcwcwdcw234930f93o"))

type Renderer struct {
	templates *template.Template
}

// Render реализует метод Render интерфейса echo.Renderer
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	// Set up HTML template rendering
	renderer := &Renderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "static")
	e.GET("/plus11", func(c echo.Context) error {
		return c.File("teml/plus11.html")
	})
	e.GET("/plus2", func(c echo.Context) error {
		return c.File("teml/plus2.html")
	})
	e.GET("/plus3", func(c echo.Context) error {
		return c.File("teml/plus3.html")
	})
	e.GET("/plus4", func(c echo.Context) error {
		return c.File("teml/plus4.html")
	})

	e.GET("/min1", func(c echo.Context) error {
		return c.File("teml/min1.html")
	})
	e.GET("/min2", func(c echo.Context) error {
		return c.File("teml/min2.html")
	})
	e.GET("/min3", func(c echo.Context) error {
		return c.File("teml/min3.html")
	})
	e.GET("/min4", func(c echo.Context) error {
		return c.File("teml/min4.html")
	})
	e.GET("/main", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.GET("/", func(c echo.Context) error {
		t, err := template.ParseFiles("main.html")
		if err != nil {
			e.Logger.Error(err)
			return err
		}
		return t.Execute(c.Response(), nil)
	})

	e.GET("/register", func(c echo.Context) error {
		return c.File("reg.html")
	})
	e.GET("/client", func(c echo.Context) error {
		// Здесь можете добавить код обработки GET-запроса на /client.html
		session, err := store.Get(c.Request(), "session-name")
		if err != nil {
			e.Logger.Error(err)
			return err
		}
		data := map[string]interface{}{
			"Name":  session.Values["Name"],
			"Email": session.Values["Email"],
		}
		return c.Render(http.StatusOK, "client.html", data)
	})
	e.POST("/login", login)
	e.POST("/register", func(c echo.Context) error {
		user := new(User)

		// Получаем данные формы
		name := c.FormValue("name")
		email := c.FormValue("email")
		password := c.FormValue("password")
		if len(password) == 0 {
			return c.File("fatal.html")
		}
		if len(name) == 0 {
			return c.File("fatal.html")
		}
		if len(email) == 0 {
			return c.File("fatal.html")
		}

		hashedPassword, err := hashPassword(password)
		if err != nil {
			e.Logger.Error()
			return err
		}

		// Заполняем структуру User
		user.Name = name
		user.Email = email
		user.Password = hashedPassword

		// Сохраняем данные в PostgreSQL
		err = saveToPostgres(user)
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		session, err := store.Get(c.Request(), "session-name")
		if err != nil {
			e.Logger.Error(err)
			return err
		}
		session.Values["Name"] = user.Name
		session.Values["Email"] = user.Email
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		return c.File("index.html")
	})
	e.POST("/buyrew1", SaveRew1)
	e.POST("/buyrew2", SaveRew2)
	e.POST("/buyrew3", SaveRew3)
	e.POST("/buyrew4", SaveRew4)
	e.POST("/buyrew5", SaveRew5)
	e.POST("/buyrew6", SaveRew6)
	e.POST("/buyrew7", SaveRew7)
	e.POST("/buyrew8", SaveRew8)

	e.GET("/login", func(c echo.Context) error {
		return c.File("kar.html")
	})
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}

func saveToPostgres(user *User) error {
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Создание таблицы, если она не существует
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT, password TEXT)")
	if err != nil {
		return err
	}

	// Вставка данных в таблицу
	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func login(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if len(password) == 0 {
		return c.File("fatal.html")
	}
	if len(name) == 0 {
		return c.File("fatal.html")
	}
	if len(email) == 0 {
		return c.File("fatal.html")
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database connection error"})
	}
	defer db.Close()
	var user User
	err = db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(401, "Неверные учётные данные")
		}
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return echo.NewHTTPError(401, "Неверные данные")
	}

	user.Name = name
	user.Email = email
	session, err := store.Get(c.Request(), "session-name")
	if err != nil {
		log.Println(err)
		return err
	}
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		log.Println(err)
		return err
	}

	return c.File("index.html")
}

func SaveRew1(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)
	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/plif_1.html")
}

func SaveRew2(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/plif_2.html")
}

func SaveRew3(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/plif_3.html")
}

func SaveRew4(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/plif_4.html")
}

func SaveRew5(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/milfi.html")
}

func SaveRew6(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/milfi2.html")
}

func SaveRew7(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/milfi3.html")
}

func SaveRew8(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=shop sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer db.Close()

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	if err != nil {
		fmt.Println("Error conection")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		fmt.Println("Error conrction")
		return c.String(http.StatusBadRequest, "Invalid value")
	}
	adLink := c.FormValue("adLink")
	keyword := c.FormValue("keyword")
	promocode := c.FormValue("promoCode") // Corrected variable name

	if len(adLink) == 0 {
		return c.File("fatal.html")
	}
	if len(keyword) == 0 {
		return c.File("fatal.html")
	}
	if len(promocode) == 0 {
		return c.File("fatal.html")
	}
	reviewData := RewsData{
		Quantity:  quantity,
		Price:     price,
		AdLink:    adLink,
		Keyword:   keyword,
		PromoCode: promocode,
	}

	_, err = db.Exec("INSERT INTO orders(quantity, price, ad_link, keyword, promocode) VALUES($1,$2, $3, $4, $5)", reviewData.Quantity, reviewData.Price, reviewData.AdLink, reviewData.Keyword, reviewData.PromoCode)

	if err != nil {
		fmt.Println("Error inserting data into the database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.File("teml/milfi4.html")
}

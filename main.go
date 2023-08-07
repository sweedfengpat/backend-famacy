package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/pborman/uuid"
)

// user ------------------ start

type AddProduct_struct struct {
	productType        string `json: "productType"`
	productCode        string `json: "productCode"`
	productName        string `json: "productName"`
	productDescription string `json: "productDescription"`
	productPrice       string `json: "productPrice"`
	productAmount      string `json: "productAmount"`
}

type AddCart struct {
	productID         string `json: "productID"`
	productName       string `json: "productName"`
	productAmount     string `json: "productAmount"`
	productPrice      string `json: "productPrice"`
	productTotal      string `json: "productTotal"`
	productImage      string `json: "productImage"`
	email             string `json: "email"`
	productAmountData string `json: "productAmountData"`
}

type Product struct {
	id                 int    `json: "id"`
	productType        string `json: "productType"`
	productCode        string `json: "productCode"`
	productName        string `json: "productName"`
	productDescription string `json: "productDescription"`
	productImage       string `json: "productImage"`
	productPrice       int    `json: "productPrice"`
	productAmount      int    `json: "productAmount"`
}
type Register struct {
	FirstName   string `json: "firstname"`
	LastName    string `json: "lastname"`
	Email       string `json: "email"`
	Password    string `json: "password"`
	Token       string `json: "token"`
	Tel         string `json: "tel"`
	Personal_id string `json: "personal_id"`
	Brithday    string `json: "brithday"`
	Gender      string `json: "gender"`
}

type UpdateProfile struct {
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
	Email     string `json: "email"`
	Password  string `json: "password"`
	Address   string `json: "address"`
	Token     string `json: "token"`
}

type Login struct {
	Email     string `json: "email"`
	Password  string `json: "password"`
	Token     string `json: "token"`
	Address   string `json: "address"`
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
}

type Order struct {
	listName   string `json: "listName"`
	listPrice  string `json: "listPrice"`
	listAmount string `json: "listAmount"`
	listImage  string `json: "listImage"`
	listId     string `json: "listId"`
	email      string `json: "email"`
	address    string `json: "address"`
	total      string `json: "total"`
	day        string `json: "day"`
	slipImage  string `json: "slipImage"`
}

type getOrder struct {
	id         int    `json: "id"`
	code       string `json: "code"`
	listName   string `json: "listName"`
	listPrice  string `json: "listPrice"`
	listAmount string `json: "listAmount"`
	listImage  string `json: "listImage"`
	listId     string `json: "listId"`
	email      string `json: "email"`
	address    string `json: "address"`
	total      string `json: "total"`
	day        string `json: "day"`
	slipImage  string `json: "slipImage"`
	status     string `json: "status"`
}

type forgotInsert struct {
	email    string `json: "email"`
	password string `json: "password"`
}

type about struct {
	logo    string `json: "logo"`
	name    string `json: "name"`
	phone   string `json: "phone"`
	address string `json: "address"`
}

type categoryGet struct {
	id   int    `json: "id"`
	name string `json: "name"`
}

type addressGet struct {
	id      int    `json: "id"`
	name    string `json: "name"`
	address string `json: "address"`
}

type addressAdd struct {
	Name     string `json: "name"`
	Address  string `json: "address"`
	Email    string `json: "email"`
	Password string `json: "password"`
}

func postRegister(context *gin.Context) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gostoredb")
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	defer db.Close()

	var newRegister Register
	if err := context.BindJSON(&newRegister); err != nil {
		return
	}
	fmt.Println("newRegister: ", newRegister)
	insert, err := db.Query("INSERT INTO users (firstName, lastName, email, password, token, Tel, personal_id, brithday, gender) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", newRegister.FirstName, newRegister.LastName, newRegister.Email, newRegister.Password, newRegister.Token, newRegister.Tel, newRegister.Personal_id, newRegister.Brithday, newRegister.Gender)
	if err != nil {
		fmt.Println("Error during insertion:", err)
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		return
	}
	defer insert.Close()
	fmt.Println("values added!")

	context.IndentedJSON(http.StatusCreated, gin.H{
		"code": 200,
	})
}

var conn = "root:root@tcp(127.0.0.1:3306)/gostoredb"

func postUpdateProfile(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var newUpdateProfile UpdateProfile
	if err := context.BindJSON(&newUpdateProfile); err != nil {
		return
	}
	update, err := db.Query("UPDATE users SET firstName=?, lastName=?, email=?, password=?, address=? WHERE email = ?", newUpdateProfile.FirstName, newUpdateProfile.LastName, newUpdateProfile.Email, newUpdateProfile.Password, newUpdateProfile.Address, newUpdateProfile.Email)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

	context.IndentedJSON(http.StatusCreated, gin.H{
		"code": 200,
	})
}

func postLogin(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the request body
	var user Login
	var newLogin Login
	if err := context.BindJSON(&newLogin); err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		return
	}

	// Query the database
	err = db.QueryRow("SELECT email, password, token FROM users where email = ? AND password = ?", newLogin.Email, newLogin.Password).Scan(&user.Email, &user.Password, &user.Token)
	if err != nil {
		// Return an error if the user didn't exist
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		return
	} else {
		// Return a success code if the user was found
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
		return
	}
}

func postProfile(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var user UpdateProfile
	var newLogin Login
	if err := context.BindJSON(&newLogin); err != nil {
		return
	}
	err = db.QueryRow("SELECT firstName, lastName, email, password, address, token FROM users WHERE email = ? AND password = ?", newLogin.Email, newLogin.Password).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Address, &user.Token)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, user)
	}
}

func CheckLogin(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var user Login
	var newLogin Login
	if err := context.BindJSON(&newLogin); err != nil {
		return
	}

	err = db.QueryRow("SELECT email, password, token, address, firstName, lastName FROM users WHERE email = ? AND password = ?", newLogin.Email, newLogin.Password).Scan(&user.Email, &user.Password, &user.Token, &user.Address, &user.FirstName, &user.LastName)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		fmt.Println(user)
		context.IndentedJSON(http.StatusCreated, user)
	}
}

func AddProduct(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	uuidWithHyphen := uuid.NewRandom()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	file, header, err := context.Request.FormFile("upload")
	filename := uuid + header.Filename
	fmt.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename)
	var product AddProduct_struct
	product.productType = context.Request.FormValue("productType")
	product.productCode = context.Request.FormValue("productCode")
	product.productName = context.Request.FormValue("productName")
	product.productDescription = context.Request.FormValue("productDescription")
	product.productPrice = context.Request.FormValue("productPrice")
	product.productAmount = context.Request.FormValue("productAmount")

	fmt.Println(product)
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
	}

	insert, err := db.Query("INSERT INTO product (productType, productCode, productName, productDescription, productImage, productPrice, productAmount) VALUES (?, ?, ?, ?, ?, ?, ?)", product.productType, product.productCode, product.productName, product.productDescription, filename, product.productPrice, product.productAmount)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		fmt.Println(err.Error())
	}
	defer insert.Close()
	fmt.Println("values added!")

	context.IndentedJSON(http.StatusCreated, gin.H{
		"code": 200,
	})

}

func AddCategory(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO category (name) VALUES (?)", context.Request.FormValue("name"))
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		fmt.Println(err.Error())
	}
	defer insert.Close()
	fmt.Println("values added!")

	context.IndentedJSON(http.StatusCreated, gin.H{
		"code": 200,
	})

}

// AddAddress Adds a new address to the database
func AddAddress(context *gin.Context) {
	// Establish a connection to the database
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var addressAdd addressAdd
	if err := context.BindJSON(&addressAdd); err != nil {
		return
	}
	fmt.Println("addressAdd: ", addressAdd)

	// Insert a new row into the database
	insert, err := db.Query("INSERT INTO address (name, address, email, password) VALUES (?, ?, ?, ?)", addressAdd.Name, addressAdd.Address, addressAdd.Email, addressAdd.Password)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":    500,
			"message": "INSERT Address fail ",
		})
		fmt.Println(err.Error())
	}
	defer insert.Close()
	fmt.Println("values added!")

	//get address id
	var addressGet addressGet
	err = db.QueryRow("SELECT id, name, address FROM address WHERE email = ? AND password = ?", addressAdd.Email, addressAdd.Password).Scan(&addressGet.id, &addressGet.name, &addressGet.address)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":    500,
			"message": "SELECT Address fail ",
		})
		fmt.Println(err.Error())
	}

	//update address a new row into the database User
	update, err := db.Query("UPDATE users SET address=? WHERE email = ?", addressGet.id, addressAdd.Email)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":    500,
			"message": "UPDATE Address fail ",
		})
		fmt.Println(err.Error())
	}
	defer update.Close()
	fmt.Println("values added!")

	// Return a success message
	context.IndentedJSON(http.StatusCreated, gin.H{
		"code":    200,
		"id":      addressGet.id,
		"name":    addressGet.name,
		"address": addressGet.address,
	})

}

func UpdateProduct(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var product AddProduct_struct
	product.productType = context.Request.FormValue("productType")
	product.productCode = context.Request.FormValue("productCode")
	product.productName = context.Request.FormValue("productName")
	product.productDescription = context.Request.FormValue("productDescription")
	product.productPrice = context.Request.FormValue("productPrice")
	product.productAmount = context.Request.FormValue("productAmount")
	var id = context.Request.FormValue("id")

	insert, err := db.Query("UPDATE product SET productType=?, productCode=?, productName=?, productDescription=?, productPrice=?, productAmount=? WHERE id=?", product.productType, product.productCode, product.productName, product.productDescription, product.productPrice, product.productAmount, id)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		fmt.Println(err.Error())
	}
	defer insert.Close()
	fmt.Println("values added!")

	update, err := db.Query("UPDATE cart SET productAmountData=? WHERE productID=?", product.productAmount, id)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
		fmt.Println(err.Error())
	}
	defer update.Close()
	fmt.Println("values added!")

	context.IndentedJSON(http.StatusCreated, gin.H{
		"code": 200,
	})

}

func AllProduct(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))
}

func AllCategory(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM category")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func AllAddress(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM address WHERE email=? AND password=?", context.Request.FormValue("email"), context.Request.FormValue("password"))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func FilterProduct(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	var productType = context.Request.FormValue("productType")

	rows, err := db.Query("SELECT * FROM product WHERE productType=?", productType)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func Cart(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM cart WHERE email=?", context.Request.FormValue("email"))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func getProductQuery(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var id = context.Request.FormValue("id")
	var product Product
	err = db.QueryRow("SELECT id, productType, productCode, productName, productDescription, productImage, productPrice, productAmount FROM product WHERE id = ?", id).Scan(&product.id, &product.productType, &product.productCode, &product.productName, &product.productDescription, &product.productImage, &product.productPrice, &product.productAmount)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":               200,
			"id":                 product.id,
			"productType":        product.productType,
			"productCode":        product.productCode,
			"productName":        product.productName,
			"productDescription": product.productDescription,
			"productImage":       product.productImage,
			"productPrice":       product.productPrice,
			"productAmount":      product.productAmount,
		})
	}

}

func addCart(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var newCart AddCart
	newCart.productID = context.Request.FormValue("productID")
	newCart.productName = context.Request.FormValue("productName")
	newCart.productAmount = context.Request.FormValue("productAmount")
	newCart.productAmountData = context.Request.FormValue("productAmountData")
	newCart.productPrice = context.Request.FormValue("productPrice")
	newCart.productTotal = context.Request.FormValue("productTotal")
	newCart.productImage = context.Request.FormValue("productImage")
	newCart.email = context.Request.FormValue("email")
	insert, err := db.Query("INSERT INTO cart (productID, productName, productAmount, productPrice, productTotal, productImage, productAmountData, email) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", newCart.productID, newCart.productName, newCart.productAmount, newCart.productPrice, newCart.productTotal, newCart.productImage, newCart.productAmountData, newCart.email)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer insert.Close()
	fmt.Println("values added!")

}

func updateAmount(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var email = context.Request.FormValue("email")
	var productTotal = context.Request.FormValue("productTotal")
	var productAmount = context.Request.FormValue("productAmount")

	insert, err := db.Query("UPDATE cart SET productAmount = ?, productTotal = ? WHERE id = ? AND email = ?", productAmount, productTotal, id, email)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer insert.Close()
	fmt.Println("values added!")

}

func deleteCart(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var email = context.Request.FormValue("email")

	update, err := db.Query("DELETE FROM cart WHERE id = ? AND email = ?", id, email)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func AddOrder(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	uuidWithHyphen := uuid.NewRandom()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	file, header, err := context.Request.FormFile("productImage")
	defer file.Close()
	filename := uuid + header.Filename
	out, err := os.Create("./tmp/" + filename)

	var order Order
	order.listName = context.Request.FormValue("listName")
	order.listPrice = context.Request.FormValue("listPrice")
	order.listAmount = context.Request.FormValue("listAmount")
	order.listImage = context.Request.FormValue("listImage")
	order.listId = context.Request.FormValue("listId")
	order.email = context.Request.FormValue("email")
	order.address = context.Request.FormValue("address")
	order.total = context.Request.FormValue("total")
	order.day = context.Request.FormValue("day")
	order.slipImage = filename

	insert, err := db.Query("INSERT INTO `order` (`code`, `listName`, `listPrice`, `listAmount`, `listImage`, `listId`, `email`, `address`, `total`, `day`, `slipImage`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", uuid, order.listName, order.listPrice, order.listAmount, order.listImage, order.listId, order.email, order.address, order.total, order.day, order.slipImage)
	if err != nil {
		fmt.Println(err.Error())

	}
	defer insert.Close()
	fmt.Println("values added!")

	fmt.Println(order)
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
	}

	update, err := db.Query("DELETE FROM cart WHERE email = ?", order.email)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("values added!")
	}

	defer update.Close()

	rows, err := db.Query("SELECT * FROM `users` WHERE `token`='admin'")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(jsonData))

	// for loop jsonData for email
	for _, item := range tableData {
		fmt.Println()

		email := item["email"].(string)

		hostname := "smtp.gmail.com"
		auth := smtp.PlainAuth("", "starterkidpiyaman1717@gmail.com", "siwpztupohujvlwq", hostname)

		msg := "From: " + "shop" + "\n" +
			"To: " + email + "\n" +
			"Subject: คำสั่งซื้อสินค้า\n\n" +
			"รหัสสินค้า: " + order.listId + "\n" +
			"ชื่อสินค้า: " + order.listName + "\n" +
			"ราคาสินค้า: " + order.listPrice + "\n" +
			"จำนวนสินค้า: " + order.listAmount + "\n" +
			"ราคารวม: " + order.total + "\n" +
			"วันที่สั่งซื้อ: " + order.day + "\n" +
			"ที่อยู่จัดส่ง: " + order.address + "\n" +
			"รอการตรวจสอบจากผู้ดูแลระบบ"

		err := smtp.SendMail(hostname+":587", auth, "from@example.com", []string{email},
			[]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Email Sent!")
	}
}

func AllOrder(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `order` WHERE `status` = 'รอตรวจสอบ'")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func UpdateOrder(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var status = context.Request.FormValue("status")

	insert, err := db.Query("UPDATE `order` SET `status` = ? WHERE `id` = ?", status, id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"status": "success",
			"code":   200,
		})
	}
	defer insert.Close()
	fmt.Println("values added!")

}

func UpdateOrderPaySuccess(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `order` WHERE NOT `status`='รอตรวจสอบ' AND NOT `status`='ล้มเหลว' AND NOT `status`='เรียบร้อย'")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func getOrderQuery(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var id = context.Request.FormValue("id")
	var order getOrder
	err = db.QueryRow("SELECT `id`, `code`, `listName`, `listPrice`, `listAmount`, `listImage`, `listId`, `email`, `address`, `total`, `day`, `slipImage`, `status` FROM `order` WHERE `id` = ?", id).Scan(&order.id, &order.code, &order.listName, &order.listPrice, &order.listAmount, &order.listImage, &order.listId, &order.email, &order.address, &order.total, &order.day, &order.slipImage, &order.status)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"id":         order.id,
			"code":       order.code,
			"listName":   order.listName,
			"listPrice":  order.listPrice,
			"listAmount": order.listAmount,
			"listImage":  order.listImage,
			"listId":     order.listId,
			"email":      order.email,
			"address":    order.address,
			"total":      order.total,
			"day":        order.day,
			"slipImage":  order.slipImage,
			"status":     order.status,
		})
	}
	fmt.Println(order)
}

func UpdateOrderPayment(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `order` WHERE NOT `status`='รอตรวจสอบ' AND NOT `status`='ล้มเหลว'")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func AllOrderPayment(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `order` WHERE NOT `status`='รอตรวจสอบ' AND NOT `status`='ล้มเหลว'")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

func DeleteProduct(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")

	update, err := db.Query("DELETE FROM product WHERE id = ?", id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func History(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	var email = context.Request.FormValue("email")

	rows, err := db.Query("SELECT * FROM `order` WHERE email = ?", email)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))
}

func AllForgot(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `forgot`")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, string(jsonData))

}

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func ForgotPassword(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var email = context.Request.FormValue("email")

	var forgot forgotInsert
	err = db.QueryRow("SELECT email, password FROM users WHERE email = ?", email).Scan(&forgot.email, &forgot.password)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		// hostname is used by PlainAuth to validate the TLS certificate.
		hostname := "smtp.gmail.com"
		auth := smtp.PlainAuth("", "starterkidpiyaman1717@gmail.com", "siwpztupohujvlwq", hostname)

		password_uuid := uuid.NewRandom()
		password := password_uuid.String()

		msg := "From: " + "shop" + "\n" +
			"To: " + email + "\n" +
			"Subject: " + "Forgot Password" + "\n\n" +
			"Your password is " + password

		err := smtp.SendMail(hostname+":587", auth, "from@example.com", []string{email},
			[]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Email Sent!")
		insert, err := db.Query("UPDATE `users` SET `password` = ? WHERE `email` = ?", password, email)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer insert.Close()

	}
}

func ChangeForgot(context *gin.Context) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var email = context.Request.FormValue("email")
	var password = context.Request.FormValue("password")

	insert, err := db.Query("UPDATE `users` SET `password` = ? WHERE `email` = ?", password, email)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer insert.Close()

	insert1, err := db.Query("UPDATE `forgot` SET `password` = ? WHERE `id` = ?", password, id)
	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"status": "success",
			"code":   200,
		})
	}
	defer insert1.Close()
	fmt.Println("values added!")

}

func DeleteForgot(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")

	update, err := db.Query("DELETE FROM forgot WHERE id = ?", id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func DeleteCategory(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")

	update, err := db.Query("DELETE FROM category WHERE id = ?", id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func DeleteAddress(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")

	update, err := db.Query("DELETE FROM address WHERE id = ?", id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func GetCategory(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")

	var category categoryGet
	err = db.QueryRow("SELECT id, name FROM category WHERE id = ?", id).Scan(&category.id, &category.name)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		fmt.Print(category)
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
			"id":   category.id,
			"name": category.name,
		})
	}
}

func GetAddress(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var email = context.Request.FormValue("email")
	var password = context.Request.FormValue("password")

	var address addressGet
	err = db.QueryRow("SELECT id, name, address FROM address WHERE id = ? AND email = ? AND password = ?", id, email, password).Scan(&address.id, &address.name, &address.address)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		fmt.Print(address)
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":    200,
			"id":      address.id,
			"name":    address.name,
			"address": address.address,
		})
	}
}

func UpdateCategory(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var name = context.Request.FormValue("name")

	update, err := db.Query("UPDATE `category` SET `name` = ? WHERE `id` = ?", name, id)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func UpdateAddress(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id = context.Request.FormValue("id")
	var name = context.Request.FormValue("name")
	var address = context.Request.FormValue("address")
	var email = context.Request.FormValue("email")
	var password = context.Request.FormValue("password")

	update, err := db.Query("UPDATE `address` SET `name` = ?, `address` = ? WHERE `id` = ? AND `email` = ? AND `password` = ? ", name, address, id, email, password)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func About(context *gin.Context) {

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM about`
	var about about
	err = db.QueryRow(sqlStatement).Scan(&about.logo, &about.name, &about.phone, &about.address)

	if err != nil {
		panic(err.Error())
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code":    200,
			"logo":    about.logo,
			"name":    about.name,
			"phone":   about.phone,
			"address": about.address,
		})
	}

}

func AboutUpdate(context *gin.Context) {

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var name = context.Request.FormValue("name")
	var phone = context.Request.FormValue("phone")
	var address = context.Request.FormValue("address")

	update, err := db.Query("UPDATE `about` SET `name` = ?, `phone` = ?, `address` = ?", name, phone, address)

	if err != nil {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 200,
		})
	}
	defer update.Close()
	fmt.Println("values added!")

}

func AddAboutImage(context *gin.Context) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// select sql and delete file
	var sqlSelect = `SELECT logo FROM about`
	var logo string
	err = db.QueryRow(sqlSelect).Scan(&logo)
	if err != nil {
		panic(err.Error())
	}
	os.Remove("./tmp/" + logo)

	uuidWithHyphen := uuid.NewRandom()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	file, header, err := context.Request.FormFile("upload")
	filename := uuid + header.Filename
	fmt.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename)

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
	}

	// update sql
	var sqlUpdate = `UPDATE about SET logo = ?`
	update, err := db.Query(sqlUpdate, filename)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
	fmt.Println("values added!")

}

func Track(context *gin.Context) {

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var track = context.Request.FormValue("track")
	var id = context.Request.FormValue("id")
	var description = context.Request.FormValue("description")

	fmt.Println(track + " " + id + " " + description)
	update, err := db.Query("UPDATE `order` SET `track` = ? WHERE `id` = ?", track, id)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()

	var order getOrder

	err = db.QueryRow("SELECT `id`, `address`, `email`, `listName`, `listPrice`, `listAmount`, `total`, `day` FROM `order` WHERE `id` = ?", id).Scan(&order.id, &order.address, &order.email, &order.listName, &order.listPrice, &order.listAmount, &order.total, &order.day)
	if err != nil {
		fmt.Println(err.Error())
		context.IndentedJSON(http.StatusCreated, gin.H{
			"code": 500,
		})
	} else {
		// hostname is used by PlainAuth to validate the TLS certificate.
		hostname := "smtp.gmail.com"
		auth := smtp.PlainAuth("", "starterkidpiyaman1717@gmail.com", "siwpztupohujvlwq", hostname)

		msg := "From: " + "shop" + "\n" +
			"To: " + order.email + "\n" +
			"Subject: " + " ข้อมูลการขนส่ง" + "\n\n" +
			"เลขพัสดุของคุณคือ " + track + "\n" +
			"รายละเอียดเพิ่มเติม " + description + "\n" +
			"รายการสั่งซื้อ " + order.listName + "\n" +
			"ราคา " + order.listPrice + "\n" +
			"จำนวน " + order.listAmount + "\n" +
			"รวม " + order.total + "\n" +
			"วันที่สั่งซื้อ " + order.day + "\n" +
			"ที่อยู่จัดส่ง " + order.address + "\n" +
			"ขอบคุณที่ใช้บริการ"

		err := smtp.SendMail(hostname+":587", auth, "from@example.com", []string{order.email},
			[]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Email Sent!")

	}

}

// user ------------------ end

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	r.Static("/", "./tmp")
	r.POST("/api/about", About)
	r.POST("/api/add-about-image", AddAboutImage)
	r.POST("/api/update-about", AboutUpdate)
	r.POST("/api/sign-up", postRegister)
	r.POST("/api/sign-in", postLogin)
	r.POST("/api/check-login", CheckLogin)
	r.POST("/api/profile", postProfile)
	r.POST("/api/update-profile", postUpdateProfile)
	r.POST("/api/add-product", AddProduct)
	r.POST("/api/update-product", UpdateProduct)
	r.POST("/api/get-product", getProductQuery)
	r.POST("/api/all-product", AllProduct)
	r.POST("/api/all-category", AllCategory)
	r.POST("/api/all-address", AllAddress)
	r.POST("/api/add-category", AddCategory)
	r.POST("/api/add-address", AddAddress)
	r.POST("/api/get-address", GetAddress)
	r.POST("/api/update-category", UpdateCategory)
	r.POST("/api/delete-category", DeleteCategory)
	r.POST("/api/get-category", GetCategory)
	r.POST("/api/get-cart", Cart)
	r.POST("/api/add-cart", addCart)
	r.POST("/api/update-amount", updateAmount)
	r.POST("/api/delete-address", DeleteAddress)
	r.POST("/api/delete-cart", deleteCart)
	r.POST("/api/add-order", AddOrder)
	r.POST("/api/all-order", AllOrder)
	r.POST("/api/update-order", UpdateOrder)
	r.POST("/api/update-address", UpdateAddress)
	r.POST("/api/all-order-pay-success", UpdateOrderPaySuccess)
	r.POST("/api/get-order", getOrderQuery)
	r.POST("/api/all-order-payment", AllOrderPayment)
	r.POST("/api/filter-product", FilterProduct)
	r.POST("/api/delete-product", DeleteProduct)
	r.POST("/api/history", History)
	r.POST("/api/all-forgot", AllForgot)
	r.POST("/api/forgot-password", ForgotPassword)
	r.POST("/api/change-forgot", ChangeForgot)
	r.POST("/api/delete-forgot", DeleteForgot)
	r.POST("/api/track", Track)

	r.Run("localhost:8080")

	fmt.Println("Hello, world.")
}

package db

import (
	"database/sql"
    "fmt"
    "log"
    "errors"
    
    _ "github.com/lib/pq"

    "crypto/sha1"
    "encoding/base64"
    "7linternational.com/gql-server/graph/model"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "XXXX" // Do it properly with a config file
  dbPassword = "XXXX"
  dbname   = "recipes"
)

func Register(firstName string, lastName string, email string, password string) (int32, error) {
	log.Println("DB: Registration attempt")

	db, err := getDB()

  if err != nil {
      return 500, err
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    return 503, err
  }

  var existingEmail string
  row := db.QueryRow("SELECT email FROM users WHERE email = $1 LIMIT 1", email)

  switch err = row.Scan(&existingEmail); err {
  case sql.ErrNoRows:
  	log.Println("DB: No user found with this email. We can register him.")
  	break
  case nil:
  	if existingEmail == email {
  		return 400, errors.New("User already exists")
  	}
  	// reduntant
  	break
  }

  hashedPassword:= hashPassword(password)
  log.Println("DB: Hash: " + hashedPassword)

  id := 0
  db.QueryRow("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)", firstName, lastName, email, hashedPassword).Scan(&id)

  log.Println("DB: New record with id: " + fmt.Sprintf("%d", id))


  return 200, nil
}

func Login(email string, password string) (*model.User, int32, error) {
	log.Println("DB: Login attempt")

	db, err := getDB()

  if err != nil {
    return &model.User{}, 500, err
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    return &model.User{}, 503, err
  }

  user := &model.User{}

  row := db.QueryRow("SELECT user_id, first_name, last_name, email FROM users WHERE email = $1 AND password = $2 LIMIT 1", email, hashPassword(password))

  switch err = row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err {
  case sql.ErrNoRows:
  	return &model.User{}, 404, errors.New("User not found")
  	// reduntant
  	break
  case nil:
  	log.Println("DB: User found, good to go.")
  	break
  default:
  	return &model.User{}, 500, err
  	// reduntant
  	break
  }

	return user, 200, nil
}

func GetUserRecipes(userId int64) (int64, error, []*model.Recipe) {
	log.Println("DB: Get recipe for user id: " + fmt.Sprintf("%d", userId))

	db, err := getDB()
  if err != nil {
      return 500, err, []*model.Recipe{}
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    return 503, err, []*model.Recipe{}
  }

  recipes := make([]*model.Recipe, 0)

  // select r.recipe_id, r.title, p.product_id, p.title FROM recipes r, products p, recipe_products rp WHERE r.user_id = 1 AND rp.recipe_id = r.recipe_id AND p.product_id = rp.product_id
  recipeRows, err := db.Query("SELECT recipe_id, title, description, image_url, user_id FROM recipes WHERE user_id = $1", userId)
  if err != nil {
  	return 500, err, []*model.Recipe{}
  }
  defer recipeRows.Close()

  for recipeRows.Next() {
  	recipe := &model.Recipe{}
  	err = recipeRows.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &recipe.ImageURL, &recipe.UserID)

  	log.Println("DB: Recipe Id: " + recipe.ID)
  	log.Println("DB: Recipe Title: " + recipe.Title)

  	recipeProductsRows, err := db.Query("SELECT product_id FROM recipe_products WHERE recipe_id = $1", recipe.ID)
  	defer recipeProductsRows.Close()

  	if err != nil {
  		log.Println("DB: Error in recipe products query: " + err.Error())
  	}

  	products := make([]*model.Product, 0)

  	for recipeProductsRows.Next() {
  		var productId int64
  		err = recipeProductsRows.Scan(&productId)

  		if err != nil {
	  		log.Println("DB: Error in recipe products scan: " + err.Error())
	  	}

  		log.Println("DB: productId: " + fmt.Sprintf("%d", productId))

  		productRows, err := db.Query("SELECT product_id, title, description, image_url FROM products WHERE product_id = $1", productId)
  		defer productRows.Close()

  		if err != nil {
	  		log.Println("DB: Error in product query: " + err.Error())
	  	}

  		for productRows.Next() {
  			product := &model.Product{}

  			err = productRows.Scan(&product.ID, &product.Title, &product.Description, &product.ImageURL)

  			if err != nil {
		  		log.Println("DB: Error in product scan: " + err.Error())
		  	}

  			log.Println("DB: Product Id: " + product.ID)
  			log.Println("DB: Product Title: " + product.Title)

  			products = append(products, product)
  		}
  	}

		recipe.Products = products

    recipes = append(recipes, recipe)
  }

	return 200, nil, recipes
}

func getDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, dbPassword, dbname)
	log.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	return db, err
}

func hashPassword(password string) string {
	bytePassword := []byte(password)

	hasher := sha1.New()
  hasher.Write(bytePassword)
  sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

  return sha
} 
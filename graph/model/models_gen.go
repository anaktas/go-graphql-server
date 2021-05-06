// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type Recipe struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageURL    string     `json:"imageUrl"`
	Products    []*Product `json:"products"`
	UserID      string     `json:"userId"`
}

type User struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     *string `json:"email"`
}

type UserRecipes struct {
	Hits []*Recipe `json:"hits"`
}

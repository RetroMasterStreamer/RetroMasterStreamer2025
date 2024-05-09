// repository/user_repository.go
package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"PortalCRG/internal/repository/entity"
)

// NewUserRepositoryMongo crea una nueva instancia de UserRepositoryMongo.
func NewUserRepositoryMongo() *UserRepositoryMongo {
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")

	db := &DataBase{}
	err := db.Connect(connectionString)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	return &UserRepositoryMongo{
		DataBase: db,
	}
}

// UserRepositoryMongo representa una implementación de UserRepository para MongoDB.
type UserRepositoryMongo struct {
	*DataBase
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (r *UserRepositoryMongo) AuthenticateUser(alias, password string) (*entity.User, error) {
	collection := r.client.Database("dbName").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"alias": alias, "password": password}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Otros métodos de UserRepositoryMongo...

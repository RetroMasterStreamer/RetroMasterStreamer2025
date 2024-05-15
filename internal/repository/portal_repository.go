// repository/user_repository.go
package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"PortalCRG/internal/repository/entity"
)

// NewPortalRepositoryMongo crea una nueva instancia de PortalRepositoryMongo.
func NewPortalRepositoryMongo() *PortalRepositoryMongo {
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")

	db := &DataBase{}
	err := db.Connect(connectionString)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	return &PortalRepositoryMongo{
		DataBase: db,
	}
}

// PortalRepositoryMongo representa una implementación de UserRepository para MongoDB.
type PortalRepositoryMongo struct {
	*DataBase
}

// Otros métodos de UserRepositoryMongo...
func (r *PortalRepositoryMongo) GetUserByAlias(userRef string) (*entity.User, error) {
	collection := r.client.Database("dbName").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"reference_text": userRef}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PortalRepositoryMongo) GetAllUsers() ([]*entity.User, error) {
	collection := r.client.Database("dbName").Collection("user")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*entity.User
	for cursor.Next(context.Background()) {
		var user entity.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

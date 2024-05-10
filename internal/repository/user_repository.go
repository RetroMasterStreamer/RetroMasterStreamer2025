// repository/user_repository.go
package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

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

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (r *UserRepositoryMongo) SetUserOnline(alias, sessionToken, hash string, online bool) (*entity.UserOnline, error) {
	collection := r.client.Database("dbName").Collection("userOnline")

	// Definir el filtro para encontrar el usuario
	filter := bson.M{"alias": alias, "sessionToken": sessionToken, "hash": hash}

	if online {

		// Definir los cambios que se van a realizar
		update := bson.M{"$set": bson.M{"online": online}}

		// Configurar la opción upsert para crear un documento si no existe
		options := options.Update().SetUpsert(true)

		// Realizar el "update or create" en la base de datos
		_, err := collection.UpdateOne(context.Background(), filter, update, options)
		if err != nil {
			return nil, err
		}

		// Después de actualizar, obtener el usuario actualizado
		var user entity.UserOnline
		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	} else {
		filter := bson.M{"alias": alias}
		_, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return nil, err
		}

		// Devolver nil ya que el usuario ya no está en línea
		return nil, nil
	}
}

func (r *UserRepositoryMongo) GetUserOnline(sessionToken, hash string) (*entity.UserOnline, error) {
	collection := r.client.Database("dbName").Collection("userOnline")

	var user entity.UserOnline
	err := collection.FindOne(context.Background(), bson.M{"sessionToken": sessionToken, "hash": hash}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Otros métodos de UserRepositoryMongo...
func (r *UserRepositoryMongo) GetUserByAlias(alias string) (*entity.User, error) {
	collection := r.client.Database("dbName").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"alias": alias}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

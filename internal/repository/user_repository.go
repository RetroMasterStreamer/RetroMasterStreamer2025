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
	collection := r.client.Database("portalRG").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"alias": bson.M{"$regex": alias, "$options": "i"}, "password": password}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (r *UserRepositoryMongo) SetUserOnline(alias, sessionToken, hash string, online bool) (*entity.UserOnline, error) {
	collection := r.client.Database("portalRG").Collection("userOnline")

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
		filter := bson.M{"alias": bson.M{"$regex": alias, "$options": "i"}}
		_, err := collection.DeleteMany(context.Background(), filter)
		if err != nil {
			return nil, err
		}

		// Devolver nil ya que el usuario ya no está en línea
		return nil, nil
	}
}

func (r *UserRepositoryMongo) GetUserOnline(sessionToken, hash string) (*entity.UserOnline, error) {
	collection := r.client.Database("portalRG").Collection("userOnline")

	var user entity.UserOnline
	err := collection.FindOne(context.Background(), bson.M{"sessionToken": sessionToken, "hash": hash}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Otros métodos de UserRepositoryMongo...
func (r *UserRepositoryMongo) GetUserByAlias(alias string) (*entity.User, error) {
	collection := r.client.Database("portalRG").Collection("user")

	// Usar una expresión regular para búsqueda insensible a mayúsculas y minúsculas
	filter := bson.M{"alias": bson.M{"$regex": alias, "$options": "i"}}

	var user entity.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryMongo) GetUserByTextRefer(text string) (*entity.User, error) {
	collection := r.client.Database("portalRG").Collection("user")

	// Usar una expresión regular para búsqueda insensible a mayúsculas y minúsculas
	filter := bson.M{"reference_text": bson.M{"$regex": text, "$options": "i"}}

	var user entity.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SaveUser guarda el contenido completo de un registro en la colección "user".
func (r *UserRepositoryMongo) SaveUser(user *entity.User) error {
	collection := r.client.Database("portalRG").Collection("user")

	// Verificar si el usuario ya existe en la base de datos
	existingUser, _ := r.GetUserByAlias(user.Alias)

	// Si el usuario existe, actualizamos su registro
	if existingUser != nil {
		filter := bson.M{"alias": bson.M{"$regex": user.Alias, "$options": "i"}}
		update := bson.M{"$set": user}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	} else {
		// Si el usuario no existe, insertamos un nuevo registro
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepositoryMongo) GetAllUsers() ([]*entity.User, error) {
	collection := r.client.Database("portalRG").Collection("user")

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

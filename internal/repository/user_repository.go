// repository/user_repository.go
package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Convertir el alias a minúsculas
	lowerAlias := strings.ToLower(alias)

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"alias": lowerAlias, "password": password}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (r *UserRepositoryMongo) SetUserOnline(alias, sessionToken, hash string, online bool) (*entity.UserOnline, error) {
	collection := r.client.Database("portalRG").Collection("userOnline")

	alias = strings.ToLower(alias)

	// Definir el filtro para encontrar el usuario
	filter := bson.M{"alias": alias, "sessionToken": sessionToken, "hash": hash}

	if online {

		// Definir los cambios que se van a realizar

		access := time.Now().Format("2006-1-2 15:4:5")
		update := bson.M{"$set": bson.M{"online": online, "access": access}}

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
		filter := bson.M{"alias": alias, "sessionToken": sessionToken, "hash": hash, "online": true}
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
	user.Alias = strings.ToLower(user.Alias)
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

func (r *UserRepositoryMongo) SaveTips(tip *entity.PostNew) (string, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	// Verificar si el tip ya existe en la base de datos por ID
	existingTip, _ := r.GetTipsByID(tip.ID)
	if existingTip != nil {
		// Si el tip existe por ID, actualizamos su registro
		filter := bson.M{"id": bson.M{"$regex": tip.ID, "$options": "i"}}

		update := bson.M{"$set": tip}

		_, err := collection.UpdateOne(context.Background(), filter, update)

		if err != nil {
			return "", err
		}
		// Retornar el ID del tip actualizado
		return tip.ID, nil
	} else {
		// Verificar si el tip ya existe por URL para evitar duplicidad
		existingTipByURL, err := r.GetTipsByURL(tip.URL) // Busca por la URL

		if (existingTipByURL != nil || err != nil) && tip.Type != "download" {
			// Si ya existe un tip con la misma URL, devolvemos un error para evitar duplicados
			return tip.ID, fmt.Errorf("tip with URL %s already exists", tip.URL)
		}
		// Si el tip no existe, insertamos un nuevo registro
		_, err = collection.InsertOne(context.Background(), tip)
		if err != nil {
			return "", err
		}
		// Retornar el ID del tip insertado
		return tip.ID, nil
	}
}

func (r *UserRepositoryMongo) GetTipsByURL(url string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	// Usar un filtro exacto para la URL
	filter := bson.M{"url": url}

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), filter).Decode(&tip)
	if err != nil {
		// Si no se encuentra un tip con esa URL, devolver nil sin error
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &tip, nil
}

func (r *UserRepositoryMongo) GetTipsByAuthor(alias string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	// Usar una expresión regular para búsqueda insensible a mayúsculas y minúsculas
	filter := bson.M{"author": bson.M{"$regex": alias, "$options": "i"}}

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), filter).Decode(&tip)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

func (r *UserRepositoryMongo) GetTipsByID(id string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	// Usar una expresión regular para búsqueda insensible a mayúsculas y minúsculas
	filter := bson.M{"id": bson.M{"$regex": id, "$options": "i"}}

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), filter).Decode(&tip)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

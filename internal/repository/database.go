package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"PortalCRG/internal/repository/entity"
)

type DataBase struct {
	ConnectionString string
	client           *mongo.Client
}

// Connect establece una conexión con la base de datos MongoDB utilizando la cadena de conexión proporcionada.
func (db *DataBase) Connect(connectionString string) error {
	/*
		clientOptions := options.Client().ApplyURI(connectionString)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return err
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			return err
		}

		db.client = client
		return nil
	*/

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("portalRG").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return err
	}
	db.client = client
	return nil
}

// FindUser busca un usuario por su nombre en la colección "user" y devuelve un puntero a la estructura User.
func (db *DataBase) FindUser(alias string) (*entity.User, error) {
	collection := db.client.Database("portalRG").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"alias": bson.M{"$regex": alias, "$options": "i"}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser crea un nuevo usuario en la colección "user".
func (db *DataBase) CreateUser(newUser *entity.User) error {
	collection := db.client.Database("portalRG").Collection("user")

	_, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return err
	}

	return nil
}

// Inicia la base de datos y consulta si existe el administrador, si no existe lo crea
func (db *DataBase) Init() {

	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")

	// Conectar a la base de datos
	err := db.Connect(connectionString)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	// Ejemplo de uso: encontrar un usuario
	user, err := db.FindUser("admin")
	//if err == nil {
	if user == nil {
		newUser := &entity.User{
			Name:          "Administrador",
			Alias:         "admin",
			Password:      "iddqd",
			ReferenceText: "Iddqd",
			UserRef:       "Administrador",
		}
		err = db.CreateUser(newUser)
		if err != nil {
			log.Fatal("Error creando usuario:", err)
		} else {
			log.Println("Nace el Admin")
		}
	} else {
		log.Println("Admin ya existe")
		log.Println("Actualizando el Avatar de los usuarios")

	}
	/*
		} else {
			log.Fatal("Error :", err.Error())
		}
	*/

}

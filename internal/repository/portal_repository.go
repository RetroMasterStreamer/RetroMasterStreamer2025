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
	collection := r.client.Database("portalRG").Collection("user")

	var user entity.User
	err := collection.FindOne(context.Background(), bson.M{"reference_text": userRef}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PortalRepositoryMongo) GetAllUsers() ([]*entity.User, error) {
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

// GetTipsWithPagination obtiene los tips con paginación
func (r *PortalRepositoryMongo) GetTipsWithPagination(skip, limit int64) ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"date", -1}})

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tips []*entity.PostNew
	for cursor.Next(context.Background()) {
		var tip entity.PostNew
		if err := cursor.Decode(&tip); err != nil {
			return nil, err
		}
		tips = append(tips, &tip)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tips, nil
}

func (r *PortalRepositoryMongo) GetAllTips() ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tips []*entity.PostNew
	for cursor.Next(context.Background()) {
		var tip entity.PostNew
		if err := cursor.Decode(&tip); err != nil {
			return nil, err
		}
		tips = append(tips, &tip)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tips, nil
}

func (r *PortalRepositoryMongo) GetTipByID(id string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&tip)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

func (r *PortalRepositoryMongo) GetTipByAuthor(author string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), bson.M{"author": author}).Decode(&tip)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

// GetTipsWithSearch obtiene los tips con paginación y búsqueda
func (r *PortalRepositoryMongo) GetTipsWithSearch(search string, skip, limit int64) ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"content": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"date", -1}})

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tips []*entity.PostNew
	for cursor.Next(context.Background()) {
		var tip entity.PostNew
		if err := cursor.Decode(&tip); err != nil {
			return nil, err
		}
		tips = append(tips, &tip)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tips, nil
}

func (r *PortalRepositoryMongo) DeleteTipByIDandAuthor(id, alias string) error {
	collection := r.client.Database("portalRG").Collection("tips")
	filter := bson.M{"id": id, "author": bson.M{"$regex": alias, "$options": "i"}}
	q, err := collection.DeleteOne(context.Background(), filter)
	log.Println(q.DeletedCount, " Eliminaciones!")
	return err
}

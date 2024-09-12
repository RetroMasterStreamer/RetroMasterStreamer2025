// repository/user_repository.go
package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

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
func (r *PortalRepositoryMongo) GetTipsWithPagination(skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"date", -1}})

	// Filtro para "type in []typeOfTips"
	filter := bson.M{"type": bson.M{"$in": typeOfTips}}

	// Ejecutar la consulta con el filtro
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

func (r *PortalRepositoryMongo) GetTipByURL(URL string) (*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	var tip entity.PostNew
	err := collection.FindOne(context.Background(), bson.M{"url": URL}).Decode(&tip)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

// GetTipsWithSearch obtiene los tips con paginación y búsqueda
func (r *PortalRepositoryMongo) GetTipsWithSearch(search string, skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	// Dividir el texto de búsqueda en palabras separadas por espacios
	searchWords := strings.Fields(search)

	// Crear filtros $regex para cada palabra en los campos title y content
	var orFilters []bson.M
	for _, word := range searchWords {
		escapedWord := regexp.QuoteMeta(word) // Escapar el texto para usarlo en la expresión regular
		orFilters = append(orFilters, bson.M{"title": bson.M{"$regex": escapedWord, "$options": "im"}})
		orFilters = append(orFilters, bson.M{"content": bson.M{"$regex": escapedWord, "$options": "im"}})

	}
	orFilters = append(orFilters, bson.M{"title": bson.M{"$regex": search, "$options": "im"}})
	orFilters = append(orFilters, bson.M{"content": bson.M{"$regex": search, "$options": "im"}})

	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": orFilters,
			},
			{
				"type": bson.M{
					"$in": typeOfTips,
				},
			},
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
		matchCount := 0
		if strings.Contains(strings.ToLower(tip.Title), strings.ToLower(search)) || strings.Contains(strings.ToLower(tip.Content), strings.ToLower(search)) {
			tip.Hash = append(tip.Hash, search)
			tip.MatchCount = 1000
			tips = append(tips, &tip)
		} else {
			for _, word := range searchWords {
				if (strings.Contains(strings.ToLower(tip.Title), strings.ToLower(word)) ||
					strings.Contains(strings.ToLower(tip.Content), strings.ToLower(word))) && len(word) > 2 {
					matchCount++
					tip.Hash = append(tip.Hash, word)
				}
			}
			tip.MatchCount = matchCount

			if matchCount > len(searchWords)/2 {
				tips = append(tips, &tip)
			}
		}

	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Ordenar los resultados por la cantidad de coincidencias (matchCount)
	sort.Slice(tips, func(i, j int) bool {
		return tips[i].MatchCount > tips[j].MatchCount
	})

	return tips, nil
}

func (r *PortalRepositoryMongo) DeleteTipByIDandAuthor(id, alias string) error {
	collection := r.client.Database("portalRG").Collection("tips")
	filter := bson.M{"id": id, "author": bson.M{"$regex": alias, "$options": "i"}}
	q, err := collection.DeleteOne(context.Background(), filter)
	log.Println(q.DeletedCount, " Eliminaciones!")
	return err
}

func (r *PortalRepositoryMongo) GetTipsByAliasWithPagination(alias string, skip, limit int64) ([]*entity.PostNew, error) {
	collection := r.client.Database("portalRG").Collection("tips")

	filter := bson.M{"author": bson.M{"$regex": alias, "$options": "i"}}
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

func (r *PortalRepositoryMongo) DeleteTipsFromDate(dateString string) error {
	collection := r.client.Database("portalRG").Collection("tips")

	// Parsear la fecha en formato "dd-mm-yyyy"
	parsedDate, err := time.Parse("02-01-2006", dateString)
	if err != nil {
		return fmt.Errorf("formato de fecha incorrecto: %v", err)
	}

	// Convertir la fecha a formato ISO para comparar en MongoDB
	filterDate := parsedDate.Format("2006-01-02T15:04:05.000Z")

	// Crear un filtro para eliminar los tips con fecha mayor o igual a la dada
	filter := bson.M{
		"date": bson.M{
			"$gte": filterDate,
		},
	}

	// Ejecutar la eliminación
	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error al eliminar tips: %v", err)
	}

	fmt.Printf("Eliminados %d documentos con fecha mayor o igual a %s\n", result.DeletedCount, dateString)
	return nil
}

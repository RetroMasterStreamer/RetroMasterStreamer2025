package util

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Eliminar documentos duplicados basados en el campo URL
func DeleteDuplicateURLs(collection *mongo.Collection) error {

	// Paso 1: Agrupar documentos por URL y contar cuántos hay en cada grupo
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$url"},
			{"count", bson.D{{"$sum", 1}}},
			{"docs", bson.D{{"$push", "$_id"}}},
		}}},
		// Paso 2: Filtrar los grupos que tienen más de 1 documento (es decir, duplicados)
		{{"$match", bson.D{
			{"count", bson.D{{"$gt", 1}}},
		}}},
	}

	// Ejecutar la consulta de agregación
	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())

	// Paso 3: Recorrer los resultados y eliminar duplicados, manteniendo solo el primero
	for cur.Next(context.Background()) {
		var result struct {
			Docs []interface{} `bson:"docs"`
		}
		if err := cur.Decode(&result); err != nil {
			return err
		}

		// Mantener el primer documento y eliminar el resto
		idsToDelete := result.Docs[1:] // Omitir el primer ID
		filter := bson.M{"_id": bson.M{"$in": idsToDelete}}

		// Eliminar documentos duplicados
		deleteResult, err := collection.DeleteMany(context.Background(), filter)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted %d duplicate documents for URL\n", deleteResult.DeletedCount)
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

func DeleteURLDuplicada() {
	// Configurar cliente MongoDB
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	// Conectar al cliente
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Seleccionar colección
	collection := client.Database("portalRG").Collection("tips")

	// Llamar a la función para eliminar duplicados basados en la URL
	if err := DeleteDuplicateURLs(collection); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Proceso de eliminación de duplicados completado.")
}

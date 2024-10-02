package mongo

import (
	"contact-api/internal/app/domain/models"
	"contact-api/internal/app/storage"
	"contact-api/internal/pkg/e"
	"contact-api/internal/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"time"
)

type DB struct {
	db *mongo.Client
}

func New(log *slog.Logger, ctx context.Context, connUrl string) (*DB, error) {
	const op = "storage.mongo.New"
	log = log.With(
		slog.String("op", op))

	log.Debug("connect uri", slog.String("url", connUrl))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connUrl))
	if err != nil {
		log.Error("Failed to connect to MongoDB", sl.Err(err))
		return nil, err
	}

	for i := 1; i <= 10; i++ {
		err = client.Ping(ctx, nil)
		if err == nil {
			log.Info("Successfully connected to MongoDB", slog.Int("try number", i))
			return &DB{client}, nil
		}
		log.Info("MongoDB connection failed, retrying in 5 seconds...", slog.Int("try number", i))
		time.Sleep(5 * time.Second)
	}

	// Если все попытки исчерпаны
	log.Error("Failed to connect to MongoDB after multiple attempts", sl.Err(err))
	return nil, err
}

func (db *DB) Close() {
	if err := db.db.Disconnect(context.Background()); err != nil {
		log.Println("Failed to disconnect MongoDB", sl.Err(err))
	}
}

func (db *DB) GetAll() ([]models.Contact, error) {
	var contactsRepo []Contact

	collection := db.db.Database("contacts").Collection("contact-list")

	//TODO: стоит перенести время на запрос в конфигурацию приложения
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, e.Err("failed to get all contacts", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &contactsRepo); err != nil {
		return nil, e.Err("failed to decode contacts", err)
	}

	contacts := RepoToContacts(contactsRepo)

	return contacts, nil
}

func (db *DB) Save(contact models.Contact) (string, error) {
	repoContact := ContactToRepoWithoutID(contact)

	collection := db.db.Database("contacts").Collection("contact-list")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, repoContact)
	if err != nil {
		return "", e.Err("failed to insert contact", err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *DB) DeleteAll() (int64, error) {
	collection := db.db.Database("contacts").Collection("contact-list")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("failed deleting contacts: %w", err)
	}

	return result.DeletedCount, nil
}

func (db *DB) ContactById(id string) (models.Contact, error) {
	var contactRepo = Contact{}

	collection := db.db.Database("contacts").Collection("contact-list")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoId, err := convertStringToObjectID(id)
	if err != nil {
		return models.Contact{}, e.Err("error convert id in mongo type", err)
	}

	filter := bson.D{{"_id", mongoId}}

	err = collection.FindOne(ctx, filter).Decode(&contactRepo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Contact{}, storage.ErrContactNotFound
		}

		return models.Contact{}, e.Err("failed to get contact", err)
	}

	contact := RepoToContact(contactRepo)

	return contact, nil
}

func (db *DB) Delete(id string) (bool, error) {
	collection := db.db.Database("contacts").Collection("contact-list")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoId, err := convertStringToObjectID(id)
	if err != nil {
		return false, e.Err("error convert id in mongo type", err)
	}

	filter := bson.D{{"_id", mongoId}}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, storage.ErrContactNotFound
		}

		return false, e.Err("failed to delete contact", err)
	}

	return true, nil
}

func (db *DB) Update(contact models.Contact) (bool, error) {
	contactRepo, err := ContactToRepo(contact)
	if err != nil {
		return false, e.Err("error convert to mongo models", err)
	}

	update := bson.M{
		"$set": contactRepo,
	}

	collection := db.db.Database("contacts").Collection("contact-list")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.UpdateByID(ctx, contactRepo.ID, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, storage.ErrContactNotFound
		}
		return false, e.Err("failed to update contact", err)
	}

	if result.MatchedCount == 0 {
		return false, storage.ErrContactNotFound
	}

	return true, nil
}

func convertStringToObjectID(idStr string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid ObjectID: %s", idStr)
	}
	return objectID, nil
}

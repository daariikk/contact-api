package mongo

import (
	"contact-api/internal/app/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserName  string             `bson:"username"`
	Email     string             `bson:"email"`
	Telephone Phone              `bson:"telephone"`
}

type Phone struct {
	Mobile string `bson:"mobile"`
	Home   string `bson:"home"`
}

// Преобразование Contact (репозиторий) в models.Contact (сервисный уровень)
func RepoToContact(repoContact Contact) models.Contact {
	return models.Contact{
		ID:       repoContact.ID.Hex(), // Конвертируем ObjectID в строку
		UserName: repoContact.UserName,
		Email:    repoContact.Email,
		Telephone: models.Phone{
			Mobile: repoContact.Telephone.Mobile,
			Home:   repoContact.Telephone.Home,
		},
	}
}

// Преобразование массива моделей репозитория в массив сервисных моделей
func RepoToContacts(repoContacts []Contact) []models.Contact {
	serviceContacts := make([]models.Contact, len(repoContacts))
	for i, repoContact := range repoContacts {
		serviceContacts[i] = RepoToContact(repoContact)
	}
	return serviceContacts
}

func ContactToRepo(serviceContact models.Contact) (Contact, error) {
	objectId, err := primitive.ObjectIDFromHex(serviceContact.ID)
	if err != nil {
		return Contact{}, err // Если некорректный ObjectID
	}

	repoContact := ContactToRepoWithoutID(serviceContact)
	repoContact.ID = objectId
	return repoContact, nil
}

func ContactToRepoWithoutID(serviceContact models.Contact) Contact {
	return Contact{
		UserName: serviceContact.UserName,
		Email:    serviceContact.Email,
		Telephone: Phone{
			Mobile: serviceContact.Telephone.Mobile,
			Home:   serviceContact.Telephone.Home,
		},
	}
}

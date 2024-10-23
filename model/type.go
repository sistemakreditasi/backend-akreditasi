package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Salt     string             `bson:"salt" json:"-"`    // Salt untuk hashing password, tidak dikirim dalam response JSON
	Role     string             `bson:"role" json:"role"` // Role: kaprodi, dosen, staff
}

type PDFDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FileName   string             `bson:"file_name"`
	FileID     string             `bson:"file_id"`
	UploadedBy string             `bson:"uploaded_by"`
	UploadedAt time.Time          `bson:"uploaded_at"`
}

package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Salt     string             `bson:"salt" json:"-"`    // Salt untuk hashing password, tidak dikirim dalam response JSON
	Role     string             `bson:"role" json:"role"` // Role: kaprodi, dosen, staff
}

// Akreditasi represents accreditation standards and file uploads
type Akreditasi struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Standard      int                `bson:"standard" json:"standard"`       // Accreditation standard (1-9)
	Description   string             `bson:"description" json:"description"` // Description of the standard
	CreatedBy     string             `bson:"created_by" json:"created_by"`   // Who created it
	FileLink      string             `bson:"file_link" json:"file_link"`     // Link to uploaded PDF
	Status        string             `bson:"status" json:"status"`           // Status (penetapan, pelaksanaan, evaluasi, etc.)
	CreationDate  primitive.DateTime `bson:"creation_date" json:"creation_date"`
	LastUpdatedBy string             `bson:"last_updated_by" json:"last_updated_by"`
}

// PDFDocument struct to handle uploaded files metadata
type PDFDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileName   string             `json:"file_name"`
	FileID     string             `json:"file_id"`
	UploadedBy string             `json:"uploaded_by"`
	UploadedAt primitive.DateTime `json:"uploaded_at"`
}

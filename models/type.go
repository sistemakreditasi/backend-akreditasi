package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User struct represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Role     string             `bson:"role" json:"role"` // Role can be "kaprodi", "dosen", "staff"
}

type Akreditasi struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Standard      int                `bson:"standard" json:"standard"`       // 1 to 9 (for each accreditation standard)
	Description   string             `bson:"description" json:"description"` // Description of the standard
	CreatedBy     string             `bson:"created_by" json:"created_by"`   // Kaprodi, Dosen, or Staff who created it
	FileLink      string             `bson:"file_link" json:"file_link"`     // Link to uploaded PDF (from Google Drive)
	Status        string             `bson:"status" json:"status"`           // Status (penetapan, pelaksanaan, evaluasi, pengendalian, peningkatan)
	CreationDate  primitive.DateTime `bson:"creation_date" json:"creation_date"`
	LastUpdatedBy string             `bson:"last_updated_by" json:"last_updated_by"` // User who last updated the record
}

package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID   primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	Capacity int                `bson:"capacity" json:"capacity"`
	From     time.Time          `bson:"from" json:"from"`
	Till     time.Time          `bson:"till" json:"till"`
}

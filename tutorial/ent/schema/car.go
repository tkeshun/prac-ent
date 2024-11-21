package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("model").SchemaType(map[string]string{"postgres": "text"}),
		field.String("car_number").SchemaType(map[string]string{"postgres": "text"}),
		field.Int("owner_id"),
		field.Time("registered_at").Default(time.Now), // Defaultでデフォルト値を指定できる
	}
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("cars").Field("owner_id").Unique().Required(),
	}
}

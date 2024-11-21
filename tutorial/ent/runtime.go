// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"
	"tutorial/ent/car"
	"tutorial/ent/schema"
	"tutorial/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	carFields := schema.Car{}.Fields()
	_ = carFields
	// carDescRegisteredAt is the schema descriptor for registered_at field.
	carDescRegisteredAt := carFields[3].Descriptor()
	// car.DefaultRegisteredAt holds the default value on creation for the registered_at field.
	car.DefaultRegisteredAt = carDescRegisteredAt.Default.(func() time.Time)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescAge is the schema descriptor for age field.
	userDescAge := userFields[0].Descriptor()
	// user.AgeValidator is a validator for the "age" field. It is called by the builders before save.
	user.AgeValidator = userDescAge.Validators[0].(func(int) error)
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[1].Descriptor()
	// user.DefaultName holds the default value on creation for the name field.
	user.DefaultName = userDescName.Default.(string)
}
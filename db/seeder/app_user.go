package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type User struct {
	ID        string
	Name      string
	FullName  string
	Email     string
	Password  string
	Active    bool
	Phone     string
	CreatedBy string
}

var users = []User{
	{
		ID:        "fcdb2142-4731-470e-8a1b-8d0037665fb2",
		Name:      "superadmin1",
		FullName:  "Super Admin Tes",
		Email:     "super@admin.com",
		Password:  "$2a$10$oPZtYt40bTpZxNS9Ayjj1OQizSbZ8f5DcOkjNFBJwkwf3qD8xQLiO",
		Phone:     "123143124",
		Active:    true,
		CreatedBy: "system",
	},
}

func userSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_user(
		id, name, full_name, email, password, phone, active, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	ON CONFLICT (id) DO NOTHING;
	`

	for i := range users {
		_, err := s.DB.Exec(statement,
			users[i].ID,
			users[i].Name,
			users[i].FullName,
			users[i].Email,
			users[i].Password,
			users[i].Phone,
			users[i].Active,
			users[i].CreatedBy,
			) 
			if err != nil {
			log.Fatalf("❌ ERROR execute user seeder : %v", err.Error())
		}
	}
}

package entity

type User struct {
	Nim          string `gorm:"primaryKey; type:varchar(15); not null; unique"`
	Email        string `gorm:"type:varchar(320); not null; unique"`
	FullName     string `gorm:"type:varchar(255); not null"`
	ProgramStudi string `gorm:"type:varchar(255); not null"`
	Role         string `gorm:"type:varchar(255); not null; default:'user'"`
}

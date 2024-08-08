package models

type City struct {
	BaseModel
	Name      string `gorm:"sizeL15;type:string;not null;"`
	CountryId int
	Country   Country `gorm:"foreignKey:CountryId"`
}

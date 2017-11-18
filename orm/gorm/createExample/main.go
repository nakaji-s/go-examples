package main

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

type Product struct {
	Id          string
	Description string
}

type ProductPtr struct {
	Id          *string
	Description *string
}

type ProductNullTag struct {
	Id          string
	Description string `sql:"default: null"`
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)

	db.DropTableIfExists(&Product{})
	db.AutoMigrate(&Product{})

	id001 := "id001"
	emptyString := ""

	//######################
	// Insert Null Patterns
	//######################
	// NG
	// INSERT INTO "products" ("id","description") VALUES ('id001','')
	//db.Create(&Product{Id: id001, Description: emptyString})

	// INSERT INTO "products" ("id","description") VALUES ('id001',NULL)
	db.Table("products").Create(&ProductPtr{Id: &id001})

	// INSERT INTO "products" ("id") VALUES ('id001')
	db.Table("products").Create(&ProductNullTag{Id: id001, Description: emptyString})

	//######################
	// Update Patterns
	//######################
	// UPDATE "products" SET "id" = 'id001'  WHERE "products"."id" = 'id001'
	db.Model(&Product{}).Update(&Product{Id: id001, Description: emptyString})

	// UPDATE "products" SET "description" = 'test', "id" = 'id001'  WHERE "products"."id" = 'id001'
	db.Model(&Product{}).Update(&Product{Id: id001, Description: "test"})

	// UPDATE "products" SET "id" = 'id001', "description" = ''  WHERE "products"."id" = 'id001'
	db.Model(&Product{}).Update(map[string]interface{}{"id": id001, "description": emptyString})

	// UPDATE "products" SET "description" = ''  WHERE "products"."id" = 'id001'
	db.Model(&Product{}).Save(&Product{Id: id001, Description: emptyString})

	// NG
	// UPDATE "products" SET "id" = '', "description" = ''
	//db.Model(&Product{}).Update(&ProductPtr{Id: &id001, Description: &emptyString})

}

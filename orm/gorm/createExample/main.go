package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kr/pretty"
)

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

	//######################
	// Update Null Patterns
	//######################
	// UPDATE "products" SET "id" = 'id001', "description" = NULL  WHERE "products"."id" = 'id001'
	db.Model(&Product{}).Update(map[string]interface{}{"id": id001, "description": gorm.Expr("NULL")})

	// UPDATE "products" SET "description" = NULL  WHERE (id = 'id001')
	db.Model(&Product{}).Where("id = ?", id001).Update("description", gorm.Expr("NULL"))

	// UPDATE "products" SET "description" = NULL  WHERE "products"."id" = 'id001'
	db.Table("products").Save(&ProductPtr{Id: &id001, Description: nil})

	// NG
	// UPDATE "products" SET "id" = 'id001'  WHERE (id = 'id001')
	//db.Table("products").Where("id = ?", id001).Update(&ProductPtr{Id: &id001, Description: nil})

	// NG
	// UPDATE "products" SET "description" = ''  WHERE "products"."id" = 'id001'
	//db.Table("products").Save(&ProductNullTag{Id: id001, Description: emptyString})

	//######################
	// Get Patterns
	//######################
	product := Product{}
	productPtr := ProductPtr{}

	// SELECT * FROM "products"   ORDER BY "products"."id" ASC LIMIT 1
	// main.Product{Id:"id001", Description:""}
	db.First(&product)
	pretty.Println(product)

	// SELECT * FROM "products"   ORDER BY "products"."id" ASC LIMIT 1
	//main.ProductPtr{
	//	Id:          &"id001",
	//	Description: (*string)(nil),
	//}
	db.Table("products").First(&productPtr)
	pretty.Println(productPtr)
}

package main

import (
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
	Raw   interface{} `gorm:"type:varchar(100)"`
}

func beforeCreate(scope *gorm.Scope) {
	//field, _ := scope.FieldByName("Raw")
	if product, ok := scope.Value.(*Product); ok {
		err := scope.SetColumn("raw", strings.Join(product.Raw.([]string), ","))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)

	db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", beforeCreate)

	var product Product
	db.DropTableIfExists(&product)

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	err = db.Create(&Product{Code: "L1212", Price: 1000, Raw: []string{"aaa", "bbb"}}).Error
	if err != nil {
		panic(err)
	}

	// Read
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}

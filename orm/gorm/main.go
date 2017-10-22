package main

import (
	"strings"

	"fmt"

	"database/sql/driver"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kr/pretty"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
	Raw   interface{}   `gorm:"type:varchar(100)"`
	Str   MyStringArray `gorm:"type:varchar(64)"`
}

type MyStringArray []string

func (a *MyStringArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		*a = strings.Split(string(src), `,`)
		return nil
	case string:
		*a = strings.Split(src, `,`)
		return nil
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to MyStringArray", src)
}

func (a MyStringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	ret := fmt.Sprint(strings.Join(a, `,`))

	return ret, nil
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
	//db, err := gorm.Open("postgres", "host=localhost user=postgres sslmode=disable password=admin")
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
	err = db.Create(&Product{Code: "L1212", Price: 1000, Raw: []string{"aaa", "bbb"}, Str: []string{"aaa", "bbb"}}).Error
	if err != nil {
		panic(err)
	}

	// Read
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212
	pretty.Println(product)

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}

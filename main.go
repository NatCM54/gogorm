package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n========================\n", sql)
}

var db *gorm.DB

func main() {
	dsn := "nut:M59P1Gjoan8z@tcp(159.89.199.251:3306)/lab_nut?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(Customer{}, Gender{}, Test{}, Employee{})
	// CreateGender("xxxx")
	// GetGenders()
	// GetGender(1)
	// GetGenderByName("Male")
	// UpdateGender2(1, "ชาย")
	// DeleteGender(4)
	// CreateTest(0, "Test1")
	// CreateTest(0, "Test2")
	// CreateTest(0, "Test3")
	// DeleteTest(2)
	// GetTests()

	// db.Migrator().CreateTable(Customer{})
	// CreateCustomer("Sim", 2)
	// GetCustomers()

}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Preload(clause.Associations).Find(&customers) //db.Preload("Gender").Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, customer := range customers {
		fmt.Printf("%v | %v | %v\n", customer.ID, customer.Name, customer.Gender.Name)
	}
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{Name: name, GenderID: genderID}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func CreateTest(code uint, name string) {
	test := Test{Code: code, Name: name}
	db.Create(&test)
}

func GetTests() {
	tests := []Test{}
	db.Find(&tests)
	for _, t := range tests {
		fmt.Printf("%v|%v\n", t.ID, t.Name)
	}
}

func DeleteTest(id uint) {
	db.Unscoped().Delete(&Test{}, id)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(int(id))
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(int(id))
}

func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id=@myid", sql.Named("myid", id)).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(int(id))
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders) //db.Where("name=?", name).Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGenderByName(name string) {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders, "name=?", name)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id int) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code  uint   `gorm:"comment:This is code"`
	Name  string `gorm:"column:myname;type:varchar(50);size:20;unique;default:hello; not null"`
	Email string
	Phone uint
}

type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderID uint
}

type Employee struct {
	gorm.Model
	Name  string
	Email string
	Phone uint
}

// func (t Test) TableName() string {
// 	return "My Test"
// }

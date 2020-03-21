package main

//one to on association
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "golangdb"
)

type UserL struct {
	ID int `gorm:"primary_key"`
	Uname string
	Languages []Language `gorm:"many2many:user_languages";"ForeignKey:UserId"`
	//Based on this 3rd table user_languages will be created
}

type Language struct {
	ID int `gorm:"primary_key"`
	Name string
}

type UserLanguages struct {
	UserLId int
	LanguageId int
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	defer db.Close()

	db.DropTableIfExists(&UserLanguages{}, &Language{}, &UserL{})
	db.AutoMigrate(&UserL{}, &Language{}, &UserLanguages{})

	//All foreign keys need to define here
	db.Model(UserLanguages{}).AddForeignKey("user_l_id", "user_ls(id)", "CASCADE", "CASCADE")
	db.Model(UserLanguages{}).AddForeignKey("language_id", "languages(id)", "CASCADE", "CASCADE")

	langs := []Language{{Name: "English"}, {Name: "French"}}
	log.Println(langs)

	user1 := UserL{Uname: "John", Languages: langs}
	user2 := UserL{Uname: "Martin", Languages: langs}
	user3 := UserL{Uname: "Ray", Languages: langs}

	db.Save(&user1) //save is happening
	db.Save(&user2)
	db.Save(&user3)

	fmt.Println("After Saving Records")
	fmt.Println("User1", &user1)
	fmt.Println("User2", &user2)
	fmt.Println("User3", &user3)

	//Fetching
	user := &UserL{}
	db.Debug().Where("uname=?", "Ray").Find(&user)
	err = db.Debug().Model(&user).Association("Languages").Find(&user.Languages).Error
	fmt.Println("User is now coming", user, err)

	//Deletion
	fmt.Println(user, "to delete")
	db.Debug().Where("uname=?", "John").Delete(&user)

	//Updation
	db.Debug().Model(&UserL{}).Where("uname=?", "Ray").Update("uname", "Martin")


}
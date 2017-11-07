package db

func fillMockData() {
	DB.DropTable(&Otto{})
	DB.DropTable(&User{})
	DB.AutoMigrate(&Otto{}, &User{})

	DB.Create(&Otto{Serial: "0ABCDE"})
	DB.Create(&Otto{Serial: "012345"})
}

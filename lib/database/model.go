package database

type Registration struct {
	RegID       uint64 `gorm:"primaryKey;autoIncrement"`
	Token       string
	Description string
	UserID      string
	UserName    string
}

type Alarm struct {
	AlarmID uint64 `gorm:"primaryKey"`
	UserID  string
	RegID   uint64
	Pattern string
	Message string
}

type FoodTag struct {
	FoodTagID uint64 `gorm:"primaryKey"`
	Name      string
}

type FoodTagRelation struct {
	FoodTagID  uint64 `gorm:"primaryKey"`
	VendorID   uint64 `gorm:"primaryKey"`
	VendorName string
	VendorCode string
	VendorURL  string
}

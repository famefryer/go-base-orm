package submode

// @GormRepository("product", "id")
type Product struct {
	ID    int `gorm:"primarykey"`
	Name  string
	Price float64
}

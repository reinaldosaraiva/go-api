package database

import (
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"gorm.io/gorm"
)

type Product struct{
	DB *gorm.DB

}

func NewProduct(db *gorm.DB) *Product{
	return &Product{DB:db}
}

func (p *Product) Create(product *entity.Product) error{
	return p.DB.Create(product).Error
}

func (p *Product) FindByID(id uint) (*entity.Product, error){
	var product entity.Product
	err := p.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product,nil
}
func (p *Product) Update(product *entity.Product) error{
	_, err := p.FindByID(product.ID)
	if err != nil {
		return err
	
	}
	return p.DB.Save(product).Error
}

func (p *Product) FindAll(page int, limit int, sort string) ([]entity.Product, error){
	var products []entity.Product
	var err error
	if sort == ""  || sort == "asc" || sort != "desc"{
		sort = "asc"
	}
	if page != 0 && limit != 0{
		err = p.DB.Limit(limit).Offset((page-1)*limit).Order("id "+sort).Find(&products).Error
	}else{
		err = p.DB.Order("CreatedAt "+sort).Find(&products).Error
	}
	return products, err
}

func (p *Product) Delete(id uint) error{
	product,err := p.FindByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}
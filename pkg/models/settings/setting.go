package settings

import (
	"github.com/jinzhu/gorm"

	"github.com/qorpress/qorpress/core/l10n"
	"github.com/qorpress/qorpress/core/location"
)

// to do
type FeeSetting struct {
	ShippingFee     uint
	GiftWrappingFee uint
	CODFee          uint `gorm:"column:cod_fee"`
	TaxRate         int
}

type Setting struct {
	gorm.Model
	FeeSetting
	location.Location `location:"name:Company Address"`
	l10n.Locale
}

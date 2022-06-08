package db

import (
	"testing"

	"github.com/reddio-com/red-adapter/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestIncreaseMintIndex(t *testing.T) {
	dsn := "root:password@tcp(127.0.0.1:3306)/reddifi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	rdb := NewDB(db)

	rdb.Exec("DROP TABLE IF EXISTS mint_models")

	db.AutoMigrate(&types.MintModel{})

	m := &types.MintModel{}
	m.Address = "0x0000000000000000000000000000000000000000"
	m.Index = 1

	var exists bool

	err = rdb.Model(&types.MintModel{}).
		Select("count(*) > 0").
		Where("address = ?", m.Address).
		Find(&exists).
		Error
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		db.Create(m)
	}

	m2 := &types.MintModel{}
	m2.Address = "0x0000000000000000000000000000000000000002"
	m2.Index = 20

	err = rdb.Model(&types.MintModel{}).
		Select("count(*) > 0").
		Where("address = ?", m2.Address).
		Find(&exists).
		Error
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		db.Create(m2)
	}

	id := rdb.IncreaseMintIndex(m)

	if id != 2 {
		t.Errorf("id should be 2, but got %d", id)
	}

	id = rdb.IncreaseMintIndex(m)

	if id != 3 {
		t.Errorf("id should be 3, but got %d", id)
	}

	id = rdb.IncreaseMintIndex(m2)

	if id != 21 {
		t.Errorf("id should be 21, but got %d", id)
	}
}

func TestGetContractInfo(t *testing.T) {
	dsn := "root:password@tcp(127.0.0.1:3306)/reddifi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	rdb := NewDB(db)

	rdb.Exec("DROP TABLE IF EXISTS contract_infos")

	db.AutoMigrate(&types.Contracts{})

	c := &types.Contracts{}
	c.Address = "0x4240e8b8c0b6e6464a13f555f6395bbfe1c4bdf1"
	c.Type = "ERC20M"
	c.Name = "Reddio20"
	c.Symbol = "R20"
	c.Decimals = 18
	c.Quantum = 1000000000000000000
	c.TotalSupply = 1000
	c.Count = 1000

	err = rdb.Create(c).Error
	if err != nil {
		t.Fatal(err)
	}

	c2 := &types.Contracts{}
	c2.Address = "0x2b071c1810ab8bf93aaaa71ce64113d57018856d"
	c2.Type = "ERC721"
	c2.Name = "Reddio721"
	c2.Symbol = "R721"
	c2.Decimals = 1
	c2.Quantum = 1
	c2.TotalSupply = 1000
	c2.TotalSupply = 10000

	err = rdb.Create(c2).Error
	if err != nil {
		t.Fatal(err)
	}

	info, _ := rdb.GetContractInfo("0x4240e8b8c0b6e6464a13f555f6395bbfe1c4bdf1")

	if info.Address != "0x4240e8b8c0b6e6464a13f555f6395bbfe1c4bdf1" {
		t.Errorf("address should be 0x4240e8b8c0b6e6464a13f555f6395bbfe1c4bdf1, but got %s", info.Address)
	}

	if info.Name != "Reddio20" {
		t.Errorf("name should be Reddio20, but got %s", info.Name)
	}

	if info.Symbol != "R20" {
		t.Errorf("symbol should be R20, but got %s", info.Symbol)
	}

	if info.Decimals != 18 {
		t.Errorf("decimals should be 18, but got %d", info.Decimals)
	}

	if info.TotalSupply != 1000 {
		t.Errorf("total supply should be 1000, but got %d", info.TotalSupply)
	}

	info, _ = rdb.GetContractInfo("0x2b071c1810ab8bf93aaaa71ce64113d57018856d")

	if info.Address != "0x2b071c1810ab8bf93aaaa71ce64113d57018856d" {
		t.Errorf("address should be 0x2b071c1810ab8bf93aaaa71ce64113d57018856d, but got %s", info.Address)
	}

	if info.Name != "Reddio721" {
		t.Errorf("name should be Reddio721, but got %s", info.Name)
	}

	if info.Symbol != "R721" {
		t.Errorf("symbol should be R721, but got %s", info.Symbol)
	}

	if info.Decimals != 1 {
		t.Errorf("decimals should be 1, but got %d", info.Decimals)
	}

	if info.TotalSupply != 10000 {
		t.Errorf("total supply should be 200, but got %d", info.TotalSupply)
	}

}

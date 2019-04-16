package gorm_test

import (
	"time"

	"github.com/housinganywhere/gorm"
	"testing"
)

func NameIn1And2(d *gorm.DB) *gorm.DB {
	return d.Where("name in (?)", []string{"ScopeUser1", "ScopeUser2"})
}

func NameIn2And3(d *gorm.DB) *gorm.DB {
	return d.Where("name in (?)", []string{"ScopeUser2", "ScopeUser3"})
}

func NameIn(names []string) func(d *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("name in (?)", names)
	}
}

func TestScopes(t *testing.T) {
	user1 := User{Name: "ScopeUser1", Age: 1}
	user2 := User{Name: "ScopeUser2", Age: 1}
	user3 := User{Name: "ScopeUser3", Age: 2}
	DB.Save(&user1).Save(&user2).Save(&user3)

	var users1, users2, users3 []User
	DB.Scopes(NameIn1And2).Find(&users1)
	if len(users1) != 2 {
		t.Errorf("Should found two users's name in 1, 2")
	}

	DB.Scopes(NameIn1And2, NameIn2And3).Find(&users2)
	if len(users2) != 1 {
		t.Errorf("Should found one user's name is 2")
	}

	DB.Scopes(NameIn([]string{user1.Name, user3.Name})).Find(&users3)
	if len(users3) != 2 {
		t.Errorf("Should found two users's name in 1, 3")
	}
}

func TestAddToVars(t *testing.T) {
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	testLoc := time.FixedZone("test TZ", secondsEastOfUTC)
	tm := time.Date(2016, 2, 10, 15, 12, 0, 0, testLoc)
	scope := DB.NewScope(nil)
	scope.AddToVars(tm)
	if len(scope.SQLVars) != 1 {
		t.Errorf("Expected sql vars slice to be 1 element")
	}
	if scope.SQLVars[0] != tm.UTC() {
		t.Errorf("Expected time to be casted to UTC")
	}
}

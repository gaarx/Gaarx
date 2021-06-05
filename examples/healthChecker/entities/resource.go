package entities

import "github.com/jinzhu/gorm"

type Resource struct {
	gorm.Model
	Url string
}

const ScopeResources = "RSR"
const EventReloadResources = "RLS"

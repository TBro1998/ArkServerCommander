package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	IP           string         `json:"ip" gorm:"not null"`
	Port         int            `json:"port" gorm:"not null;default:7777"`
	QueryPort    int            `json:"query_port" gorm:"not null;default:27015"`
	RCONPort     int            `json:"rcon_port" gorm:"not null;default:27020"`
	RCONPassword string         `json:"-" gorm:"not null"`
	Map          string         `json:"map" gorm:"default:'TheIsland'"`
	MaxPlayers   int            `json:"max_players" gorm:"default:70"`
	Status       string         `json:"status" gorm:"default:'stopped'"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ServerRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	IP           string `json:"ip" binding:"required"`
	Port         int    `json:"port" binding:"required,min=1,max=65535"`
	QueryPort    int    `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort     int    `json:"rcon_port" binding:"required,min=1,max=65535"`
	RCONPassword string `json:"rcon_password" binding:"required"`
	Map          string `json:"map"`
	MaxPlayers   int    `json:"max_players" binding:"required,min=1,max=200"`
}

type ServerResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	QueryPort   int    `json:"query_port"`
	RCONPort    int    `json:"rcon_port"`
	Map         string `json:"map"`
	MaxPlayers  int    `json:"max_players"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ServerUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IP          string `json:"ip"`
	Port        int    `json:"port" binding:"min=1,max=65535"`
	QueryPort   int    `json:"query_port" binding:"min=1,max=65535"`
	RCONPort    int    `json:"rcon_port" binding:"min=1,max=65535"`
	Map         string `json:"map"`
	MaxPlayers  int    `json:"max_players" binding:"min=1,max=200"`
}

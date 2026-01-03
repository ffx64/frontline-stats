package dtos

import (
	"time"
)

type KillsDTO struct {
	ID               string    `json:"id"`
	ServerID         string    `json:"server_id"`
	RoundID          string    `json:"round_id"`
	KillerID         string    `json:"killer_id"`
	VictimID         string    `json:"victim_id"`
	VictimWeaponName string    `json:"victim_weapon_name"`
	VictimWeaponType string    `json:"victim_weapon_type"`
	KillerWeaponName string    `json:"killer_weapon_name"`
	KillerWeaponType string    `json:"killer_weapon_type"`
	HitZone          string    `json:"hit_zone"`
	Distance         float64   `json:"distance"`
	IsHeadshot       bool      `json:"is_headshot"`
	IsFriendly       bool      `json:"is_friendly"`
	IsVehicle        bool      `json:"is_vehicle"`
	KillerTeam       string    `json:"killer_team"`
	VictimTeam       string    `json:"victim_team"`
	Timestamp        time.Time `json:"timestamp"`
	CreatedAt        time.Time `json:"created_at"`
}

type KillsSaveDTO struct {
	ServerID         string     `json:"server_id" validate:"required,uuid4"`
	RoundID          string     `json:"round_id" validate:"required,uuid4"`
	KillerID         string     `json:"killer_id" validate:"required,uuid4"`
	VictimID         string     `json:"victim_id" validate:"required,uuid4"`
	VictimWeaponName string     `json:"victim_weapon_name" validate:"required"`
	VictimWeaponType string     `json:"victim_weapon_type" validate:"required"`
	KillerWeaponName string     `json:"killer_weapon_name" validate:"required"`
	KillerWeaponType string     `json:"killer_weapon_type" validate:"required"`
	HitZone          string     `json:"hit_zone" validate:"required"`
	Distance         float64    `json:"distance"`
	IsHeadshot       bool       `json:"is_headshot"`
	IsFriendly       bool       `json:"is_friendly"`
	IsVehicle        bool       `json:"is_vehicle"`
	KillerTeam       string     `json:"killer_team" validate:"required"`
	VictimTeam       string     `json:"victim_team" validate:"required"`
	Timestamp        *time.Time `json:"timestamp" validate:"required"`
}

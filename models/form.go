package models

import "time"

type CreateLinkForm struct {
	TargetUrl         string `form:"targetUrl" binding:"required"`
	Ttl               int    `form:"ttl" binding:"required"`
	SingleUse         string `form:"single-use"`
	PasswordProtected string `form:"password-protected"`
	Password          string `form:"password"`
}

type CreateLink struct {
	TargetUrl         string
	Ttl               time.Duration
	SingleUse         bool
	PasswordProtected bool
	Password          string
}

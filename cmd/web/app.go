package main

import (
	"github.com/alexedwards/scs"
)

type App struct {
	DSN       string
	HTMLDir   string
	Sessions  *scs.Manager
	StaticDir string
	StoreDir  string
	SecretKey string
}

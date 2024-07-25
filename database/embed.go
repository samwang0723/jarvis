package database

import "embed"

//go:embed migrations/*
var MigrationFiles embed.FS

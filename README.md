# Simple Reminder App

A minimal Go web app for reminders, built using Clean Architecture principles.

## Features
- Add reminders
- List reminders
- Remove reminders

## Tech
- Go 1.21+
- Gorilla Mux

## Structure
- `/cmd` - Main app entrypoint
- `/internal/core` - Entities and interfaces
- `/internal/usecase` - Business logic
- `/internal/adapter/http` - HTTP handlers
- `/internal/adapter/repo/mem` - In-memory repository

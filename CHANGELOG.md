# Changelog

All notable changes to this project are documented here. Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [0.1.0] - first prototype

### Added 
- Core domain model: applications, contacts, documents, reminders, activity history, status lifecycle rules
- Application workflow services: create, list, get, update details, update status, remove 
- Contact and document workflows: add, list, update, remove
- Reminder workflows: schedule, list, complete, update, remove
- PostgreSQL persistence layer with migrations
- Application search and reminder prioritization
- Activity and status history tracking
- Full HTTP API covering the workflows above
- React/TypeScript frontend: Dashboard, applications workbench, appliication detail and editing, contact/documentation/reminder management, search and filtering for view tables
- Demo video and GIF walkthroughs in README

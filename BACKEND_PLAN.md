# GlobalShot Backend Implementation Plan

## System Overview
A construction progress tracking system where Admins and Companies manage construction sites, upload 360° photos for specific rooms, and Clients view the progress of their assigned units.

## 1. Domain Model & Database Schema

### A. Users & Roles
We need a unified `users` table with roles to handle authentication and basic profile info.
- **Roles**: `SUPER_ADMIN` (GlobalShot), `COMPANY_ADMIN` (Client Companies), `CLIENT` (End users).

**Table: `users`**
- `id` (UUID)
- `email` (Unique)
- `password_hash`
- `role` (Enum: ADMIN, COMPANY, CLIENT)
- `company_id` (FK to `companies`, nullable - for Company Admins and Clients of Companies)
- `created_by` (FK to `users`, who created this user)

### B. Hierarchy Entities

**Table: `companies`**
- `id` (UUID)
- `name`
- `created_at`

**Table: `construction_sites`**
- `id` (UUID)
- `name`
- `address`
- `company_id` (FK to `companies` - owner of the site. If NULL, owned by GlobalShot Admin?)

**Table: `units`** (Represents a House or Flat)
- `id` (UUID)
- `name` (e.g., "Apartment 101", "Villa A")
- `type` (Enum: HOUSE, FLAT)
- `site_id` (FK to `construction_sites`)
- `client_id` (FK to `users` - the Client assigned to this unit)

**Table: `rooms`**
- `id` (UUID)
- `name` (e.g., "Living Room", "Kitchen")
- `unit_id` (FK to `units`)

### C. Media Content

**Table: `media`**
- `id` (UUID)
- `room_id` (FK to `rooms`)
- `url` (Storage URL of the 360 image)
- `thumbnail_url`
- `uploaded_by` (FK to `users`)
- `taken_at` (Date the photo was taken, for progress tracking)
- `created_at`

## 2. API Endpoints (REST)

### Auth
- `POST /auth/login` (Already exists)
- `POST /auth/reset-password`

### Companies Management (Admin Only)
- `GET /companies`
- `POST /companies`
- `GET /companies/:id`
- `PUT /companies/:id`
- `DELETE /companies/:id`

### User Management
- `GET /users` (Filter by role, company)
- `POST /users` (Create Client or Company Admin)
- `PUT /users/:id`
- `DELETE /users/:id`

### Construction Management
- `GET /sites`
- `POST /sites`
- `GET /sites/:id`
// ... standard CRUD for Sites, Units, Rooms

### Media
- `POST /media/upload` (Multipart form data)
- `GET /rooms/:id/media` (List history of images for a room)

## 3. Technology Stack
- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM (recommended for rapid relationship handling) or raw SQL
- **Storage**: Local filesystem (dev) / S3-compatible (prod)

## 4. Next Steps
1.  Define GORM models in `internal/model`.
2.  Set up database migrations.
3.  Implement CRUD Handlers for `Company` and `Site` first.

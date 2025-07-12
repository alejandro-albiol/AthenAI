# Database Schema Documentation

## Overview

This document describes the database schema for the Athenai multitenancy gym application with AI-powered exercise generation.

## Public Schema Tables

### 1. `public.gym`

Stores information about gym tenants in the multitenancy system.

| Column            | Type                     | Constraints             | Description                    |
| ----------------- | ------------------------ | ----------------------- | ------------------------------ |
| `id`              | UUID                     | PRIMARY KEY             | Unique identifier              |
| `name`            | TEXT                     | NOT NULL                | Gym name                       |
| `domain`          | TEXT                     | NOT NULL, UNIQUE        | Subdomain for tenant isolation |
| `email`           | TEXT                     | NOT NULL                | Contact email                  |
| `address`         | TEXT                     | NOT NULL                | Physical address               |
| `phone`           | TEXT                     | NOT NULL                | Contact phone                  |
| `is_active`       | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status                  |
| `deleted_at`      | TIMESTAMP WITH TIME ZONE | NULL                    | Soft delete timestamp          |
| `created_at`      | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp             |
| `updated_at`      | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp          |
| `business_hours`  | JSONB                    | NOT NULL, DEFAULT '[]'  | Operating hours                |
| `social_links`    | JSONB                    | NOT NULL, DEFAULT '[]'  | Social media links             |
| `payment_methods` | JSONB                    | NOT NULL, DEFAULT '[]'  | Accepted payment methods       |

### 2. `public.admin`

Stores system administrators who can manage the platform.

| Column          | Type                     | Constraints             | Description           |
| --------------- | ------------------------ | ----------------------- | --------------------- |
| `id`            | UUID                     | PRIMARY KEY             | Unique identifier     |
| `username`      | TEXT                     | NOT NULL, UNIQUE        | Admin username        |
| `password_hash` | TEXT                     | NOT NULL                | Hashed password       |
| `email`         | TEXT                     | NOT NULL, UNIQUE        | Admin email           |
| `is_active`     | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status         |
| `created_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp    |
| `updated_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp |

### 3. `public.exercise`

Core table for storing exercise information, including AI-generated exercises.

| Column             | Type                     | Constraints             | Description                                                  |
| ------------------ | ------------------------ | ----------------------- | ------------------------------------------------------------ |
| `id`               | UUID                     | PRIMARY KEY             | Unique identifier                                            |
| `name`             | TEXT                     | NOT NULL                | Exercise name                                                |
| `synonyms`         | TEXT[]                   | NOT NULL, DEFAULT '{}'  | Alternative names (AI-generated)                             |
| `muscular_groups`  | TEXT[]                   | NOT NULL, DEFAULT '{}'  | Affected muscle groups                                       |
| `equipment_needed` | TEXT[]                   | NOT NULL, DEFAULT '{}'  | Required equipment                                           |
| `difficulty_level` | TEXT                     | NOT NULL, CHECK         | 'beginner', 'intermediate', 'advanced'                       |
| `exercise_type`    | TEXT                     | NOT NULL, CHECK         | 'strength', 'cardio', 'flexibility', 'balance', 'functional' |
| `instructions`     | TEXT                     | NOT NULL                | How to perform the exercise                                  |
| `video_url`        | TEXT                     | NULL                    | Instructional video URL                                      |
| `image_url`        | TEXT                     | NULL                    | Exercise image URL                                           |
| `ai_generated`     | BOOLEAN                  | NOT NULL, DEFAULT FALSE | Whether AI generated this exercise                           |
| `ai_model_version` | TEXT                     | NULL                    | AI model version used                                        |
| `created_by`       | UUID                     | REFERENCES admin(id)    | Admin who created the exercise                               |
| `is_active`        | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status                                                |
| `created_at`       | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                                           |
| `updated_at`       | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp                                        |

### 4. `public.muscular_group`

Reference table for muscle groups.

| Column        | Type                     | Constraints             | Description                                     |
| ------------- | ------------------------ | ----------------------- | ----------------------------------------------- |
| `id`          | UUID                     | PRIMARY KEY             | Unique identifier                               |
| `name`        | TEXT                     | NOT NULL, UNIQUE        | Muscle group name                               |
| `description` | TEXT                     | NULL                    | Description                                     |
| `body_part`   | TEXT                     | NOT NULL, CHECK         | 'upper_body', 'lower_body', 'core', 'full_body' |
| `is_active`   | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status                                   |
| `created_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                              |
| `updated_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp                           |

### 5. `public.equipment`

Reference table for exercise equipment.

| Column        | Type                     | Constraints             | Description                                                       |
| ------------- | ------------------------ | ----------------------- | ----------------------------------------------------------------- |
| `id`          | UUID                     | PRIMARY KEY             | Unique identifier                                                 |
| `name`        | TEXT                     | NOT NULL, UNIQUE        | Equipment name                                                    |
| `description` | TEXT                     | NULL                    | Description                                                       |
| `category`    | TEXT                     | NOT NULL, CHECK         | 'free_weights', 'machines', 'cardio', 'accessories', 'bodyweight' |
| `is_active`   | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status                                                     |
| `created_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                                                |
| `updated_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp                                             |



## Indexes

The following indexes are created for optimal performance:

- `idx_gym_domain` - Fast gym lookup by domain
- `idx_gym_active` - Filter active gyms
- `idx_exercise_muscular_groups` - GIN index for array searches
- `idx_exercise_equipment` - GIN index for equipment searches
- `idx_exercise_difficulty` - Filter by difficulty level
- `idx_exercise_type` - Filter by exercise type
- `idx_admin_username` - Fast admin lookup
- `idx_admin_email` - Fast admin email lookup
- `idx_muscular_group_body_part` - Filter by body part
- `idx_equipment_category` - Filter by equipment category

## Setup Instructions

### 1. Using the Setup Script

```powershell
.\scripts\setup-db.ps1
```

### 2. Manual Setup

```bash
go run cmd/setup-db/main.go
```

### 3. Programmatic Setup

```go
db, err := database.NewPostgresDB()
if err != nil {
    log.Fatal(err)
}
defer db.Close()

err = database.CreatePublicTables(db)
if err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. **UUIDs**: All primary keys use UUIDs for better distribution and security
2. **Soft Deletes**: Use `deleted_at` for soft deletes instead of hard deletes
3. **Timestamps**: All tables include `created_at` and `updated_at` for audit trails
4. **Constraints**: Use CHECK constraints to ensure data integrity
5. **Indexes**: GIN indexes for array columns, B-tree for regular columns
6. **JSONB**: Use JSONB for flexible data structures (business hours, social links, etc.)
7. **Foreign Keys**: Proper referential integrity with foreign key constraints

## Future Considerations

- **Tenant Isolation**: Each gym gets its own schema for complete data isolation
- **AI Integration**: Built-in support for tracking AI model performance and costs
- **Scalability**: Indexes optimized for common query patterns
- **Audit Trail**: Complete logging of AI exercise generation for compliance

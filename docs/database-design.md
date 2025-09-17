# Database Design

## Overview

AthenAI uses a PostgreSQL database with a multi-tenant schema design that provides complete data isolation between gyms while sharing common resources efficiently.

## 🏗️ Schema Architecture

### Multi-Tenant Design Pattern

```
PostgreSQL Database: athenai
├── public schema                    # Platform-wide shared data
│   ├── gym                         # Tenant registry & metadata
│   ├── admin                       # Platform administrators
│   ├── exercise                    # Global exercise library
│   ├── equipment                   # Global equipment catalog
│   ├── muscular_group              # Muscle group definitions
│   ├── template_block              # Shared workout components
│   ├── workout_template            # Public workout templates
│   ├── exercise_equipment          # Global exercise-equipment relationships
│   └── exercise_muscular_group     # Global exercise-muscle relationships
│
└── {gym_uuid} schemas              # One schema per gym tenant
    ├── users                       # Gym members, trainers, admins
    ├── custom_equipment            # Gym-specific equipment
    ├── custom_exercise             # Gym-created exercises
    ├── custom_exercise_equipment   # Custom exercise equipment links
    ├── custom_exercise_muscular_group # Custom exercise muscle targeting
    ├── custom_template_block       # Gym-specific template components
    ├── custom_workout_template     # Gym workout templates
    ├── custom_member_workout       # Workout assignments to members
    ├── custom_workout_instance     # Active/completed workout sessions
    └── custom_workout_exercise     # Individual exercises within workouts
```

## 📋 Table Specifications

### **Public Schema Tables**

#### `public.gym` - Tenant Registry

Central registry for all gym tenants in the system.

| Column            | Type                     | Constraints             | Description                                 |
| ----------------- | ------------------------ | ----------------------- | ------------------------------------------- |
| `id`              | UUID                     | PRIMARY KEY             | Unique gym identifier (used as schema name) |
| `name`            | TEXT                     | NOT NULL                | Gym business name                           |
| `domain`          | TEXT                     | NOT NULL, UNIQUE        | Subdomain for tenant routing                |
| `email`           | TEXT                     | NOT NULL                | Primary contact email                       |
| `address`         | TEXT                     | NOT NULL                | Physical gym address                        |
| `phone`           | TEXT                     | NOT NULL                | Contact phone number                        |
| `is_active`       | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Operational status                          |
| `business_hours`  | JSONB                    | NOT NULL, DEFAULT '[]'  | Operating schedule                          |
| `social_links`    | JSONB                    | NOT NULL, DEFAULT '[]'  | Social media profiles                       |
| `payment_methods` | JSONB                    | NOT NULL, DEFAULT '[]'  | Accepted payment types                      |
| `deleted_at`      | TIMESTAMP WITH TIME ZONE | NULL                    | Soft delete timestamp                       |
| `created_at`      | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Registration date                           |
| `updated_at`      | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last modification                           |

#### `public.admin` - Platform Administrators

System administrators who manage the entire platform.

| Column          | Type                     | Constraints             | Description             |
| --------------- | ------------------------ | ----------------------- | ----------------------- |
| `id`            | UUID                     | PRIMARY KEY             | Unique admin identifier |
| `username`      | TEXT                     | NOT NULL, UNIQUE        | Admin login username    |
| `password_hash` | TEXT                     | NOT NULL                | Bcrypt hashed password  |
| `email`         | TEXT                     | NOT NULL, UNIQUE        | Admin email address     |
| `is_active`     | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Account status          |
| `created_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Account creation        |
| `updated_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update             |

#### `public.exercise` - Global Exercise Library

Comprehensive exercise database shared across all gyms.

| Column             | Type                     | Constraints             | Description                                                  |
| ------------------ | ------------------------ | ----------------------- | ------------------------------------------------------------ |
| `id`               | UUID                     | PRIMARY KEY             | Unique exercise identifier                                   |
| `name`             | TEXT                     | NOT NULL                | Exercise name                                                |
| `synonyms`         | TEXT[]                   | NOT NULL, DEFAULT '{}'  | Alternative names                                            |
| `difficulty_level` | TEXT                     | NOT NULL, CHECK         | 'beginner', 'intermediate', 'advanced'                       |
| `exercise_type`    | TEXT                     | NOT NULL, CHECK         | 'strength', 'cardio', 'flexibility', 'balance', 'functional' |
| `instructions`     | TEXT                     | NOT NULL                | Step-by-step exercise instructions                           |
| `video_url`        | TEXT                     | NULL                    | Instructional video link                                     |
| `image_url`        | TEXT                     | NULL                    | Exercise demonstration image                                 |
| `ai_generated`     | BOOLEAN                  | NOT NULL, DEFAULT FALSE | AI-created exercise flag                                     |
| `ai_model_version` | TEXT                     | NULL                    | AI model version used                                        |
| `created_by`       | UUID                     | REFERENCES admin(id)    | Admin who added exercise                                     |
| `is_active`        | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Availability status                                          |
| `created_at`       | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                                           |
| `updated_at`       | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last modification                                            |

#### `public.equipment` - Global Equipment Catalog

Standard gym equipment available across all facilities.

| Column        | Type                     | Constraints             | Description                                                       |
| ------------- | ------------------------ | ----------------------- | ----------------------------------------------------------------- |
| `id`          | UUID                     | PRIMARY KEY             | Unique equipment identifier                                       |
| `name`        | TEXT                     | NOT NULL, UNIQUE        | Equipment name                                                    |
| `description` | TEXT                     | NULL                    | Equipment description                                             |
| `category`    | TEXT                     | NOT NULL, CHECK         | 'free_weights', 'machines', 'cardio', 'accessories', 'bodyweight' |
| `is_active`   | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Availability status                                               |
| `created_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                                                |
| `updated_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last modification                                                 |

#### `public.muscular_group` - Muscle Group Definitions

Anatomical muscle group classifications.

| Column        | Type                     | Constraints             | Description                                     |
| ------------- | ------------------------ | ----------------------- | ----------------------------------------------- |
| `id`          | UUID                     | PRIMARY KEY             | Unique muscle group identifier                  |
| `name`        | TEXT                     | NOT NULL, UNIQUE        | Muscle group name                               |
| `description` | TEXT                     | NULL                    | Anatomical description                          |
| `body_part`   | TEXT                     | NOT NULL, CHECK         | 'upper_body', 'lower_body', 'core', 'full_body' |
| `is_active`   | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Active status                                   |
| `created_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Creation timestamp                              |
| `updated_at`  | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last modification                               |

#### Relationship Tables

**`public.exercise_equipment`** - Links exercises to required equipment

- `id` (UUID, PRIMARY KEY)
- `exercise_id` (UUID, REFERENCES exercise(id))
- `equipment_id` (UUID, REFERENCES equipment(id))

**`public.exercise_muscular_group`** - Maps exercises to target muscles

- `id` (UUID, PRIMARY KEY)
- `exercise_id` (UUID, REFERENCES exercise(id))
- `muscular_group_id` (UUID, REFERENCES muscular_group(id))

**`public.template_block`** - Reusable workout components

- Standard template building blocks for workout structure

**`public.workout_template`** - Shareable workout templates

- Pre-built workout structures available to all gyms

### **Tenant Schema Tables** (`{gym_uuid}`)

#### `{gym_uuid}.users` - Gym Members and Staff

All users associated with a specific gym.

| Column              | Type                     | Constraints             | Description                         |
| ------------------- | ------------------------ | ----------------------- | ----------------------------------- |
| `id`                | UUID                     | PRIMARY KEY             | Unique user identifier              |
| `username`          | TEXT                     | NOT NULL, UNIQUE        | User login name                     |
| `email`             | TEXT                     | NOT NULL, UNIQUE        | User email address                  |
| `password_hash`     | TEXT                     | NOT NULL                | Bcrypt hashed password              |
| `role`              | TEXT                     | NOT NULL, CHECK         | 'admin', 'trainer', 'member'        |
| `verified`          | BOOLEAN                  | NOT NULL, DEFAULT FALSE | Email verification status           |
| `description`       | TEXT                     | NULL                    | User bio/description                |
| `training_phase`    | TEXT                     | NOT NULL                | Current training focus              |
| `motivation`        | TEXT                     | NOT NULL                | Fitness motivation/goals            |
| `special_situation` | TEXT                     | NOT NULL                | Medical considerations, limitations |
| `is_active`         | BOOLEAN                  | NOT NULL, DEFAULT TRUE  | Account status                      |
| `created_at`        | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Registration date                   |
| `updated_at`        | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last profile update                 |

#### Custom Entity Tables

Each gym can extend the global catalogs with custom entries:

- **`{gym_uuid}.custom_equipment`** - Gym-specific equipment not in global catalog
- **`{gym_uuid}.custom_exercise`** - Gym-created exercises with same structure as public.exercise
- **`{gym_uuid}.custom_exercise_equipment`** - Links custom exercises to equipment
- **`{gym_uuid}.custom_exercise_muscular_group`** - Maps custom exercises to muscles

#### Workout Management Tables

- **`{gym_uuid}.custom_template_block`** - Gym-specific workout components
- **`{gym_uuid}.custom_workout_template`** - Gym workout templates
- **`{gym_uuid}.custom_member_workout`** - Workout plans assigned to specific members
- **`{gym_uuid}.custom_workout_instance`** - Actual workout sessions (filled templates)
- **`{gym_uuid}.custom_workout_exercise`** - Individual exercises within workout instances

## 🔗 Key Relationships

### Cross-Schema References

```sql
-- Custom exercises can reference public equipment and muscle groups
{gym_uuid}.custom_exercise_equipment.equipment_id → public.equipment.id
{gym_uuid}.custom_exercise_muscular_group.muscular_group_id → public.muscular_group.id

-- Workout instances can use both public and custom exercises
{gym_uuid}.custom_workout_exercise.public_exercise_id → public.exercise.id
{gym_uuid}.custom_workout_exercise.gym_exercise_id → {gym_uuid}.custom_exercise.id
```

### Workout Hierarchy

```
custom_workout_template (gym's workout structure)
    ↓
custom_workout_instance (filled template for member)
    ↓
custom_workout_exercise (individual exercises in workout)
    ↓ references
public.exercise OR {gym_uuid}.custom_exercise
```

## 🔐 Security & Access Control

### Schema-Level Isolation

- **Database users** have access only to their gym's schema
- **Application-level** validation ensures gym_id matches JWT claims
- **Row-level security** policies can be added for additional protection

### Data Integrity

```sql
-- Example: Ensure workout exercises belong to correct gym
CREATE POLICY gym_isolation ON {gym_uuid}.custom_workout_exercise
FOR ALL TO app_user
USING (
  workout_instance_id IN (
    SELECT id FROM {gym_uuid}.custom_workout_instance
  )
);
```

## 📊 Performance Optimizations

### Indexes

```sql
-- Public schema indexes
CREATE INDEX idx_exercise_search ON public.exercise USING GIN(name gin_trgm_ops);
CREATE INDEX idx_exercise_type ON public.exercise(exercise_type);
CREATE INDEX idx_exercise_difficulty ON public.exercise(difficulty_level);
CREATE INDEX idx_equipment_category ON public.equipment(category);
CREATE INDEX idx_gym_domain ON public.gym(domain);
CREATE INDEX idx_gym_active ON public.gym(is_active);

-- Tenant schema indexes (applied to each gym schema)
CREATE INDEX idx_user_role ON {gym_uuid}.users(role);
CREATE INDEX idx_user_active ON {gym_uuid}.users(is_active);
CREATE INDEX idx_workout_instance_user ON {gym_uuid}.custom_workout_instance(user_id);
CREATE INDEX idx_workout_exercise_instance ON {gym_uuid}.custom_workout_exercise(workout_instance_id);
```

### Query Optimization

- **Prepared statements** for frequently used queries
- **Connection pooling** to manage database connections
- **Schema-specific** connection routing for tenant operations
- **Materialized views** for complex aggregate queries

## 🔄 Migration Strategy

### Schema Creation

```sql
-- Create new tenant schema
CREATE SCHEMA IF NOT EXISTS {gym_uuid};

-- Apply all tenant tables to new schema
-- Copy structure from template schema

-- Set up permissions
GRANT USAGE ON SCHEMA {gym_uuid} TO app_user;
GRANT ALL ON ALL TABLES IN SCHEMA {gym_uuid} TO app_user;
```

### Data Migration

- **Backup strategies** per tenant schema
- **Schema versioning** for consistent updates
- **Migration scripts** applied across all tenant schemas
- **Rollback procedures** for safe deployments

---

This database design ensures complete tenant isolation while maximizing resource sharing and maintaining high performance across all gym operations.

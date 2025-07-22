# Tenant Schema Design

## Overview

Each gym gets its own schema (`{gym_uuid}.*`) with complete workout management capabilities.

## Table Structure

### 1. `{gym_uuid}.user` ✅ (Already exists)

- Gym staff and members

### 2. `{gym_uuid}.custom_exercise`

- Gym-specific exercises that extend the public library
- References public.exercise for global exercises

### 3. `{gym_uuid}.custom_equipment`

- Gym-specific equipment that extends the public library
- References public.equipment for global equipment

### 4. `{gym_uuid}.workout_template`

- Gym's own templates (similar to public.workout_templates)
- Extends the global template library with gym-specific templates

### 5. `{gym_uuid}.template_block`

- Blocks for gym-specific templates
- Similar to public.template_blocks but for gym templates

### 6. `{gym_uuid}.workout_instance`

- Actual workouts created from templates (public or gym-specific)
- The "filled" templates with actual exercises assigned

### 7. `{gym_uuid}.workout_exercise`

- Individual exercises within a workout instance
- Links to either public.exercise or {gym_uuid}.custom_exercises

### 8. `{gym_uuid}.member_workout`

- Member workout sessions (tracking when members do workouts)
- Links users to workout_instances with session data

## Data Flow

1. **Template Creation**: Gym creates template structure (blocks, counts)
2. **Template Filling**: AI/Manual selection fills blocks with actual exercises
3. **Workout Assignment**: Completed templates become workout instances
4. **Member Sessions**: Members execute workout instances, creating session records

## AI Integration Points

- Exercise selection for template blocks
- Workout customization based on member profile
- Progress tracking and adaptation

-- AthenAI Public (Admin) Schema for Supabase/PostgreSQL
-- Directly extracted from internal/database/public.go

CREATE TABLE IF NOT EXISTS public.gym (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.admin (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.exercise (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    synonyms TEXT[] NOT NULL,
    difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
    exercise_type TEXT NOT NULL CHECK (exercise_type IN ('strength', 'cardio', 'flexibility', 'balance', 'functional')),
    instructions TEXT NOT NULL,
    video_url TEXT,
    image_url TEXT,
    created_by UUID REFERENCES public.admin(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.muscular_group (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    body_part TEXT NOT NULL CHECK (body_part IN ('upper_body', 'lower_body', 'core', 'full_body')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    category TEXT NOT NULL CHECK (category IN ('free_weights', 'machines', 'cardio', 'accessories', 'bodyweight')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.exercise_muscular_group (
    exercise_id UUID NOT NULL REFERENCES public.exercise(id) ON DELETE CASCADE,
    muscular_group_id UUID NOT NULL REFERENCES public.muscular_group(id) ON DELETE RESTRICT,
    PRIMARY KEY (exercise_id, muscular_group_id)
);

CREATE TABLE IF NOT EXISTS public.exercise_equipment (
    exercise_id UUID NOT NULL REFERENCES public.exercise(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES public.equipment(id) ON DELETE RESTRICT,
    PRIMARY KEY (exercise_id, equipment_id)
);

CREATE TABLE IF NOT EXISTS public.refresh_token (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    token TEXT NOT NULL UNIQUE,
    user_type VARCHAR(50) NOT NULL CHECK (user_type IN ('platform_admin', 'tenant_user')),
    gym_id UUID REFERENCES public.gym(id),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, user_type, gym_id)
);

-- Table: workout_template
CREATE TABLE IF NOT EXISTS public.workout_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
    estimated_duration_minutes INTEGER,
    target_audience TEXT, -- e.g., 'weight_loss', 'muscle_building', 'endurance', 'general_fitness'
    created_by UUID REFERENCES public.admin(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_public BOOLEAN NOT NULL DEFAULT TRUE, -- If true, available to all gyms
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Table: template_block
CREATE TABLE IF NOT EXISTS public.template_block (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES public.workout_template(id) ON DELETE CASCADE,
    block_name TEXT NOT NULL, -- e.g., 'Pre-Warmup', 'Warmup', 'Main Block 1', 'Core', 'Cool Down'
    block_type TEXT NOT NULL CHECK (block_type IN ('warmup', 'main', 'core', 'cardio', 'cooldown', 'custom')),
    block_order INTEGER NOT NULL, -- Order of blocks in the template
    exercise_count INTEGER NOT NULL, -- Number of exercises for this block (e.g., 3 warmup exercises, 5 main exercises)
    estimated_duration_minutes INTEGER, -- Optional estimated time for this block
    instructions TEXT, -- Special instructions for this block type
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(template_id, block_order)
);

-- Indexes for refresh_token
CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON public.refresh_token(token);
CREATE INDEX IF NOT EXISTS idx_refresh_token_user ON public.refresh_token(user_id, user_type);
CREATE INDEX IF NOT EXISTS idx_refresh_token_expires ON public.refresh_token(expires_at);
CREATE INDEX IF NOT EXISTS idx_workout_template_active ON public.workout_template(is_active);
CREATE INDEX IF NOT EXISTS idx_workout_template_public ON public.workout_template(is_public);
CREATE INDEX IF NOT EXISTS idx_workout_template_difficulty ON public.workout_template(difficulty_level);
CREATE INDEX IF NOT EXISTS idx_template_block_template ON public.template_block(template_id);
CREATE INDEX IF NOT EXISTS idx_template_block_order ON public.template_block(template_id, block_order);

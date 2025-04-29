CREATE TABLE IF NOT EXISTS profiles (
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE PRIMARY KEY,
    description TEXT,
    age INT,
    location VARCHAR(255),
    avatar_url VARCHAR(512),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

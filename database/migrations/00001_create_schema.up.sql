CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE CHECK (email <>'') NOT NULL,
                                     phone_no TEXT NOT NULL,
                                     password TEXT NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                     updated_at TIMESTAMP WITH TIME ZONE,
                                     archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS role (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     role TEXT NOT NULL,
                                     user_id UUID REFERENCES users(id) NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                     updated_at TIMESTAMP WITH TIME ZONE,
                                     archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS notes (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    heading TEXT,
                                    note TEXT,
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE,
                                    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS user_note (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    user_id UUID REFERENCES users(id) NOT NULL,
                                    note_id UUID REFERENCES notes(id) NOT NULL,
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE,
                                    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS sessions (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        user_id UUID REFERENCES users(id) NOT NULL,
                                        start_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                        end_time TIMESTAMP WITH TIME ZONE,
                                        archived_at TIMESTAMP WITH TIME ZONE
);

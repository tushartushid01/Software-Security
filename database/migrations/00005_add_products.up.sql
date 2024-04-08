CREATE TABLE IF NOT EXISTS product (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     name TEXT,
                                     description TEXT,
                                     created_by UUID REFERENCES users(id) NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                     updated_at TIMESTAMP WITH TIME ZONE,
                                     archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS user_product (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         user_id UUID REFERENCES users(id) NOT NULL,
                                         product_id UUID REFERENCES product(id) NOT NULL,
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                         updated_at TIMESTAMP WITH TIME ZONE,
                                         archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS feedback (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       description TEXT,
                                       created_by UUID REFERENCES users(id) NOT NULL,
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       updated_at TIMESTAMP WITH TIME ZONE,
                                       archived_at TIMESTAMP WITH TIME ZONE
);
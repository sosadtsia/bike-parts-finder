-- Initialize the database schema for the Bike Parts Finder

-- Create tables if they don't exist
CREATE TABLE IF NOT EXISTS parts (
    id VARCHAR(36) PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    sub_category VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    msrp DECIMAL(10, 2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    in_stock BOOLEAN NOT NULL DEFAULT FALSE,
    rating DECIMAL(3, 2),
    num_reviews INTEGER,
    description TEXT,
    url VARCHAR(2048) NOT NULL,
    source VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS part_specs (
    id SERIAL PRIMARY KEY,
    part_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE CASCADE,
    CONSTRAINT unique_part_spec UNIQUE (part_id, name)
);

CREATE TABLE IF NOT EXISTS part_images (
    id SERIAL PRIMARY KEY,
    part_id VARCHAR(36) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scrape_requests (
    id VARCHAR(36) PRIMARY KEY,
    url VARCHAR(2048) NOT NULL,
    source VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_parts_brand ON parts(brand);
CREATE INDEX IF NOT EXISTS idx_parts_category ON parts(category);
CREATE INDEX IF NOT EXISTS idx_parts_sub_category ON parts(sub_category);
CREATE INDEX IF NOT EXISTS idx_parts_in_stock ON parts(in_stock);
CREATE INDEX IF NOT EXISTS idx_parts_price ON parts(price);
CREATE INDEX IF NOT EXISTS parts_search_idx ON parts USING GIN (
    to_tsvector('english', brand || ' ' || model || ' ' || COALESCE(description, ''))
);

-- Create a function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create a trigger to automatically update updated_at
CREATE TRIGGER update_parts_updated_at
BEFORE UPDATE ON parts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create a view for parts with their specifications
CREATE OR REPLACE VIEW parts_with_specs AS
SELECT
    p.*,
    jsonb_agg(jsonb_build_object('name', ps.name, 'value', ps.value)) AS specifications
FROM
    parts p
LEFT JOIN
    part_specs ps ON p.id = ps.part_id
GROUP BY
    p.id;

-- Create a view for parts with their images
CREATE OR REPLACE VIEW parts_with_images AS
SELECT
    p.*,
    jsonb_agg(pi.url ORDER BY pi.position) AS images
FROM
    parts p
LEFT JOIN
    part_images pi ON p.id = pi.part_id
GROUP BY
    p.id;

-- Create a view for complete parts information
CREATE OR REPLACE VIEW complete_parts AS
SELECT
    p.*,
    pws.specifications,
    pwi.images
FROM
    parts p
LEFT JOIN
    parts_with_specs pws ON p.id = pws.id
LEFT JOIN
    parts_with_images pwi ON p.id = pwi.id;

-- Insert some sample data
INSERT INTO parts (id, brand, model, category, sub_category, price, msrp, currency, in_stock, description, url, source)
VALUES
    ('1', 'Shimano', 'Deore XT M8100 12-speed Cassette', 'Drivetrain', 'Cassettes', 159.99, 179.99, 'USD', TRUE, 'High-performance 12-speed cassette for mountain bikes.', 'https://example.com/shimano-xt-cassette', 'Sample'),
    ('2', 'SRAM', 'GX Eagle Derailleur', 'Drivetrain', 'Derailleurs', 129.99, 149.99, 'USD', TRUE, '12-speed rear derailleur with Eagle technology.', 'https://example.com/sram-gx-eagle', 'Sample'),
    ('3', 'RockShox', 'Pike Ultimate Fork', 'Suspension', 'Forks', 899.99, 949.99, 'USD', FALSE, 'Premium trail fork with Charger 2.1 damper.', 'https://example.com/rockshox-pike', 'Sample')
ON CONFLICT (id) DO NOTHING;

INSERT INTO part_specs (part_id, name, value)
VALUES
    ('1', 'Speed', '12'),
    ('1', 'Range', '10-51T'),
    ('1', 'Weight', '460g'),
    ('2', 'Speed', '12'),
    ('2', 'Cage Length', 'Long'),
    ('2', 'Weight', '290g'),
    ('3', 'Travel', '140mm'),
    ('3', 'Wheel Size', '29"'),
    ('3', 'Axle', '15x110mm Boost')
ON CONFLICT ON CONSTRAINT unique_part_spec DO NOTHING;

INSERT INTO part_images (part_id, url, position)
VALUES
    ('1', 'https://example.com/images/shimano-xt-cassette-1.jpg', 0),
    ('1', 'https://example.com/images/shimano-xt-cassette-2.jpg', 1),
    ('2', 'https://example.com/images/sram-gx-eagle-1.jpg', 0),
    ('3', 'https://example.com/images/rockshox-pike-1.jpg', 0),
    ('3', 'https://example.com/images/rockshox-pike-2.jpg', 1),
    ('3', 'https://example.com/images/rockshox-pike-3.jpg', 2);

-- Grant permissions to the postgres user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- Notify that the database has been initialized
DO $$
BEGIN
    RAISE NOTICE 'Database initialization complete';
END $$;

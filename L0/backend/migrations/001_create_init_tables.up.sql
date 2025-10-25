--CREATION SCRIPTS
-- 1.New schema
CREATE SCHEMA IF NOT EXISTS L0;

-- 2.New tables
-- Create deliveries table
CREATE TABLE IF NOT EXISTS L0.Deliveries (
delivery_id INT PRIMARY KEY,
order_name VARCHAR(100),
phone VARCHAR(100),
zip VARCHAR(100),
city VARCHAR(100),
address VARCHAR(1000) NOT NULL,
region VARCHAR(100),
email VARCHAR(100) NOT NULL
);

-- Create orders table
CREATE TABLE IF NOT EXISTS L0.Orders (
order_uid VARCHAR(100) PRIMARY KEY,
track_number VARCHAR(100) UNIQUE NOT NULL,
entry VARCHAR(100) NOT NULL,
delivery_id INT,
FOREIGN KEY(delivery_id) REFERENCES L0.Deliveries(delivery_id) ON DELETE SET NULL,
locale VARCHAR(20) NOT NULL,
internal_signature VARCHAR(100),
customer_id VARCHAR(100) NOT NULL,
delivery_service VARCHAR(100),
shard_key VARCHAR(20),
sm_id INT,
date_created TIMESTAMP DEFAULT NOW(),
oof_shard VARCHAR(20)
);

-- Create payments table
CREATE TABLE IF NOT EXISTS L0.Payments (
pay_transaction VARCHAR(100) PRIMARY KEY,
FOREIGN KEY(pay_transaction) REFERENCES L0.Orders(order_uid) ON DELETE CASCADE,
request_id VARCHAR(100),
currency VARCHAR(20) NOT NULL,
provider VARCHAR(100) NOT NULL,
amount DECIMAL NOT NULL,
payment_dt INT NOT NULL,
bank VARCHAR(100) NOT NULL,
delivery_cost DECIMAL NOT NULL,
goods_total DECIMAL NOT NULL,
custom_fee DECIMAL NOT NULL
);

-- Create orderItems table
CREATE TABLE IF NOT EXISTS L0.OrderItems (
chrt_id SERIAL PRIMARY KEY,
track_number VARCHAR(100) NOT NULL,
FOREIGN KEY(track_number) REFERENCES L0.Orders(track_number) ON DELETE CASCADE,
price DECIMAL NOT NULL,
rid VARCHAR(100) NOT NULL,
item_name VARCHAR(100) NOT NULL,
sale INT NOT NULL,
item_size VARCHAR(20) NOT NULL,
totaL_PRICE DECIMAL NOT NULL,
nm_id INT NOT NULL,
brand VARCHAR(100) NOT NULL,
status INT NOT NULL
);

-- 3. Create user
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'testUser') THEN
        CREATE USER testUser WITH PASSWORD 'password';
    END IF;
END
$$

-- 4. Grant ownership to testUser
GRANT USAGE ON SCHEMA L0 TO testUser;

-- Grant privileges on all existing tables in the schema
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA L0 TO testUser;
GRANT TRUNCATE ON ALL TABLES IN SCHEMA L0 TO testUser;

-- Grant privileges on all existing sequences in the schema
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA L0 TO testUser;

-- Set default privileges for future tables created by the current user in the schema
ALTER DEFAULT PRIVILEGES IN SCHEMA L0
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO testUser;

-- Set default privileges for future sequences created by the current user in the schema
ALTER DEFAULT PRIVILEGES IN SCHEMA L0
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO testUser;
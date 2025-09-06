--DELETION SCRIPTS
DROP TABLE Deliveries CASCADE;
DROP TABLE Orders CASCADE;
DROP TABLE Payments CASCADE;
DROP TABLE OrderItems CASCADE;


--CREATION SCRIPTS
-- Create deliveries table
CREATE TABLE Deliveries (
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
CREATE TABLE Orders (
order_uid VARCHAR(100) PRIMARY KEY,
track_number VARCHAR(100) UNIQUE NOT NULL,
entry VARCHAR(100) NOT NULL,
delivery_id INT,
FOREIGN KEY(delivery_id) REFERENCES Deliveries(delivery_id) ON DELETE SET NULL,
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
CREATE TABLE Payments (
pay_transaction VARCHAR(100) PRIMARY KEY,
FOREIGN KEY(pay_transaction) REFERENCES Orders(order_uid) ON DELETE CASCADE,
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
CREATE TABLE OrderItems (
chrt_id SERIAL PRIMARY KEY,
track_number VARCHAR(100) NOT NULL,
FOREIGN KEY(track_number) REFERENCES Orders(track_number) ON DELETE CASCADE,
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
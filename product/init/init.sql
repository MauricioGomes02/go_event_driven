CREATE TABLE IF NOT EXISTS products (
    product_id CHAR(36) PRIMARY KEY NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description VARCHAR(255) NOT NULL,
    product_amount DECIMAL(10,2) NOT NULL,
    product_created_at DATETIME NOT NULL,
    product_active BIT NOT NULL
)
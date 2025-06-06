CREATE TABLE IF NOT EXISTS products (
    product_id CHAR(36) PRIMARY KEY NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description VARCHAR(255) NOT NULL,
    product_amount DECIMAL(10,2) NOT NULL,
    product_created_at DATETIME NOT NULL,
    product_active BIT NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox_events (
    id CHAR(36) PRIMARY KEY NOT NULL,               
    aggregate_id CHAR(36) NOT NULL,                 
    aggregate_type VARCHAR(50) NOT NULL,            
    event_type VARCHAR(100) NOT NULL,               
    payload JSON NOT NULL,                          
    status ENUM('pending', 'sent', 'error') NOT NULL DEFAULT 'pending',
    retries INT NOT NULL DEFAULT 0,
    error_message VARCHAR(1024) NULL,                             
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    sent_at DATETIME(6) NULL
);
-- Create the database if it doesn't exist
CREATE DATABASE device;

-- Connect to the database
\c device;

-- Create the device table
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    device_model_id INTEGER NOT NULL,
    house_id INTEGER NOT NULL,
    serial_number VARCHAR(50),
    name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create the device table
CREATE TABLE IF NOT EXISTS devices_attributes (
    id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL REFERENCES devices(id),
    key VARCHAR(100) NOT NULL,
    value VARCHAR(20)
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_devices_device_model_id ON devices(device_model_id);
CREATE INDEX IF NOT EXISTS idx_devices_house_id ON devices(house_id);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_attributes_device_id ON devices_attributes(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_attributes_key ON devices_attributes(key);
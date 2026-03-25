-- migrate:up

-- Dealerships
INSERT INTO dealerships (id, uuid, name, timezone, is_weekend_open)
OVERRIDING SYSTEM VALUE VALUES
(1, 'd1000000-0000-4000-a000-000000000001', 'Downtown Auto', 'Asia/Ho_Chi_Minh', true),
(2, 'd2000000-0000-4000-a000-000000000002', 'Uptown Motors', 'Asia/Ho_Chi_Minh', false);

SELECT setval(pg_get_serial_sequence('dealerships', 'id'), (SELECT MAX(id) FROM dealerships));

-- Service Bays
INSERT INTO service_bays (id, uuid, dealership_id, name)
OVERRIDING SYSTEM VALUE VALUES
(1, 'b1000000-0000-4000-a000-000000000001', 1, 'Bay 1'),
(2, 'b2000000-0000-4000-a000-000000000002', 1, 'Bay 2'),
(3, 'b3000000-0000-4000-a000-000000000003', 2, 'Bay 1');

SELECT setval(pg_get_serial_sequence('service_bays', 'id'), (SELECT MAX(id) FROM service_bays));

-- Service Types
INSERT INTO service_types (id, uuid, name, duration_minutes)
OVERRIDING SYSTEM VALUE VALUES
(1, 's1000000-0000-4000-a000-000000000001', 'Oil Change', 60),
(2, 's2000000-0000-4000-a000-000000000002', 'Tire Rotation', 30),
(3, 's3000000-0000-4000-a000-000000000003', 'Full Service', 180);

SELECT setval(pg_get_serial_sequence('service_types', 'id'), (SELECT MAX(id) FROM service_types));

-- Technicians
INSERT INTO technicians (id, uuid, dealership_id, name)
OVERRIDING SYSTEM VALUE VALUES
(1, 't1000000-0000-4000-a000-000000000001', 1, 'Alice'),
(2, 't2000000-0000-4000-a000-000000000002', 1, 'Bob'),
(3, 't3000000-0000-4000-a000-000000000003', 2, 'Charlie');

SELECT setval(pg_get_serial_sequence('technicians', 'id'), (SELECT MAX(id) FROM technicians));

-- Technician Service Types
INSERT INTO technician_service_types (technician_id, service_type_id) VALUES
(1, 1), -- Alice does Oil Change
(1, 2), -- Alice does Tire Rotation
(2, 3), -- Bob does Full Service
(3, 1), -- Charlie does Oil Change
(3, 3); -- Charlie does Full Service

-- Customers
INSERT INTO customers (id, uuid, name, email, phone)
OVERRIDING SYSTEM VALUE VALUES
(1, 'c1000000-0000-4000-a000-000000000001', 'John Doe', 'john@example.com', '123456789'),
(2, 'c2000000-0000-4000-a000-000000000002', 'Jane Smith', 'jane@example.com', '987654321');

SELECT setval(pg_get_serial_sequence('customers', 'id'), (SELECT MAX(id) FROM customers));

-- Vehicles
INSERT INTO vehicles (id, uuid, customer_id, name)
OVERRIDING SYSTEM VALUE VALUES
(1, 'v1000000-0000-4000-a000-000000000001', 1, 'Toyota Camry'),
(2, 'v2000000-0000-4000-a000-000000000002', 2, 'Honda Civic');

SELECT setval(pg_get_serial_sequence('vehicles', 'id'), (SELECT MAX(id) FROM vehicles));

-- migrate:down

DELETE FROM vehicles;
DELETE FROM customers;
DELETE FROM technician_service_types;
DELETE FROM technicians;
DELETE FROM service_types;
DELETE FROM service_bays;
DELETE FROM dealerships;

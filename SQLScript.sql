-- Create database
CREATE DATABASE IF NOT EXISTS CarSharingSystem;
USE CarSharingSystem;

-- Users Table
CREATE TABLE Users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    membership_id ENUM('Basic', 'Premium', 'VIP') DEFAULT 'Basic',
    verification_code VARCHAR(6),
    verified BOOLEAN DEFAULT FALSE,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_memberships FOREIGN KEY (membership_id) REFERENCES MembershipTiers(membership_id) 
);

-- Create RentalHistory table
CREATE TABLE RentalHistory (
    rental_id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL, 
    model VARCHAR(255) NOT NULL,
	rental_start DATETIME NOT NULL,
    rental_end DATETIME NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_history FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_vehicle_history FOREIGN KEY (vehicle_id) REFERENCES Vehicles(vehicle_id) ON DELETE CASCADE
);

-- Vehicles Table
CREATE TABLE Vehicles (
    vehicle_id INT AUTO_INCREMENT PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reservations Table
CREATE TABLE Reservations (
    reservation_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('active', 'cancelled', 'completed') DEFAULT 'active',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_vehicle FOREIGN KEY (vehicle_id) REFERENCES Vehicles(vehicle_id) ON DELETE CASCADE
);

-- Billing Table
CREATE TABLE Billing (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    initial_amount DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0.00,
    final_amount DECIMAL(10, 2) AS (initial_amount * (1 - discount / 100)) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reservation FOREIGN KEY (reservation_id) REFERENCES Reservations(reservation_id) ON DELETE CASCADE
);

CREATE TABLE MembershipTiers (
    membership_id ENUM('Basic', 'Premium', 'VIP') NOT NULL PRIMARY KEY,
    discount DECIMAL(5, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE Promotions (
    promotion_id VARCHAR(5) PRIMARY KEY NOT NULL,
    promotion_code VARCHAR(20) NOT NULL,
    discount DECIMAL(5, 2) NOT NULL
);

-- Insert Mock Data into Promotions Table
INSERT INTO Promotions (promotion_id, promotion_code, discount)
VALUES
    ('P001', 'HOLIDAY20', 20.00), -- 20% holiday discount
    ('P002', 'NEWYEAR10', 10.00); -- 10% New Year discount

-- Insert values into memberships table
INSERT INTO MembershipTiers (membership_id, booking_limit) 
VALUES 
    ('Basic', 5),  -- 5 bookings for Basic membership
    ('Premium', 10),  -- 10 bookings for Premium
    ('VIP', 20);  -- 20 bookings for VIP

-- Insert values into users
INSERT INTO Users (name, email, password, phone_number, membership_id, verification_code, verified)
VALUES (
    'John Doe', 
    'johndoe@example.com', 
    '$2a$10$E4C7hXOZCqQiNlsFbYg3a.nFDoNC7.Y3DAf5KntTSD2Sk68ZX4n0C', -- Example hashed password
    '+1234567890', 
    'Premium', 
    '123456', 
    FALSE
);

-- Dummy inserts for RentalHistory
INSERT INTO RentalHistory (rental_id, user_id, model, vehicle_id, rental_start, rental_end, total_amount)
VALUES
    ('R003', 1, 101, 'Toyota Corolla', '2024-12-03 10:00:00', '2024-12-03 12:00:00', 50.00),
    ('R004', 2, 102, 'Honda Civic', '2024-12-01 13:00:00', '2024-12-01 15:00:00', 40.00),
    ('R009', 1, 101, 'Toyota Corolla', '2024-12-03 10:00:00', '2024-12-03 12:00:00', 50.00);
-- Dummy inserts for Vehicles
INSERT INTO Vehicles (vehicle_id, model, location, status)
VALUES
    (101, 'Toyota Corolla', 'Location A'),
    (102, 'Honda Civic', 'Location B'),
    (103, 'Ford Fiesta', 'Location C'),
    (104, 'Toyota Corolla', 'Location B');
    
-- Dummy inserts for Reservations
INSERT INTO Reservations (user_id, vehicle_id, start_time, end_time, status, total_amount)
VALUES
	(1, 101, '2024-12-01 12:00:00', '2024-12-01 14:00:00', 'completed', 25.00),
    (1, 102, '2024-12-02 15:00:00', '2024-12-02 17:00:00', 'cancelled', 20.00),
    (2, 101, '2024-12-01 10:00:00', '2024-12-01 12:00:00', 'completed', 25.00),
    (2, 102, '2024-12-02 13:00:00', '2024-12-02 15:00:00', 'cancelled', 20.00);

-- Dummy inserts for Billing
INSERT INTO Billing (reservation_id, initial_amount, discount)
VALUES
    (14, 30.00, 10.00),
    (22, 65.00, 15.00);
-- for testing
SELECT COUNT(*) FROM Reservations WHERE vehicle_id = '101' AND start_time > '2024-12-10 17:00:10' AND end_time < '2024-12-10 17:20:00' AND status = 'active';
SELECT * FROM Reservations WHERE start_time > '2023-12-10 :00:10';
SELECT * FROM Reservations;

SELECT final_amount FROM billing WHERE reservation_id = 39
SELECT * FROM Users;
select * from vehicles
select * from reservations;
select * from rentalhistory;
select * from billing
select * from membershiptiers
select * from promotions
SELECT COUNT(*) FROM Reservations WHERE vehicle_id = 101 AND start_time <= 2024-12-10 17:00:00 AND end_time >= 2024-12-10 19:00:00 AND status = 'active'
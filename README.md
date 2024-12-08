# 🚗 **Car Sharing System**

## **Overview**

Car sharing system designed to enable user registration, vehicle reservation, billing.. The system follows a microservices architecture with distinct services handling user management, vehicle reservations, and billing processes.

---

## 📋 **Table of Contents**

1. [User Management](#user-management)  
   - User Registration and Authentication  
   - Membership Tiers  
   - User Profile Management  
2. [Vehicle Reservation System](#vehicle-reservation-system)  
   - Car Availability and Booking  
   - Booking Modification and Cancellation  
3. [Billing and Payment Processing](#billing-and-payment-processing)  
   - Tier-Based Pricing and Discounts  
   - Real-Time Billing Calculation  
   - Invoicing and Receipts  
4. [Microservices Architecture](#microservices-architecture)  
   - Service Decomposition  
   - Inter-Service Communication  
5. [Database Management](#database-management)  
   - Database Schema Design

---

## 👥 **User Management**

### **3.1.1 User Registration and Authentication**

The system implements secure user registration with email/phone verification and authentication. Passwords are securely encrypted to ensure the safety of user credentials.

**Key Features:**
- User registration with email/phone verification.
- Authentication with password hashing using bcrypt.

---

### **3.1.2 Membership Tiers**

Users can subscribe to different membership levels with varying benefits:

- **Basic Membership:** Standard rates.
- **Premium Membership:** Reduced rates (5% discount)
- **VIP Membership:** Reduced rates (10% discount)

---

### **3.1.3 User Profile Management**

Users can:

- Update personal details like contact information and name.  
- View their membership tier and status.  
- Track their rental history, including reservations made.

---

## 🚗 **Vehicle Reservation System**

### **3.2.1 Car Availability and Booking**

Users can:

- View available cars in real-time.  
- Book vehicles for specific date and time ranges.  

**User Flow:**  
1. Check available cars.  
2. Select time range for booking.  
3. Confirm reservation.

---

### **3.2.2 Booking Modification and Cancellation**

The system allows users to:

- Modify their existing reservations (time ranges or other preferences) within defined policies.  
- Cancel reservations, triggering automatic updates to vehicle availability.

The system ensures a smooth user experience by automatically managing updates.

---

## 💸 **Billing and Payment Processing**

### **3.3.1 Tier-Based Pricing and Discounts**

The system calculates billing based on:
- Membership level (Basic, Premium, VIP).  
- Rental duration by minute.  
- Applicable discounts through promotions and membership status.

Discounts are dynamically applied to ensure fair pricing for all users.

---

### **3.3.2 Real-Time Billing Calculation**

Before a user confirms their reservation, estimated costs are shown dynamically to ensure transparency.

**Features:**
- Estimate costs before confirming a reservation. 

---

### **3.3.3 Invoicing and Receipts**

Upon the successful completion of each rental:
- A detailed invoice is generated automatically.  
- Invoices are either emailed to the user or made available in their user profile.

This ensures that users have proper records for their bookings and payments.

---

## 🏗️ **Microservices Architecture**

### **3.4.1 Service Decomposition**

The system has been broken into microservices for scalability and maintainability:

1. **User Service:** Handles user registration, authentication, membership tiers, and user profile operations.  
2. **Vehicle Service:** Manages vehicle availability, car booking, and reservation management.  
3. **Billing Service:** Handles cost calculations, invoice generation, payment processing, and discount application.

Each service communicates independently, ensuring modularity.

---

### **3.4.2 Inter-Service Communication**

Communication between microservices is facilitated via **RESTful APIs**. All APIs have standardized endpoints with proper documentation to ensure smooth communication and clarity for future integrations.

---

## 🛢️ **Database Management**

### **3.5.1 Database Schema Design**

The system uses a **relational database** to store structured data such as:

- **Users:** Tracks user registration details, membership tier, and contact information.  
- **Vehicles:** Contains vehicle availability and location details.  
- **Reservations:** Manages user booking history and reservation times.  
- **Billing:** Tracks costs associated with a user’s reservation and payment history.  
- **Promotions:** Manages promotions/discounts applied to billing.

The database schema adheres to normalization to reduce redundancy and ensure data integrity.

---

## 📊 **Tech Stack**

The system is built using the following technologies:

- **Programming Language:** Go (Golang)  
- **Database:** MySQL  
- **Web Framework:** Standard Go net/http  
- **Microservices:** Decomposed architecture with independent service endpoints  
- **Authentication:** Secure user password encryption with bcrypt  
- **APIs:** RESTful API endpoints for service communication  

---

## **Architecture Diagram**
![test](https://github.com/user-attachments/assets/2130dc94-5223-4ffc-81e3-c041ac75e883)



-- Bảng Users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       name VARCHAR(50) NOT NULL,
                       password TEXT NOT NULL,
                       phone_number VARCHAR(11) UNIQUE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Roles
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT
);

-- Bảng Permissions
CREATE TABLE permissions (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(100) UNIQUE NOT NULL,
                             description TEXT
);

-- Bảng RolePermissions
CREATE TABLE rolepermissions (
                                 role_id INT NOT NULL,
                                 permission_id INT NOT NULL,
                                 PRIMARY KEY (role_id, permission_id)
);

-- Bảng UserRoles
CREATE TABLE userroles (
                           user_id INT NOT NULL,
                           role_id INT NOT NULL,
                           PRIMARY KEY (user_id, role_id)
);

-- Bảng Movies
CREATE TABLE movies (
                        id SERIAL PRIMARY KEY,
                        title VARCHAR(100) UNIQUE NOT NULL,
                        description VARCHAR(255) NOT NULL,
                        duration INT NOT NULL,
                        cover_url TEXT,
                        background_url TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Theaters
CREATE TABLE theaters (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(10) UNIQUE NOT NULL
);

-- Bảng Seats
CREATE TABLE seats (
                       id SERIAL PRIMARY KEY,
                       type VARCHAR(255) UNIQUE NOT NULL,
                       description TEXT,
                       price DOUBLE PRECISION NOT NULL,
                       number INT NOT NULL,
                       event_id INT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Payment
CREATE TABLE payment (
                         id SERIAL PRIMARY KEY,
                         user_id INT NOT NULL,
                         total_amount DOUBLE PRECISION NOT NULL,
                         status VARCHAR(20) NOT NULL DEFAULT 'pending',
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng PaymentDetails
CREATE TABLE paymentdetails (
                                id SERIAL PRIMARY KEY,
                                payment_id INT NOT NULL,
                                seat_id INT NOT NULL,
                                quantity INT NOT NULL DEFAULT 1,
                                total_price DOUBLE PRECISION NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng MoviesCategories
CREATE TABLE moviescategories (
                                  id SERIAL PRIMARY KEY,
                                  name VARCHAR(100) UNIQUE NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng MovieCategoryMappings
CREATE TABLE moviecategorymappings (
                                       movie_id INT NOT NULL,
                                       category_id INT NOT NULL,
                                       PRIMARY KEY (movie_id, category_id)
);

-- Khóa ngoại
ALTER TABLE rolepermissions ADD FOREIGN KEY (role_id) REFERENCES roles (id);
ALTER TABLE rolepermissions ADD FOREIGN KEY (permission_id) REFERENCES permissions (id);
ALTER TABLE userroles ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE userroles ADD FOREIGN KEY (role_id) REFERENCES roles (id);
ALTER TABLE seats ADD FOREIGN KEY (event_id) REFERENCES theaters (id);
ALTER TABLE payment ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE paymentdetails ADD FOREIGN KEY (payment_id) REFERENCES payment (id);
ALTER TABLE paymentdetails ADD FOREIGN KEY (seat_id) REFERENCES seats (id);
ALTER TABLE moviecategorymappings ADD FOREIGN KEY (movie_id) REFERENCES movies (id);
ALTER TABLE moviecategorymappings ADD FOREIGN KEY (category_id) REFERENCES moviescategories (id);

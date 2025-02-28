-- 1. Tạo bảng Users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       name VARCHAR(50) NOT NULL,
                       password TEXT NOT NULL,
                       phone_number VARCHAR(11) UNIQUE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tạo bảng Roles
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT
);

-- 3. Tạo bảng Permissions
CREATE TABLE permissions (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(100) UNIQUE NOT NULL,
                             description TEXT
);

-- 4. Tạo bảng RolePermissions
CREATE TABLE rolepermissions (
                                 role_id INT NOT NULL,
                                 permission_id INT NOT NULL,
                                 PRIMARY KEY (role_id, permission_id)
);

-- 5. Tạo bảng UserRoles
CREATE TABLE userroles (
                           user_id INT NOT NULL,
                           role_id INT NOT NULL,
                           PRIMARY KEY (user_id, role_id)
);

-- 6. Tạo bảng Movies
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

-- 7. Tạo bảng Theaters
CREATE TABLE theaters (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(10) UNIQUE NOT NULL
);

-- 8. Tạo bảng Seat Types (lưu trữ hạng ghế, giá vé và mô tả)
CREATE TABLE seat_types (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(50) NOT NULL,  -- Ví dụ: 'VIP', 'Thường'
                            description TEXT,
                            price DOUBLE PRECISION NOT NULL,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. Tạo bảng Seats (chỉ giữ lại id, number, theater_id, seat_type_id, created_at và updated_at)
CREATE TABLE seats (
                       id SERIAL PRIMARY KEY,
                       number VARCHAR(50) NOT NULL,  -- Mã ghế, ví dụ: "A1", "B12", ...
                       theater_id INT NOT NULL,
                       seat_type_id INT,             -- Ban đầu có thể để NULL, sau này admin cập nhật
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Tạo bảng Payment
CREATE TABLE payment (
                         id SERIAL PRIMARY KEY,
                         user_id INT NOT NULL,
                         total_amount DOUBLE PRECISION NOT NULL,
                         status VARCHAR(20) NOT NULL DEFAULT 'pending',
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Tạo bảng PaymentDetails
CREATE TABLE paymentdetails (
                                id SERIAL PRIMARY KEY,
                                payment_id INT NOT NULL,
                                seat_id INT NOT NULL,
                                quantity INT NOT NULL DEFAULT 1,
                                showtime_id INT NOT NULL,
                                total_price DOUBLE PRECISION NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. Tạo bảng MoviesCategories
CREATE TABLE moviescategories (
                                  id SERIAL PRIMARY KEY,
                                  name VARCHAR(100) UNIQUE NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 13. Tạo bảng MovieCategoryMappings
CREATE TABLE moviecategorymappings (
                                       movie_id INT NOT NULL,
                                       category_id INT NOT NULL,
                                       PRIMARY KEY (movie_id, category_id)
);

-- 14. Tạo bảng Showtimes
CREATE TABLE IF NOT EXISTS showtimes (
                                         id SERIAL PRIMARY KEY,
                                         movie_id INT NOT NULL,
                                         theater_id INT NOT NULL,
                                         start_time TIMESTAMP NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         FOREIGN KEY (movie_id) REFERENCES movies(id),
                                         FOREIGN KEY (theater_id) REFERENCES theaters(id)
);



-- 16. Thiết lập khóa ngoại

-- Khóa ngoại cho RolePermissions
ALTER TABLE rolepermissions ADD FOREIGN KEY (role_id) REFERENCES roles(id);
ALTER TABLE rolepermissions ADD FOREIGN KEY (permission_id) REFERENCES permissions(id);

-- Khóa ngoại cho UserRoles
ALTER TABLE userroles ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE userroles ADD FOREIGN KEY (role_id) REFERENCES roles(id);

-- Khóa ngoại cho Seats
ALTER TABLE seats ADD FOREIGN KEY (theater_id) REFERENCES theaters(id);
ALTER TABLE seats ADD FOREIGN KEY (seat_type_id) REFERENCES seat_types(id);

-- Khóa ngoại cho Payment
ALTER TABLE payment ADD FOREIGN KEY (user_id) REFERENCES users(id);

-- Khóa ngoại cho PaymentDetails
ALTER TABLE paymentdetails ADD FOREIGN KEY (payment_id) REFERENCES payment(id);
ALTER TABLE paymentdetails ADD FOREIGN KEY (seat_id) REFERENCES seats(id);
ALTER TABLE paymentdetails
    ADD CONSTRAINT fk_paymentdetails_showtime
        FOREIGN KEY (showtime_id)
            REFERENCES showtimes(id);

-- Khóa ngoại cho MovieCategoryMappings
ALTER TABLE moviecategorymappings ADD FOREIGN KEY (movie_id) REFERENCES movies(id);
ALTER TABLE moviecategorymappings ADD FOREIGN KEY (category_id) REFERENCES moviescategories(id);


-- 1. Insert basic permissions
INSERT INTO permissions (name, description) VALUES
                                                ('view_movies', 'Allows viewing movies'),
                                                ('purchase_tickets', 'Allows ticket purchase'),
                                                ('manage_movies', 'Allows managing movie information'),
                                                ('manage_users', 'Allows managing users'),
                                                ('manage_roles', 'Allows managing roles and permissions');

-- 2. Insert basic roles
INSERT INTO roles (name, description) VALUES
                                          ('customer', 'Basic customer role'),
                                          ('staff', 'Staff role with limited permissions'),
                                          ('admin', 'Administrator role with highest permissions');

-- 3. Assign permissions cho từng role

-- Customer: chỉ được phép xem phim và mua vé
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'customer'
  AND p.name IN ('view_movies', 'purchase_tickets');

-- Staff: được phép xem phim, mua vé và quản lý phim
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff'
  AND p.name IN ('view_movies', 'purchase_tickets', 'manage_movies');

-- Admin: gán tất cả các quyền
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin';

-- 4. Insert super user (quản trị viên đặc biệt)
INSERT INTO users (email, name, password, phone_number) VALUES
    ('superuser@example.com', 'Super User', 'supersecret', '01234567890');

-- 5. Gán super user vào vai trò admin
INSERT INTO userroles (user_id, role_id)
VALUES (
           (SELECT id FROM users WHERE email = 'superuser@example.com'),
           (SELECT id FROM roles WHERE name = 'admin')
       );

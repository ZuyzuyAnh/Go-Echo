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

-- 3. Assign permissions to each role

-- Customer: only allowed to view movies and purchase tickets
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'customer'
  AND p.name IN ('view_movies', 'purchase_tickets');

-- Staff: allowed to view movies, purchase tickets, and manage movies
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff'
  AND p.name IN ('view_movies', 'purchase_tickets', 'manage_movies');

-- Administrator: assign all permissions
INSERT INTO rolepermissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin';

-- 4. Insert super user (special administrator)
INSERT INTO users (email, name, password, phone_number) VALUES
    ('superuser@example.com', 'Super User', 'supersecret', '01234567890');

-- 5. Assign super user to admin role
INSERT INTO userroles (user_id, role_id)
VALUES (
           (SELECT id FROM users WHERE email = 'superuser@example.com'),
           (SELECT id FROM roles WHERE name = 'admin')
       );

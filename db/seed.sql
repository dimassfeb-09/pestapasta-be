-- Seed data for pestapasta-be (PostgreSQL)
-- Run after AutoMigrate has created the tables.

-- Categories
INSERT INTO categories (id, category_name, description) VALUES
(1, 'Pasta', 'Berbagai menu pasta dengan saus pilihan'),
(2, 'Pizza', 'Pizza panggang dengan topping premium'),
(3, 'Drinks', 'Minuman segar pendamping makanan'),
(4, 'Dessert', 'Hidangan penutup manis'),
(5, 'Sides', 'Menu pendamping / appetizer')
ON CONFLICT (id) DO NOTHING;

-- Payment Methods
INSERT INTO payment_methods (id, name, account_number, account_name, code) VALUES
(1, 'BCA', '1234567890', 'Pesta Pasta', 'bca'),
(2, 'QRIS', '0987654321', 'Pesta Pasta', 'qris')
ON CONFLICT (id) DO NOTHING;

-- Menus / Products (min 10)
INSERT INTO menus (id, name, price, description, category_id, image_url, rating, is_available) VALUES
(1, 'Spaghetti Bolognese', 45000, 'Spaghetti dengan saus bolognese daging sapi cincang', 1, 'https://picsum.photos/seed/pestapasta1/400/300', 5, true),
(2, 'Carbonara', 50000, 'Pasta carbonara creamy dengan bacon dan parmesan', 1, 'https://picsum.photos/seed/pestapasta2/400/300', 5, true),
(3, 'Aglio Olio', 40000, 'Spaghetti aglio olio dengan cabai dan olive oil', 1, 'https://picsum.photos/seed/pestapasta3/400/300', 4, true),
(4, 'Lasagna', 60000, 'Lasagna panggang dengan saus bechamel dan daging', 1, 'https://picsum.photos/seed/pestapasta4/400/300', 5, true),
(5, 'Fettuccine Alfredo', 52000, 'Fettuccine dengan saus alfredo creamy', 1, 'https://picsum.photos/seed/pestapasta5/400/300', 4, true),
(6, 'Margherita Pizza', 75000, 'Pizza margherita dengan mozzarella dan basil', 2, 'https://picsum.photos/seed/pestapasta6/400/300', 5, true),
(7, 'Pepperoni Pizza', 85000, 'Pizza pepperoni dengan keju mozzarella', 2, 'https://picsum.photos/seed/pestapasta7/400/300', 5, true),
(8, 'Iced Lemon Tea', 15000, 'Teh lemon dingin segar', 3, 'https://picsum.photos/seed/pestapasta8/400/300', 4, true),
(9, 'Mineral Water', 8000, 'Air mineral 600ml', 3, 'https://picsum.photos/seed/pestapasta9/400/300', 3, true),
(10, 'Tiramisu', 35000, 'Dessert tiramisu khas Italia', 4, 'https://picsum.photos/seed/pestapasta10/400/300', 5, true),
(11, 'Garlic Bread', 20000, 'Roti bawang panggang crispy', 5, 'https://picsum.photos/seed/pestapasta11/400/300', 4, true),
(12, 'Caesar Salad', 30000, 'Salad caesar dengan crouton dan parmesan', 5, 'https://picsum.photos/seed/pestapasta12/400/300', 4, true)
ON CONFLICT (id) DO NOTHING;

-- Optional: reset sequences so future inserts get correct IDs
SELECT setval(pg_get_serial_sequence('categories', 'id'), COALESCE((SELECT MAX(id) FROM categories), 1));
SELECT setval(pg_get_serial_sequence('payment_methods', 'id'), COALESCE((SELECT MAX(id) FROM payment_methods), 1));
SELECT setval(pg_get_serial_sequence('menus', 'id'), COALESCE((SELECT MAX(id) FROM menus), 1));

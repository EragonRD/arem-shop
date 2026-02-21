-- 1. Create the second shop (completely empty)
INSERT INTO shops (id, name, active, whatsapp_number) 
VALUES ('22222222-2222-2222-2222-222222222222', 'Shop 2', true, '+212600000002')
ON CONFLICT (id) DO NOTHING;

-- 2. Create the SuperAdmin for Shop 2
-- We use a pre-calculated bcrypt hash for the password "Password456!" (cost=12)
INSERT INTO users (id, name, email, password, role, shop_id)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    'Owner Shop 2',
    'owner2@shopdemo.com',
    '$2a$12$e4c2jFXY3Lq2sBcR5kHb3Oqk5N7s4k5N7s4k5N7s4k5N7s4k5N7s', -- This is a mocked hash. Postgres pgcrypto is better if installed, but for simplicity we rely on a pre-generated bcrypt.
    'SuperAdmin',
    '22222222-2222-2222-2222-222222222222'
) ON CONFLICT (id) DO NOTHING;

-- Update the hashed password accurately using pgcrypto if the extension is available
CREATE EXTENSION IF NOT EXISTS pgcrypto;
UPDATE users 
SET password = crypt('Password456!', gen_salt('bf', 12)) 
WHERE email = 'owner2@shopdemo.com';

-- 3. We also need to seed the basic categories for Shop 2 so it can use the Product form right away!
INSERT INTO categories (name, shop_id) VALUES
  ('Laptops', '22222222-2222-2222-2222-222222222222'),
  ('Desktop PCs', '22222222-2222-2222-2222-222222222222'),
  ('Smartphones', '22222222-2222-2222-2222-222222222222'),
  ('Tablets', '22222222-2222-2222-2222-222222222222'),
  ('Accessories', '22222222-2222-2222-2222-222222222222'),
  ('Components', '22222222-2222-2222-2222-222222222222'),
  ('Networking', '22222222-2222-2222-2222-222222222222'),
  ('Printers', '22222222-2222-2222-2222-222222222222'),
  ('Monitors', '22222222-2222-2222-2222-222222222222'),
  ('Audio', '22222222-2222-2222-2222-222222222222');

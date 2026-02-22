-- ──────────────────────────────────────────────────────────────
-- migrations/000003_add_categories.up.sql
-- Create dedicated categories table and migrate existing data
-- ──────────────────────────────────────────────────────────────

-- 1. Create the categories table
CREATE TABLE categories (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name           VARCHAR(80) NOT NULL,
  shop_id        UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_categories_name_shop UNIQUE (name, shop_id)
);

CREATE INDEX idx_categories_shop_id ON categories (shop_id);

-- 2. Seed basic electronics categories for ALL existing shops
DO $$
DECLARE
  current_shop_id UUID;
BEGIN
  FOR current_shop_id IN SELECT id FROM shops LOOP
    INSERT INTO categories (name, shop_id) VALUES
      ('Smartphones & Mobiles', current_shop_id),
      ('Tablettes & Accessoires', current_shop_id),
      ('Ordinateurs Portables', current_shop_id),
      ('Composants PC', current_shop_id),
      ('Écrans & Moniteurs', current_shop_id),
      ('Audio & Casques', current_shop_id),
      ('Consoles & Gaming', current_shop_id),
      ('Objets Connectés', current_shop_id),
      ('Réseau & Stockage', current_shop_id),
      ('Accessoires divers', current_shop_id)
    ON CONFLICT DO NOTHING;
  END LOOP;
END $$;

-- 3. Add temporary category_id to products
ALTER TABLE products ADD COLUMN category_id UUID REFERENCES categories(id) ON DELETE RESTRICT;

-- 4. Migrate existing strings to UUIDs
-- Any existing category in `products` that doesn't exist in `categories` will be created first.
DO $$
DECLARE
  prod RECORD;
  cat_id UUID;
BEGIN
  FOR prod IN SELECT id, category, shop_id FROM products LOOP
    -- Try to find existing category by string match
    SELECT id INTO cat_id FROM categories WHERE name = prod.category AND shop_id = prod.shop_id LIMIT 1;
    
    -- If string doesn't exist in our seeded list, insert it and get its UUID
    IF cat_id IS NULL THEN
      INSERT INTO categories (name, shop_id) VALUES (prod.category, prod.shop_id) RETURNING id INTO cat_id;
    END IF;
    
    -- Update the product with the mapped category_id
    UPDATE products SET category_id = cat_id WHERE id = prod.id;
  END LOOP;
END $$;

-- 5. Make category_id NOT NULL
ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;

-- 6. Drop the old string index and column
DROP INDEX IF EXISTS idx_products_shop_category;
ALTER TABLE products DROP COLUMN category;

-- 7. Add the new UUID index
CREATE INDEX idx_products_shop_category_id ON products (shop_id, category_id);

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('SuperAdmin', 'Admin');
CREATE TYPE transaction_type AS ENUM ('Sale', 'Expense', 'Withdrawal');

CREATE TABLE shops (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(120) NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  whatsapp_number VARCHAR(20) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_shops_whatsapp_format
    CHECK (whatsapp_number ~ '^\+[1-9][0-9]{7,14}$')
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(120) NOT NULL,
  email VARCHAR(254) NOT NULL,
  password TEXT NOT NULL,
  role user_role NOT NULL,
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_users_shop_email ON users (shop_id, LOWER(email));
CREATE INDEX idx_users_shop_id ON users (shop_id);

CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(160) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  category VARCHAR(80) NOT NULL,
  purchase_price NUMERIC(12,2) NOT NULL CHECK (purchase_price >= 0),
  selling_price NUMERIC(12,2) NOT NULL CHECK (selling_price >= 0),
  stock INT NOT NULL CHECK (stock >= 0),
  image_url TEXT NOT NULL DEFAULT '',
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_products_id_shop UNIQUE (id, shop_id)
);

CREATE INDEX idx_products_shop_id ON products (shop_id);
CREATE INDEX idx_products_shop_category ON products (shop_id, category);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type transaction_type NOT NULL,
  product_id UUID NULL,
  quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_transactions_product_shop
    FOREIGN KEY (product_id, shop_id)
    REFERENCES products (id, shop_id)
    ON DELETE RESTRICT,
  CONSTRAINT chk_transactions_rules
    CHECK (
      (type = 'Sale' AND product_id IS NOT NULL AND quantity > 0)
      OR
      (type IN ('Expense', 'Withdrawal') AND product_id IS NULL AND quantity = 0)
    )
);

CREATE INDEX idx_transactions_shop_created ON transactions (shop_id, created_at DESC);
CREATE INDEX idx_transactions_shop_type ON transactions (shop_id, type);

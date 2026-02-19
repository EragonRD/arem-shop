-- ──────────────────────────────────────────────────────────────
-- migrations/000002_seed_demo.sql
-- Seed de démonstration : Shop + SuperAdmin pour tester l'API
-- Ce seed s'exécute automatiquement dans Docker via initdb.d
-- ──────────────────────────────────────────────────────────────

-- Shop de démonstration
INSERT INTO shops (id, name, active, whatsapp_number)
VALUES ('11111111-1111-1111-1111-111111111111', 'Shop Demo', true, '+212600000000')
ON CONFLICT (id) DO NOTHING;

-- SuperAdmin (mot de passe : ChangeMe123!)
-- Hash bcrypt généré avec cost=12
INSERT INTO users (name, email, password, role, shop_id)
VALUES (
  'Owner Demo',
  'owner@shopdemo.com',
  '$2a$12$ObVk8.mVXqVxv/84Zbse0Od2BlNQ4OJDUSw/3OZ2ZRH23/VGd7tFy',
  'SuperAdmin',
  '11111111-1111-1111-1111-111111111111'
)
ON CONFLICT DO NOTHING;

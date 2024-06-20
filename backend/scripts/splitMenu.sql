ALTER TABLE tovars DROP COLUMN cost;
ALTER TABLE tovars DROP COLUMN profit;
ALTER TABLE tovars DROP COLUMN margin;
ALTER TABLE ingredients DROP COLUMN cost;
ALTER TABLE tech_carts DROP COLUMN cost;
ALTER TABLE tech_carts DROP COLUMN profit;
ALTER TABLE tech_carts DROP COLUMN margin;
ALTER TABLE ingredient_nabors DROP COLUMN name;
ALTER TABLE ingredient_nabors DROP COLUMN image;

SELECT setval('public.tovar_masters_id_seq',125, true);
SELECT setval('public.tech_cart_masters_id_seq',473, true);
SELECT setval('public.ingredient_masters_id_seq',210, true);
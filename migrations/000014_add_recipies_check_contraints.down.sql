ALTER TABLE recipies DROP CONSTRAINT IF EXISTS recipies_portions_check;

ALTER TABLE recipe_ingredients DROP CONSTRAINT IF EXISTS recipe_ingredients_amount_check;
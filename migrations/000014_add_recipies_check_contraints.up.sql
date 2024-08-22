ALTER TABLE recipies ADD CONSTRAINT recipies_portions_check CHECK (portions >= 1);

ALTER TABLE recipe_ingredients ADD CONSTRAINT recipe_ingredients_amount_check CHECK (amount >= 1);
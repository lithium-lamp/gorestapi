create table IF NOT EXISTS recipies (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  description text NOT NULL,
  cooking_steps text [] NOT NULL DEFAULT '{}',
  cook_time_minutes int NOT NULL,
  portions int NOT NULL,
  tags text [] NOT NULL DEFAULT '{}',
  version integer NOT NULL DEFAULT 1
);

create table IF NOT EXISTS ingredients (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  tags text [] NOT NULL DEFAULT '{}',
  version integer NOT NULL DEFAULT 1
);

create table recipe_ingredients (
  recipe_id bigint REFERENCES recipies(id) NOT NULL,
  ingredient_id bigint REFERENCES ingredients(id) NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  amount int NOT NULL,
  measurement bigint REFERENCES measurements(id) NOT NULL,
  version integer NOT NULL DEFAULT 1,
  PRIMARY KEY (recipe_id, ingredient_id)
);
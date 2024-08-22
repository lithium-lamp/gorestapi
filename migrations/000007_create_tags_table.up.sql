CREATE TABLE tags (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    itemtype bigint REFERENCES itemtypes(id) NOT NULL,
    name text NOT NULL UNIQUE,
    version integer NOT NULL DEFAULT 1
);

INSERT INTO tags (itemtype, name)
VALUES
    (1, 'milk'),
    (1, 'fruit'),
    (1, 'drinkable'),
    (2, 'cleaning'),
    (2, 'bathroom');

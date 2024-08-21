CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES
    ('availableitems:read'),
    ('availableitems:write'),
    ('itemtypes:read'),
    ('itemtypes:write'),
    ('measurements:read'),
    ('measurements:write'),
    ('tags:read'),
    ('tags:write'),
    ('knownitems:read'),
    ('knownitems:write'),
    ('recipeingredients:read'),
    ('recipeingredients:write'),
    ('ingredients:read'),
    ('ingredients:write'),
    ('recipies:read'),
    ('recipies:write'),
    ('metrics:view');
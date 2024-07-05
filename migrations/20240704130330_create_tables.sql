-- +goose Up
-- +goose StatementBegin
CREATE TABLE burgers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    measure VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE burger_ingredients (
    burger_id INT NOT NULL,
    ingredient_id INT NOT NULL,
    measure VARCHAR(255),
    PRIMARY KEY (burger_id, ingredient_id),
    FOREIGN KEY (burger_id) REFERENCES burgers(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE burger_ingredients;

DROP TABLE ingredients;

DROP TABLE burgers;
-- +goose StatementEnd

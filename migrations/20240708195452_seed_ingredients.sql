-- +goose Up
-- +goose StatementBegin
INSERT INTO ingredients (name, description)
VALUES
    ('Beef Patty', 'Juicy beef patty'),
    ('Cheddar Cheese', 'Sharp cheddar cheese'),
    ('Lettuce', 'Crisp lettuce'),
    ('Tomato', 'Fresh tomato slice'),
    ('Pickles', 'Sliced pickles'),
    ('Onion', 'Sliced onion rings'),
    ('Burger Bun', 'Toasted burger bun'),
    ('Bacon Strips', 'Crispy bacon strips'),
    ('Avocado', 'Fresh avocado slices'),
    ('Swiss Cheese', 'Creamy Swiss cheese'),
    ('Sautéed Mushrooms', 'Seasoned mushrooms'),
    ('Grilled Pineapple', 'Juicy grilled pineapple slice'),
    ('Chipotle Sauce', 'Spicy chipotle sauce'),
    ('Pepper Jack Cheese', 'Spicy pepper jack cheese'),
    ('Guacamole', 'Creamy guacamole'),
    ('Mayonnaise', 'Creamy mayonnaise'),
    ('Blue Cheese Crumbles', 'Tangy blue cheese crumbles'),
    ('Caramelized Onions', 'Sweet caramelized onions'),
    ('Parmesan Cheese', 'Grated parmesan cheese'),
    ('Blue Cheese Dressing', 'Creamy blue cheese dressing'),
    ('Italian Sausage', 'Spicy Italian sausage patty'),
    ('Marinara Sauce', 'Rich marinara sauce'),
    ('Thai Peanut Sauce', 'Savory Thai peanut sauce'),
    ('Ham', 'Sliced ham'),
    ('Ranch Dressing', 'Creamy ranch dressing'),
    ('Gruyere Cheese', 'Melty gruyere cheese'),
    ('Korean BBQ Sauce', 'Spicy-sweet Korean BBQ sauce'),
    ('Kimchi', 'Fermented cabbage'),
    ('Feta Cheese', 'Crumbly feta cheese'),
    ('Tzatziki Sauce', 'Cool tzatziki sauce'),
    ('Fresh Basil', 'Fresh basil leaves'),
    ('Mozzarella Cheese', 'Fresh mozzarella cheese');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM ingredients;
-- +goose StatementEnd

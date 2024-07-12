-- +goose Up
-- +goose StatementBegin
INSERT INTO burgers (name, description, is_vegan, image_url, updated_at)
VALUES
    ('Classic Cheeseburger', 'A classic cheeseburger with cheddar cheese', true, 'assets/classic_cheeseburger.jpg', NOW()),
    ('Bacon Avocado Burger', 'A gourmet burger with crispy bacon and fresh avocado', true, 'assets/avocado_burger.jpg', NOW()),
    ('Mushroom Swiss Burger', 'A savory burger with sautéed mushrooms and Swiss cheese', true, 'assets/Boca_Burger_2.jpg', NOW()),
    ('Spicy Jalapeño Burger', 'A fiery burger with spicy jalapeños and pepper jack cheese', false, 'assets/jalapeno_burger.jpg', NOW()),
    ('BBQ Bacon Burger', 'A hearty burger with BBQ sauce and crispy bacon', false, 'assets/Burger_3.jpg', NOW()),
    ('Mozzarella Mushroom Burger', 'A gourmet burger with mozzarella cheese and sautéed mushrooms', true, 'assets/classic_cheeseburger.jpg', NOW()),
    ('Teriyaki Pineapple Burger', 'A tangy burger with teriyaki sauce and grilled pineapple', false, 'assets/pineapple_burger.jpg', NOW()),
    ('Guacamole Burger', 'A fresh burger with guacamole and salsa', true, 'assets/classic_cheeseburger.jpg', NOW()),
    ('Crispy Chicken Burger', 'A crunchy chicken burger with lettuce and mayo', false, 'assets/Burger_3.jpg', NOW()),
    ('Vegetarian Portobello Burger', 'A flavorful vegetarian burger with grilled portobello mushroom', true, 'assets/vegan_lettuce_burger.jpg', NOW()),
    ('Garlic Parmesan Burger', 'A savory burger with garlic and parmesan cheese', false, 'assets/classic_cheeseburger.jpg', NOW()),
    ('Buffalo Chicken Burger', 'A spicy buffalo chicken burger with blue cheese dressing', false, 'assets/classic_cheeseburger.jpg', NOW()),
    ('Italian Sausage Burger', 'An Italian-inspired burger with sausage and marinara sauce', false, 'assets/Burger_3.jpg', NOW()),
    ('Thai Peanut Burger', 'A unique burger with Thai peanut sauce and cilantro', false, 'assets/Burger_3.jpg', NOW()),
    ('Hawaiian Burger', 'A tropical burger with ham and pineapple', false, 'assets/pineapple_burger.jpg', NOW()),
    ('Ranch Burger', 'A classic burger with ranch dressing and crispy bacon', false, 'assets/double_hamburger.jpg', NOW()),
    ('French Onion Burger', 'A rich burger with caramelized onions and gruyere cheese', false, 'assets/double_hamburger.jpg', NOW()),
    ('Korean BBQ Burger', 'A spicy-sweet burger with Korean BBQ sauce and kimchi', false, 'assets/mega_hamburger.jpg', NOW()),
    ('Mediterranean Burger', 'A fresh burger with feta cheese and tzatziki sauce', true, 'assets/avocado_burger.jpg', NOW()),
    ('Caprese Burger', 'A light burger with mozzarella, tomato, and basil', true, 'assets/vegan_lettuce_burger.jpg', NOW());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM burgers;
-- +goose StatementEnd

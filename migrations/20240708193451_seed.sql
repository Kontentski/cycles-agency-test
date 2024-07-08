-- +goose Up
-- +goose StatementBegin
INSERT INTO burgers (name, description, is_vegan, image_url, updated_at)
VALUES
    ('Classic Cheeseburger', 'A classic cheeseburger with cheddar cheese', true, 'localhost:8080/api/assets/classic_cheeseburger.jpg', NOW()),
    ('Bacon Avocado Burger', 'A gourmet burger with crispy bacon and fresh avocado', true, 'localhost:8080/api/assets/avocado_burger.jpg', NOW()),
    ('Mushroom Swiss Burger', 'A savory burger with sautéed mushrooms and Swiss cheese', true, 'localhost:8080/api/assets/Boca_Burger_2.jpg', NOW()),
    ('Spicy Jalapeño Burger', 'A fiery burger with spicy jalapeños and pepper jack cheese', false, 'localhost:8080/api/assets/jalapeno_burger.jpg', NOW()),
    ('BBQ Bacon Burger', 'A hearty burger with BBQ sauce and crispy bacon', false, 'localhost:8080/api/assets/Burger_3.jpg', NOW()),
    ('Mozzarella Mushroom Burger', 'A gourmet burger with mozzarella cheese and sautéed mushrooms', true, 'localhost:8080/api/assets/classic_cheeseburger.jpg', NOW()),
    ('Teriyaki Pineapple Burger', 'A tangy burger with teriyaki sauce and grilled pineapple', false, 'localhost:8080/api/assets/pineapple_burger.jpg', NOW()),
    ('Guacamole Burger', 'A fresh burger with guacamole and salsa', true, 'localhost:8080/api/assets/classic_cheeseburger.jpg', NOW()),
    ('Crispy Chicken Burger', 'A crunchy chicken burger with lettuce and mayo', false, 'localhost:8080/api/assets/Burger_3.jpg', NOW()),
    ('Vegetarian Portobello Burger', 'A flavorful vegetarian burger with grilled portobello mushroom', true, 'localhost:8080/api/assets/vegan_lettuce_burger.jpg', NOW()),
    ('Garlic Parmesan Burger', 'A savory burger with garlic and parmesan cheese', false, 'localhost:8080/api/assets/classic_cheeseburger.jpg', NOW()),
    ('Buffalo Chicken Burger', 'A spicy buffalo chicken burger with blue cheese dressing', false, 'localhost:8080/api/assets/classic_cheeseburger.jpg', NOW()),
    ('Italian Sausage Burger', 'An Italian-inspired burger with sausage and marinara sauce', false, 'localhost:8080/api/assets/Burger_3.jpg', NOW()),
    ('Thai Peanut Burger', 'A unique burger with Thai peanut sauce and cilantro', false, 'localhost:8080/api/assets/Burger_3.jpg', NOW()),
    ('Hawaiian Burger', 'A tropical burger with ham and pineapple', false, 'localhost:8080/api/assets/pineapple_burger.jpg', NOW()),
    ('Ranch Burger', 'A classic burger with ranch dressing and crispy bacon', false, 'localhost:8080/api/assets/double_hamburger.jpg', NOW()),
    ('French Onion Burger', 'A rich burger with caramelized onions and gruyere cheese', false, 'localhost:8080/api/assets/double_hamburger.jpg', NOW()),
    ('Korean BBQ Burger', 'A spicy-sweet burger with Korean BBQ sauce and kimchi', false, 'localhost:8080/api/assets/mega_hamburger.jpg', NOW()),
    ('Mediterranean Burger', 'A fresh burger with feta cheese and tzatziki sauce', true, 'localhost:8080/api/assets/avocado_burger.jpg', NOW()),
    ('Caprese Burger', 'A light burger with mozzarella, tomato, and basil', true, 'localhost:8080/api/assets/vegan_lettuce_burger.jpg', NOW());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM burgers;
-- +goose StatementEnd

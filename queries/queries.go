package queries

const GetBurgers string = `
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
`

const GetBurgerById string = `
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE b.id=$1
`

const GetBurgerByName string = `
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE LOWER(b.name) LIKE '%' || LOWER($1) || '%'
`

const GetBurgerByLetter string = `
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE LOWER(b.name) LIKE $1
`

const GetBurgerByRandom string = `
WITH random_burger AS (
	SELECT id
	FROM burgers
	ORDER BY RANDOM()
	LIMIT 1
)
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE b.id = (SELECT id FROM random_burger)
`

const GetBurgersByRandom string = `
WITH random_burgers AS (
	SELECT id
	FROM burgers
	ORDER BY RANDOM()
	LIMIT 10
)
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE b.id IN (SELECT id FROM random_burgers)
`

const GetLatestBurgers string = `
WITH latest_burgers AS (
	SELECT id
	FROM burgers
	ORDER BY updated_at DESC
	LIMIT 10
)
SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
	i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
FROM burgers b
LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
LEFT JOIN ingredients i ON bi.ingredient_id = i.id
WHERE b.id IN (SELECT id FROM latest_burgers)
ORDER BY b.updated_at DESC
`

const GetIngredientByName string = `
SELECT id, name, description
FROM ingredients
WHERE LOWER(name) LIKE LOWER($1)
LIMIT 1
`

const GetIngredientByID string = `
SELECT id, name, description
FROM ingredients
WHERE id = $1
LIMIT 1
`

const GetBurgersByIngredientName = `
SELECT b.name, b.image_url, b.id
FROM burgers b
JOIN burger_ingredients bi ON b.id = bi.burger_id
JOIN ingredients i ON bi.ingredient_id = i.id
WHERE LOWER(i.name) LIKE LOWER($1)
`

const GetBurgersByVeganStatus = `
SELECT b.id, b.name, b.image_url
FROM burgers b
WHERE b.is_vegan = $1
`

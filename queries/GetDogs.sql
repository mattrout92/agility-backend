SELECT
    id,
    name,
    height,
    grade,
    handler,
    user_id
FROM
    dogs
WHERE
    user_id = ?
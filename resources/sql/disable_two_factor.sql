UPDATE
    users
SET
    secret = ''
WHERE
    name = $1;


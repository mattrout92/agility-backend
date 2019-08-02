SELECT
    id,
    date,
    location,
    judging_from,
    entries_close
FROM
    shows
WHERE
    UNIX_TIMESTAMP() < entries_close
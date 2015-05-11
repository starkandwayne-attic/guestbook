SELECT
concat('INSERT INTO posts (orig_id, url, title) VALUES(',
id, ',''',
concat(
'https://blog.starkandwayne.com/'::text,
to_char(published_at, 'YYYY'), '/',
to_char(published_at, 'MM'), '/',
to_char(published_at, 'DD'), '/',
slug
), ''',''',
REPLACE(title,'''',''''''), ''');')
FROM posts WHERE status = 'published'

select id from posts where id = 51
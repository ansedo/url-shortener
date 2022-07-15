UPDATE urls SET is_deleted = TRUE WHERE short_url_id = $1 AND uid = $2

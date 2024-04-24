// пакет содержит начальную миграцию для создания базы.
package postgres

// функция Up возвращает sql строку для создания схемы.
func Up() string {
	return `
		CREATE TABLE shot_url (
		id         SERIAL NOT NULL,
		hash_value BIGSERIAL NOT NULL,
		origin_url VARCHAR(2048) UNIQUE,
		PRIMARY KEY (id)
	);
	`
}

// функция Down возвращает sql строку для отката миграции.
func Down() string {
	return `
	DROP TABLE shot_url;
	`
}

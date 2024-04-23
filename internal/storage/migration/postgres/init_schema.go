// пакет содержит начальную миграцию для создания базы.
package postgres

// функция Up возвращает sql строку для создания схемы
func Up() string {
	return `
		CREATE TABLE shot_url (
		id         SERIAL  NOT NULL,
		hash_value BIGINT NOT NULL,
		origin_url TEXT    UNIQUE,
		PRIMARY KEY (id)
	);
	`
}

// функция Down возвращает sql строку для отката миграции
func Down() string {
	return `
	DROP TABLE shot_url;
	`
}

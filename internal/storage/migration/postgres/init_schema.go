// пакет содержит начальную миграцию для создания базы.
package postgres

// функция Up возвращает sql строку для создания схемы.
func Migration() map[int]string {
	return map[int]string{
		1: `
			CREATE TABLE shot_url (
			id         SERIAL NOT NULL,
			hash_value BIGSERIAL NOT NULL,
			origin_url VARCHAR(2048) UNIQUE,
			PRIMARY KEY (id)
		);
		`,
		2: `
			ALTER TABLE shot_url ADD user_id SERIAL;
		`,
		3: `
			ALTER TABLE shot_url ADD is_deleted BOOLEAN;		
		`,
	}
}

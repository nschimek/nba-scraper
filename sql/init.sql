GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON nba.* TO 'go';

CREATE TABLE `nba`.`players` (
  `id` VARCHAR(9) NOT NULL,
  `name` VARCHAR(30) NOT NULL,
  `shoots` ENUM('L', 'R') NOT NULL,
  `birth_place` VARCHAR(50) NOT NULL,
  `birth_country_code` VARCHAR(2) NOT NULL,
  `birth_date` DATE NOT NULL,
  `height` TINYINT UNSIGNED NOT NULL,
  `weight` SMALLINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);

CREATE TABLE `nba`.`teams` (
  `id` VARCHAR(3) NOT NULL,
  `name` VARCHAR(30) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);

CREATE TABLE `nba`.`team_players` (
  `team_id` VARCHAR(3) NOT NULL,
  `player_id` VARCHAR(9) NOT NULL,
  `season` TINYINT UNSIGNED NOT NULL,
  `position` VARCHAR(2) NOT NULL,
  `number` TINYINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`team_id`, `player_id`, `season`));

CREATE TABLE `nba`.`team_player_salaries` (
  `team_id` VARCHAR(3) NOT NULL,
  `player_id` VARCHAR(9) NOT NULL,
  `season` TINYINT UNSIGNED NOT NULL,
  `salary` BIGINT UNSIGNED NOT NULL,
  `rank` TINYINT UNSIGNED NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`team_id`, `player_id`, `season`));
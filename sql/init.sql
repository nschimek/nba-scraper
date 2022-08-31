GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON nba.* TO 'go';

CREATE TABLE `game_four_factors` (
  `game_id` varchar(12) NOT NULL,
  `team_id` varchar(3) NOT NULL,
  `pace` float NOT NULL,
  `effective_fg_pct` float NOT NULL,
  `turnover_pct` float NOT NULL,
  `offensive_rb_pct` float NOT NULL,
  `ft_per_fga` float NOT NULL,
  `offensive_rating` float NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`game_id`,`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `game_line_scores` (
  `game_id` varchar(12) NOT NULL,
  `team_id` varchar(3) NOT NULL,
  `quarter` tinyint unsigned NOT NULL,
  `score` smallint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`game_id`,`team_id`,`quarter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `game_player_advanced_stats` (
  `game_id` varchar(12) NOT NULL,
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `true_shooting_pct` float NOT NULL,
  `effective_fg_pct` float NOT NULL,
  `three_pt_attempt_rate` float NOT NULL,
  `free_throw_attempt_rate` float NOT NULL,
  `offensive_rb_pct` float NOT NULL,
  `defensive_rb_pct` float NOT NULL,
  `total_rb_pct` float NOT NULL,
  `assist_pct` float NOT NULL,
  `steal_pct` float NOT NULL,
  `block_pct` float NOT NULL,
  `turnover_pct` float NOT NULL,
  `usage_pct` float NOT NULL,
  `box_plus_minus` float NOT NULL,
  `offensive_rating` smallint NOT NULL,
  `defensive_rating` smallint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`game_id`,`team_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `game_player_basic_stats` (
  `game_id` varchar(12) NOT NULL,
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `quarter` tinyint unsigned NOT NULL,
  `time_played` smallint NOT NULL,
  `field_goals` tinyint unsigned NOT NULL,
  `field_goals_attempted` tinyint unsigned NOT NULL,
  `field_goal_pct` float NOT NULL,
  `three_pointers` tinyint unsigned NOT NULL,
  `three_pointers_attempted` tinyint unsigned NOT NULL,
  `three_pointers_pct` float NOT NULL,
  `free_throws` tinyint unsigned NOT NULL,
  `free_throws_attempted` tinyint unsigned NOT NULL,
  `free_throws_pct` float NOT NULL,
  `offensive_rb` tinyint unsigned NOT NULL,
  `defensive_rb` tinyint unsigned NOT NULL,
  `total_rb` tinyint unsigned NOT NULL,
  `assists` tinyint unsigned NOT NULL,
  `steals` tinyint unsigned NOT NULL,
  `blocks` tinyint unsigned NOT NULL,
  `turnovers` tinyint unsigned NOT NULL,
  `personal_fouls` tinyint unsigned NOT NULL,
  `points` tinyint unsigned NOT NULL,
  `plus_minus` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`game_id`,`team_id`,`player_id`,`quarter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `game_players` (
  `game_id` varchar(12) NOT NULL,
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `status` enum('S','R','D','I') NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`game_id`,`team_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `games` (
  `id` varchar(12) NOT NULL,
  `location` varchar(255) NOT NULL,
  `type` enum('R','P') NOT NULL,
  `season` smallint unsigned NOT NULL,
  `quarters` tinyint unsigned NOT NULL,
  `start_time` datetime NOT NULL,
  `home_team_id` varchar(3) NOT NULL,
  `home_score` smallint unsigned NOT NULL,
  `home_result` enum('W','L') NOT NULL,
  `home_wins` tinyint unsigned NOT NULL,
  `home_losses` tinyint unsigned NOT NULL,
  `away_team_id` varchar(3) NOT NULL,
  `away_score` smallint unsigned NOT NULL,
  `away_result` enum('W','L') NOT NULL,
  `away_wins` tinyint unsigned NOT NULL,
  `away_losses` tinyint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`) /*!80000 INVISIBLE */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `player_injuries` (
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `season` smallint unsigned NOT NULL,
  `source_update_date` date NOT NULL,
  `description` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`team_id`,`player_id`,`season`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `players` (
  `id` varchar(9) NOT NULL,
  `name` varchar(30) NOT NULL,
  `shoots` enum('L','R') NOT NULL,
  `birth_place` varchar(50) NOT NULL,
  `birth_country_code` varchar(2) NOT NULL,
  `birth_date` date NOT NULL,
  `height` smallint unsigned NOT NULL,
  `weight` smallint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `team_player_salaries` (
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `season` smallint unsigned NOT NULL,
  `salary` bigint unsigned NOT NULL,
  `rank` tinyint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`team_id`,`player_id`,`season`),
  KEY `team_player_salaries.id2team.id_idx` (`team_id`),
  CONSTRAINT `team_player_salaries.id2team.id` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `team_players` (
  `team_id` varchar(3) NOT NULL,
  `player_id` varchar(9) NOT NULL,
  `season` smallint unsigned NOT NULL,
  `position` varchar(2) NOT NULL,
  `number` tinyint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`team_id`,`player_id`,`season`),
  KEY `team_players.id2team.id_idx` (`team_id`),
  CONSTRAINT `team_players.id2team.id` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `team_standings` (
  `team_id` varchar(3) NOT NULL,
  `season` smallint unsigned NOT NULL,
  `rank` tinyint unsigned DEFAULT NULL,
  `overall_wins` tinyint unsigned DEFAULT NULL,
  `overall_losses` tinyint unsigned DEFAULT NULL,
  `home_wins` tinyint unsigned DEFAULT NULL,
  `home_losses` tinyint unsigned DEFAULT NULL,
  `road_wins` tinyint unsigned DEFAULT NULL,
  `road_losses` tinyint DEFAULT NULL,
  `east_wins` tinyint unsigned DEFAULT NULL,
  `east_losses` tinyint unsigned DEFAULT NULL,
  `west_wins` tinyint unsigned DEFAULT NULL,
  `west_losses` tinyint unsigned DEFAULT NULL,
  `atlantic_wins` tinyint unsigned DEFAULT NULL,
  `atlantic_losses` tinyint unsigned DEFAULT NULL,
  `central_wins` tinyint unsigned DEFAULT NULL,
  `central_losses` tinyint unsigned DEFAULT NULL,
  `southeast_wins` tinyint unsigned DEFAULT NULL,
  `southeast_losses` tinyint unsigned DEFAULT NULL,
  `northwest_wins` tinyint unsigned DEFAULT NULL,
  `northwest_losses` tinyint unsigned DEFAULT NULL,
  `pacific_wins` tinyint unsigned DEFAULT NULL,
  `pacific_losses` tinyint unsigned DEFAULT NULL,
  `southwest_wins` tinyint unsigned DEFAULT NULL,
  `southwest_losses` tinyint unsigned DEFAULT NULL,
  `pre_all_star_wins` tinyint unsigned DEFAULT NULL,
  `pre_all_star_losses` tinyint unsigned DEFAULT NULL,
  `post_all_star_wins` tinyint unsigned DEFAULT NULL,
  `post_all_star_losses` tinyint unsigned DEFAULT NULL,
  `margin_less_3_wins` tinyint unsigned DEFAULT NULL,
  `margin_less_3_losses` tinyint unsigned DEFAULT NULL,
  `margin_greater_10_wins` tinyint unsigned DEFAULT NULL,
  `margin_greater_10_losses` tinyint unsigned DEFAULT NULL,
  `oct_wins` tinyint unsigned DEFAULT NULL,
  `oct_losses` tinyint unsigned DEFAULT NULL,
  `nov_wins` tinyint unsigned DEFAULT NULL,
  `nov_losses` tinyint unsigned DEFAULT NULL,
  `dec_wins` tinyint unsigned DEFAULT NULL,
  `dec_losses` tinyint unsigned DEFAULT NULL,
  `jan_wins` tinyint unsigned DEFAULT NULL,
  `jan_losses` tinyint unsigned DEFAULT NULL,
  `feb_wins` tinyint unsigned DEFAULT NULL,
  `feb_losses` tinyint unsigned DEFAULT NULL,
  `mar_wins` tinyint unsigned DEFAULT NULL,
  `mar_losses` tinyint unsigned DEFAULT NULL,
  `apr_wins` tinyint unsigned DEFAULT NULL,
  `apr_losses` tinyint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`team_id`,`season`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `teams` (
  `id` varchar(3) NOT NULL,
  `name` varchar(30) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `nba`.`team_players` 
ADD CONSTRAINT `team_players.team_id2teams.id`
  FOREIGN KEY (`team_id`)
  REFERENCES `nba`.`teams` (`id`)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;

ALTER TABLE `nba`.`team_players` 
ADD CONSTRAINT `team_players.player_id2players.id`
  FOREIGN KEY (`player_id`)
  REFERENCES `nba`.`players` (`id`)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;

ALTER TABLE `nba`.`team_player_salaries` 
ADD CONSTRAINT `team_player_salaries.player_id2teams.id`
  FOREIGN KEY (`team_id`)
  REFERENCES `nba`.`teams` (`id`)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;

ALTER TABLE `nba`.`team_standings` 
ADD CONSTRAINT `team_standings.team_id2team.id`
  FOREIGN KEY (`team_id`)
  REFERENCES `nba`.`teams` (`id`)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;

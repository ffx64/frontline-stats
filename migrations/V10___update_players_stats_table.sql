ALTER TABLE players_stats
    DROP COLUMN IF EXISTS kill_streak_x3,
    DROP COLUMN IF EXISTS kill_streak_x5,
    DROP COLUMN IF EXISTS kill_streak_x10,
    DROP COLUMN IF EXISTS kill_streak_x20,
    DROP COLUMN IF EXISTS kill_streak_x30,
    DROP COLUMN IF EXISTS max_kill_streak,
    DROP COLUMN IF EXISTS bullets_shot,
    DROP COLUMN IF EXISTS friendly_kills;

ALTER TABLE players_stats
    ADD COLUMN friendly_fire_kills        INT DEFAULT 0 NOT NULL,
    ADD COLUMN friendly_fire_deaths       INT DEFAULT 0 NOT NULL,
    ADD COLUMN headshots_made             INT DEFAULT 0 NOT NULL,
    ADD COLUMN headshots_taken            INT DEFAULT 0 NOT NULL,
    ADD COLUMN vehicle_kills              INT DEFAULT 0 NOT NULL,
    ADD COLUMN vehicle_deaths             INT DEFAULT 0 NOT NULL,
    ADD COLUMN longest_kill_distance      FLOAT8 DEFAULT 0.0 NOT NULL,
    ADD COLUMN average_kill_distance      FLOAT8 DEFAULT 0.0 NOT NULL,
    ADD COLUMN average_death_distance     FLOAT8 DEFAULT 0.0 NOT NULL,
    ADD COLUMN weapons_most_used          VARCHAR(100) DEFAULT '' NOT NULL,
    ADD COLUMN vehicle_most_used          VARCHAR(100) DEFAULT '' NOT NULL,
    ADD COLUMN hit_zones_most_killed      VARCHAR(50) DEFAULT '' NOT NULL,
    ADD COLUMN hit_zones_most_died        VARCHAR(50) DEFAULT '' NOT NULL,
    ADD COLUMN ratio_kdr                  NUMERIC(6,3) DEFAULT 0.0 NOT NULL,
    ADD COLUMN ratio_headshot             NUMERIC(6,3) DEFAULT 0.0 NOT NULL,
    ADD COLUMN ratio_friendly_fire        NUMERIC(6,3) DEFAULT 0.0 NOT NULL,
    ADD COLUMN ratio_vehicle              NUMERIC(6,3) DEFAULT 0.0 NOT NULL;

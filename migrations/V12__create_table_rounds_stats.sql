CREATE TABLE IF NOT EXISTS rounds_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    round_id UUID NOT NULL,
    server_id UUID NOT NULL,
    player_id UUID NOT NULL,
    team VARCHAR(20) NOT NULL,
    total_kills BIGINT NOT NULL DEFAULT 0,
    total_deaths BIGINT NOT NULL DEFAULT 0,
    total_suicides BIGINT NOT NULL DEFAULT 0,
    total_vehicle_kills BIGINT NOT NULL DEFAULT 0,
    total_vehicle_deaths BIGINT NOT NULL DEFAULT 0,
    average_kill_distance DOUBLE PRECISION NOT NULL DEFAULT 0,
    total_team_kills BIGINT NOT NULL DEFAULT 0,
    total_headshots BIGINT NOT NULL DEFAULT 0,
    most_used_weapon VARCHAR(255) DEFAULT '',
    most_hit_zone VARCHAR(255) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE rounds_stats ADD CONSTRAINT unique_round_player UNIQUE (round_id, player_id);

CREATE TABLE IF NOT EXISTS servers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS server_settings (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
	setting_key TEXT NOT NULL,
	UNIQUE (server_id, setting_key)
);

CREATE TABLE IF NOT EXISTS players (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guid VARCHAR(64) UNIQUE NOT NULL,
	username TEXT NOT NULL,
	admin INTEGER DEFAULT 0,
	premium INTEGER DEFAULT 0,
	premium_start_at TIMESTAMPTZ DEFAULT NULL,
	premium_expire_at TIMESTAMPTZ DEFAULT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	is_banned BOOLEAN DEFAULT FALSE,
	last_server_id UUID REFERENCES servers(id),
	last_login TIMESTAMPTZ,
	platform TEXT,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS players_stats (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	player_id UUID NOT NULL REFERENCES players(id),
	level INTEGER DEFAULT 0,
	xp INTEGER DEFAULT 0,
	kills INTEGER DEFAULT 0,
	deaths INTEGER DEFAULT 0,
	friendly_kills INTEGER DEFAULT 0,
	bullets_shot INTEGER DEFAULT 0,
	grenades_thrown INTEGER DEFAULT 0,
	max_kill_streak INTEGER DEFAULT 0,
	max_kill_distance INTEGER DEFAULT 0,
	insertion_bonus INTEGER DEFAULT 0,
	kill_streak_x3 INTEGER DEFAULT 0,
	kill_streak_x5 INTEGER DEFAULT 0,
	kill_streak_x10 INTEGER DEFAULT 0,
	kill_streak_x20 INTEGER DEFAULT 0,
	kill_streak_x30 INTEGER DEFAULT 0,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS players_bans (
    ban_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_guid VARCHAR(64) NOT NULL,
    issued_by_guid VARCHAR(64) NOT NULL,
    ban_reason TEXT,
	is_active BOOLEAN DEFAULT TRUE,
    ban_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clans (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL,
	tag TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rounds (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
	current_mode TEXT,
	mission_header TEXT,
	status TEXT,
	winner_faction TEXT DEFAULT 'null',
	ended_at TIMESTAMPTZ,
	start_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS kill_logs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
	round_id UUID NOT NULL REFERENCES rounds(id) ON DELETE CASCADE,
	killer_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
	victim_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
	killer_weapon TEXT NOT NULL,
	victim_weapon TEXT NOT NULL,
	distance FLOAT NOT NULL,
	is_friendly BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clan_members (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	clan_id UUID NOT NULL REFERENCES clans(id) ON DELETE CASCADE,
	player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
	is_owner BOOLEAN DEFAULT FALSE,
	joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (clan_id, player_id)
);

CREATE INDEX idx_players_username ON players(username);
CREATE INDEX idx_kill_logs_server_id ON kill_logs(server_id);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_player_timestamp
BEFORE UPDATE ON players
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

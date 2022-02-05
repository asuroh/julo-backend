
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table wallet (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
	balance integer NOT NULL,
	owned_by uuid NOT NULL,
	status TEXT CHECK (char_length(status) <= 8),
	enabled_at TIMESTAMP WITH TIME ZONE,
	disabled_at TIMESTAMP WITH TIME ZONE
);

create table balance (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
	amount integer NOT NULL,
	status TEXT NOT NULL CHECK (char_length(status) <= 20),
	reference_id uuid NOT NULL,
	deposited_by uuid,
	deposited_at TIMESTAMP WITH TIME ZONE,
	withdrawn_by uuid,
	withdrawn_at TIMESTAMP WITH TIME ZONE
);
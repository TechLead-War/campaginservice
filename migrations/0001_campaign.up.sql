CREATE TABLE campaigns (
    campaign_id TEXT PRIMARY KEY,
    campaign_name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    call_to_action TEXT NOT NULL,
    campaign_status TEXT NOT NULL CHECK (campaign_status IN ('ACTIVE', 'INACTIVE')),
    cdate TIMESTAMPTZ NOT NULL DEFAULT now(),
    udate TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE targeting_rules (
    campaign_id TEXT REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
    dimension TEXT NOT NULL CHECK (dimension IN ('country', 'os', 'app_id')),
    type TEXT NOT NULL CHECK (type IN ('include', 'exclude')),
    value TEXT NOT NULL,
    cdate TIMESTAMPTZ NOT NULL DEFAULT now(),
    udate TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (campaign_id, dimension, type, value)
);

-- This is indexed only for (dimension, type, value), not for campaign_id since the lookup for delivery query will be done on these 3 fields only.
CREATE INDEX idx_targeting_lookup ON targeting_rules (dimension, type, value);
CREATE TABLE IF NOT EXISTS default.logs (
    `timestamp` DateTime64(3) CODEC(Delta, ZSTD),
    `service_name` String CODEC(ZSTD),
    `level` LowCardinality(String) CODEC(ZSTD),
    `message` String CODEC(ZSTD),
    `metadata` Map(String, String) CODEC(ZSTD)
) ENGINE = MergeTree() PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, service_name);
CREATE KEYSPACE IF NOT EXISTS koperasi_analytics
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};


Connect to your Cassandra container and create the keyspace manually:

```bash
# Connect to Cassandra container
docker exec -it cassandra_db cqlsh -u cassandra -p cassandra

# Create keyspace
CREATE KEYSPACE IF NOT EXISTS koperasi_analytics
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

# Exit cqlsh
exit
```

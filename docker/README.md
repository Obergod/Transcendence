## In the case of emergency the procedure to recover data from services should look like this:

1. Run *docker exec -it backup bash*.
2. Find directory called *backup*.
3. To extract needed recovery files run *tar -xvzf name_of_the_file.tar.gz*
4. Copy files to the needed folders. In the case of database: /var/lib/postgresql/data. In the case of prometheus: /prometheus. In the case of Elasticsearch: /usr/share/elasticsearch/data.


## In the case of restoring snapshots from Elasticsearch use:
1. POST _snapshot/s3-repo/_restore
`    {
    "indices": "*",
    "include_global_state": false
    }`
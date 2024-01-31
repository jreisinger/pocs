Work-queue shows how to use Kubernetes jobs to process a work queue.

```
producer ---> [ work queue ] ---> consumer
                                  consumer
                                  consumer
                                  consumer
                                  consumer
```

```
# Launch a centralized queue service.
k apply -f queue.yaml

# In different terminal.
k port-forward rs/queue 8080:8080

# Create a work queue called 'keygen'.
curl -X PUT localhost:8080/memq/server/queues/keygen

# Create 100 work items and load up the queue.
for i in work-item-{0..99}; do
    curl -X POST localhost:8080/memq/server/queues/keygen/enqueue -d $i
done

# Start the consumers.
k apply -f consumers.yaml

# Check what's going on.
curl localhost:8080/memq/server/stats

# Clean up.
k delete rs,svc,job -l app=work-queue
```

For more see Kubernetes: Up and Running, ch. 12 Jobs
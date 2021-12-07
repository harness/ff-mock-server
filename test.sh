ID=$(docker run -d -p 9090:3000 ff-mock-server:latest)
docker exec ${ID} /bin/bash /app/wait-for-it.sh
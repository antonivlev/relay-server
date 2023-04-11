Start db:

```bash
# kill all container
docker kill $(docker ps -a -q)

# run
docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=my-secret-pw -dp 3306:3306 mysql:latest
# or
docker start mysql-container

# get into db
docker exec -it mysql-container mysql -uroot -pmy-secret-pw

# create db
USE relay;
CREATE TABLE users (
  id INT AUTO_INCREMENT,
  created_at DATETIME,
  email TEXT NOT NULL,
  number_of_tokens FLOAT,
  password TEXT NOT NULL,
  PRIMARY KEY (id)
);

INSERT INTO users 
  (email, password, created_at) 
VALUES 
  ("anton@mail.com", "12345", NOW());
```
# dynamic-user-segmentation-service

## Common app schema

Database: PostgreSQL
Database gui: adminer
Web-framework: Gin


Common shema

![image](https://github.com/mark47B/dynamic-user-segmentation-service/assets/43784470/79c8d544-3d6e-49ae-961f-e5ec9ca6e86f)


## Run ✈️

```bash
    git clone git@github.com:mark47B/dynamic-user-segmentation-service.git
```

```bash
    cd dynamic-user-segmentation-service
```

```bash
    docker-compose --env-file "config/local.env" up
```

## Tests 
### Database 

adminer is available at http://0.0.0.0:8091

> System: PostgreSQL </br>
> Server: postgres </br>
> Username: user </br>
> Password: 1234 </br>
> Database: database </br>

#### Shema: 
![image](https://github.com/mark47B/dynamic-user-segmentation-service/assets/43784470/294ec6b0-0da2-4ce8-8bbf-3bdd4a811e7e)

### API
Base task
![image](https://github.com/mark47B/dynamic-user-segmentation-service/assets/43784470/c1c63138-5196-4a91-8a4e-68b26f9291ec)

All available API
```
GET    /api/v1/user/:uuid
GET    /api/v1/user/:uuid/slugs  
PUT    /api/v1/user/:uuid     
GET    /api/v1/user/
POST   /api/v1/user/ 
POST   /api/v1/slug/
DELETE /api/v1/slug/:name
```

### Curl commands
#### Base task
1. Create Slug
    ```bash
    curl http://localhost:8080/api/v1/slug/ \
    --include --header "Content-Type: application/json" \
    --request "POST" --data \
    '{"name": "AVITO_TEST_SLUG"}'
    ```

   ```bash
    curl http://localhost:8080/api/v1/slug/ \
    --include --header "Content-Type: application/json" \
    --request "POST" --data \
    '{"name": "AVITO_VOICE_MESSAGES"}'
    ```

   ```bash
    curl http://localhost:8080/api/v1/slug/ \
    --include --header "Content-Type: application/json" \
    --request "POST" --data \
    '{"name": "AVITO_PERFORMANCE_VAS"}'
    ```

2. Delete Slug
    ```bash
    curl http://localhost:8080/api/v1/slug/AVITO_TEST_SLUG \
    --include --header "Content-Type: application/json" \
    --request "DELETE"
    ```
    
3. Change User's slugs
    ```bash
    curl http://localhost:8080/api/v1/user/83f31292-dca2-4870-87aa-76d21c750739 \
    --include --header "Content-Type: application/json" \
    --request "PUT" --data \
    '{"delete_slugs": [], "add_slugs": ["AVITO_PERFORMANCE_VAS", "AVITO_TEST_SLUG", "AVITO_VOICE_MESSAGES"]}'
    ```

   ```bash
    curl http://localhost:8080/api/v1/user/83f31292-dca2-4870-87aa-76d21c750739 \
    --include --header "Content-Type: application/json" \
    --request "PUT" --data \
    '{"delete_slugs": ["AVITO_PERFORMANCE_VAS"], "add_slugs": []}'
    ```
5. Get active Slugs for User
    ```bash
    curl http://localhost:8080/api/v1/user/{UUID}/slugs \
    --include --header "Content-Type: application/json" \
    --request "GET"
    ```

#### Additioonal API features
1. Get user by UUID
    ```bash
    curl http://localhost:8080/api/v1/user/{UUID} \
    --include --header "Content-Type: application/json" \
    --request "GET"
    ```
2. Get all users
    ```bash
    curl http://localhost:8080/api/v1/user/ \
    --include --header "Content-Type: application/json" \
    --request "GET"
    ```
3. Create user
    ```bash
    curl http://localhost:8080/api/v1/user/ \
    --include --header "Content-Type: application/json" \
    --request "POST" --data \
    '{"username": "Alexandr"}'
    ```
    Response: UUID




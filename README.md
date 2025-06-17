# Setup service:
1. Clone the repository:
```bash
git clone git@github.com:TechLead-War/campaignservice.git
```
2. Run migrations.
```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/campaign_service?sslmode=disable" up
```

3. Seed some data, with the following command:
```bash

##On Local:
go run seed.go -records=2000 -workers=20

##Docker based:
docker-compose run --rm seed
```

4. Run the service:
```bash
##On Local:
air

## Docker based:
docker-compose down -v      
docker-compose build
docker-compose up
```
## DB Schema:
![Screenshot 2025-06-17 at 4 54 49 PM](https://github.com/user-attachments/assets/f2d70357-229f-4c0b-81da-b019fa4e7f5d)


### Note:<br>
This repo is tested with Go 1.23.5 and Postgres 15.4. With a 30L+ records.

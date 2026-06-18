# Microservice Monitoring

## Swagger docs genereren

Voer uit in de root van het project:

```powershell
swag init -g main.go -o docs
```

## Lokaal opstarten met Docker Compose

Voer uit in de root van het project:

```powershell
docker compose up -d
```

Dit start onder andere:

- timescaledb (PostgreSQL/Timescale)
- nats
- migrate (eenmalig, stopt daarna met exit code 0)
- monitoring-data-generator
- monitoring-service
- dashboard (Angular)

## Controleren of alles draait

```powershell
docker compose ps
```

Belangrijke endpoints:

- Angular dashboard: http://localhost:4200
- Monitoring API: http://localhost:8080
- NATS monitor: http://localhost:8222

Logs bekijken:

```powershell
docker compose logs -f dashboard
docker compose logs -f monitoring-service
docker compose logs -f monitoring-data-generator
```

## Stoppen

```powershell
docker compose down
```

Wil je ook de database-volume opruimen (let op: data gaat verloren):

```powershell
docker compose down -v
```
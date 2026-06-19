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

Wanneer je code changes heb gedaan en dan opnieuw docker wil builden:
```powershell
docker compose down 
docker compose up -d --build
```

Met generator erbij:
```powershell
docker compose --profile generator down 
docker compose --profile generator up -d --build
```

Dit start onder andere:

- timescaledb (PostgreSQL/Timescale)
- nats
- migrate (eenmalig, stopt daarna met exit code 0)
- monitoring-service
- dashboard (Angular)

## Generator aan/uit zetten

De `monitoring-data-generator` is optioneel gemaakt via een Docker Compose profile.

Generator uit (default):

```powershell
docker compose up -d
```

De generator kan gestart worden met het volgende commando:

```powershell
docker compose --profile generator up -d
```

Data Gen stoppen:
```powershell
docker compose stop monitoring-data-generator
```

Data Gen weer aanzetten:
```powershell
docker compose start monitoring-data-generator
```

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

Als de generator ook aan stond, moet je hem meegeven anders wordt hij niet gestopt/verwijderd:

```powershell
docker compose --profile generator down
```

Wil je ook de database-volume opruimen (let op: data gaat verloren):

```powershell
docker compose --profile generator down -v
```
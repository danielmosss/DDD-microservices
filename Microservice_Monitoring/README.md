# Microservice Monitoring

## Swagger docs genereren

Voer uit in de root van het project:

```powershell
swag init -g main.go -o docs
```

## Go unit tests

Run alle Go unit tests:

```powershell
go test ./...
```

Run alleen tests voor v1 REST API package:

```powershell
go test ./internal/app/restapi/v1 -v
```

## Teststrategie

Deze repo gebruikt 3 testlagen:

- Unit tests: pure functies/validatie zonder externe afhankelijkheden.
- Interface tests: HTTP contract tests op handler-niveau met `httptest`.
- Integratie tests: end-to-end tegen een draaiende API (Docker Compose).

Run alleen unit + interface tests:

```powershell
go test ./internal/... -v
```

Run integratie tests (API moet draaien op localhost:8080):

```powershell
docker compose up -d --build
go test -tags=integration ./tests/integration -v
```

## GitHub Actions (automatisch testen)

Omdat dit project in een monorepo staat, moet de workflow in de root van de monorepo staan.

De workflow draait automatisch op push en pull request (alleen bij wijzigingen in `Microservice_Monitoring/**`):

- Unit + interface tests: `go test ./internal/... -v`
- Integratie tests: start Docker Compose, wacht op API, run `go test -tags=integration ./tests/integration -v`

Workflow bestand:

- `DDD-microservices/.github/workflows/monitoring-go-tests.yml`

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
COMPOSE := docker compose -f docker-compose.yml
.PHONY: infra up down ps logs demo
infra:
	$(COMPOSE) up -d --wait postgres redis redpanda redpanda-init
up:
	$(COMPOSE) --profile app up -d --build --wait
down:
	$(COMPOSE) --profile app down
ps:
	$(COMPOSE) --profile app ps
logs:
	$(COMPOSE) --profile app logs -f pub sub act ret
demo:
	$(COMPOSE) exec -T pub demo

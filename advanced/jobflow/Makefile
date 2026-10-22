PYTHON ?= python3
LESSON ?= 1
PROJECT ?= starter
.PHONY: check check-all check-mutants smoke docs validate up down demo
check:
	$(PYTHON) tools/coursecheck.py --lesson $(LESSON) --project "$(PROJECT)"
check-all:
	$(PYTHON) tools/check_all.py
check-mutants:
	$(PYTHON) tools/check_mutants.py
smoke:
	$(PYTHON) tools/smoke_compose.py
docs:
	$(PYTHON) tools/render_course.py
validate:
	$(PYTHON) tools/validate_course.py
up:
	docker compose -f deploy/docker-compose.yml --profile app up -d --build --wait
down:
	docker compose -f deploy/docker-compose.yml --profile app down
demo:
	docker compose -f deploy/docker-compose.yml exec -T pub demo

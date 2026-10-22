PYTHON ?= python3
TRACK ?= core
LESSON ?= 1
PROJECT ?= starter
.PHONY: docs validate check check-all check-advanced check-mutants smoke up down
# Only the explicitly marked course navigation is regenerated.
docs:
	$(PYTHON) tools/render_course.py
validate:
	$(PYTHON) tools/validate_course.py
check:
	$(PYTHON) tools/coursecheck.py --track $(TRACK) --lesson $(LESSON) --project "$(PROJECT)"
check-all:
	$(PYTHON) tools/check_all.py
check-advanced:
	$(PYTHON) advanced/jobflow/tools/validate_course.py
	$(PYTHON) advanced/jobflow/tools/check_all.py
up:
	docker compose -p combat-core -f deploy/docker-compose.yml up -d --build --wait
down:
	docker compose -p combat-core -f deploy/docker-compose.yml down

check-mutants:
	$(PYTHON) tools/check_mutants.py
smoke:
	$(PYTHON) tools/smoke_compose.py

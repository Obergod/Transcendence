# Noms
NAME = transcendance
MODULE_NAME = $(NAME)

# Dossiers
BACKEND_DIR = backend
GAME_DIR = game
FRONTEND_DIR = frontend
WASM_DIR = $(FRONTEND_DIR)/public
WASM_OUT = $(WASM_DIR)/main.wasm

# Affichage
GREEN = \033[32m
YELLOW = \033[33m
CYAN = \033[36m
RESET = \033[0m
CLEAR = \033[2K\r

.PHONY: all clean fclean re up down start stop ps core launch_core reset

SRCS=docker/docker-compose.yml

all: build up clean # TODO clean has been deleted from all launch, check if it is needed 

# --- For Dockers ---

bootstrap:
	./docker/start_elastic.sh

up: bootstrap
	@docker compose -f ${SRCS} up -d

build:
	@docker compose -f ${SRCS} build backend && docker compose -f ${SRCS} build

launch_core:
	@docker compose -f ${SRCS} up -d database backend frontend

core: build launch_core clean

start:
	@docker compose -f ${SRCS} start

stop_ng:
	@docker compose -f ${SRCS} kill nginx

stop: stop_ng
	@docker compose -f ${SRCS} stop

down:
	@docker compose -f ${SRCS} down

ps:
	@docker compose -f ${SRCS} ps

# --- Nettoyage ---

clean:
	@docker image rm backend-builder docker_frontend custom-logstash -f
# 	docker image prune -af

fclean: stop down clean
	@docker container prune -f
# 	@ docker image prune -af
# 	@ docker compose -f ${SRCS} down --rmi 'all'

death: fclean
	@rm -rf docker/.env docker/.secrets .env .secrets docker/backups docker/certs docker/src/nginx/logs
	@docker system prune -af
	@docker volume prune -f

re: fclean all

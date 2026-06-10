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
SETUP=docker/docker-compose.setup.yml

all: bootstrap build up clean # build clean <up> # TODO clean has been deleted from all launch, check if it is needed 

# --- For Dockers ---

bootstrap:
	@docker compose -f ${SETUP} up certs
	@./docker/start_elastic.sh

up:
	@docker compose -f ${SRCS} up -d

build:
	@docker compose -f ${SRCS} build #backend && docker compose -f ${SRCS} build

start:
	@docker compose -f ${SRCS} start

stop: 
	@docker compose -f ${SRCS} down --remove-orphans

ps:
	@docker compose -f ${SRCS} ps

# --- Nettoyage ---

clean:
	@docker image rm backend frontend -f
# 	docker image prune -af

fclean: stop clean # down at beginning
# 	@ docker image prune -af
# 	@ docker compose -f ${SRCS} down --rmi 'all'

death: fclean
	@rm -rf docker/backups docker/certs docker/src/elasticsearch/ssl/* docker/src/kibana/ssl/*
	@echo > docker/.env
	@docker system prune -af
	@docker volume prune -f

re: fclean all
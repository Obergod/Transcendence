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

.PHONY: all clean fclean re up down start stop ps core

SRCS=docker/docker-compose.yml

all: build up

# --- For Dockers ---
up:
	@docker compose -f ${SRCS} up -d

build:
	@docker compose -f ${SRCS} build backend && docker compose -f ${SRCS} build

down:
	@docker compose -f ${SRCS} down

start:
	@docker compose -f ${SRCS} start

stop:
	@docker compose -f ${SRCS} stop

ps:
	@docker compose -f ${SRCS} ps

# --- Nettoyage ---
clean:
	@printf "$(YELLOW)Cleaning Go cache...$(RESET)"
	@$(GOCLEAN) -cache
	@printf "$(CLEAR)$(GREEN)✓ Go cache cleaned$(RESET)\n"

fclean: clean
	@printf "$(YELLOW)Deleting $(NAME)...$(RESET)"
	@rm -f $(NAME)
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting WASM files...$(RESET)"
	@rm -f $(WASM_OUT) $(WASM_DIR)/wasm_exec.js
	@printf "$(CLEAR)$(GREEN)✓ WASM files deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting go.mod and go.sum...$(RESET)"
	@rm -f go.mod go.sum
	@printf "$(CLEAR)$(GREEN)✓ go.mod and go.sum deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting frontend/dist and frontend/node_modules...$(RESET)"
	@rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules
	@printf "$(CLEAR)$(GREEN)✓ frontend/dist and node_modules deleted$(RESET)\n"
	-docker system prune -af
	@ docker compose -f ${SRCS} down --rmi 'all'

re: fclean all

# --- To launch without reinstalling dockers ---

core: build
	@docker compose -f ${SRCS} up -d backend frontend database
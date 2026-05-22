# Noms
NAME = transcendance
MODULE_NAME = $(NAME)

# Dossiers
BACKEND_DIR = backend
GAME_DIR = game
FRONTEND_DIR = frontend
WASM_DIR = $(FRONTEND_DIR)/public
WASM_OUT = $(WASM_DIR)/main.wasm

# Le vrai Goinfre local (SSD rapide !)
GOINFRE_DIR = /goinfre/$(USER)

# Affichage
GREEN = \033[32m
YELLOW = \033[33m
CYAN = \033[36m
RESET = \033[0m
CLEAR = \033[2K\r

.PHONY: all clean fclean re up build down start stop ps core setup_goinfre

SRCS=docker/docker-compose.yml

all: setup_goinfre build up

# --- MAGIE POUR 42 : Prépare le VRAI Goinfre local automatiquement ---
setup_goinfre:
	@printf "$(CYAN)Verifying local Goinfre configuration for Podman...$(RESET)\n"
	@mkdir -p $(GOINFRE_DIR)/containers
	@mkdir -p $(HOME)/.config/containers
	@printf "[storage]\ndriver = \"overlay\"\ngraphroot = \"$(GOINFRE_DIR)/containers\"\n\n[storage.options.overlay]\nignore_chown_errors = \"true\"\n" > $(HOME)/.config/containers/storage.conf
	@podman system migrate 2>/dev/null || true

# --- For Dockers ---
up: setup_goinfre
	@docker compose -f ${SRCS} up -d

build: setup_goinfre
	@docker compose -f ${SRCS} build backend && docker compose -f ${SRCS} build frontend

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
	@go clean -cache
	@printf "$(CLEAR)$(GREEN)✓ Go cache cleaned$(RESET)\n"

fclean: clean
	@printf "$(YELLOW)Deleting $(NAME)...$(RESET)"
	@rm -f $(NAME)
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting WASM files...$(RESET)"
	@rm -f $(WASM_OUT) $(WASM_DIR)/wasm_exec.js
	@printf "$(CLEAR)$(GREEN)✓ WASM files deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting frontend/dist and frontend/node_modules...$(RESET)"
	@rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules
	@printf "$(CLEAR)$(GREEN)✓ frontend/dist and node_modules deleted$(RESET)\n"
	@printf "$(CYAN)Purging Podman cache on local Goinfre...$(RESET)\n"
	@podman system prune -af 2>/dev/null || true
	@docker compose -f ${SRCS} down --rmi 'all' 2>/dev/null || true

re: fclean all

core: setup_goinfre build
	@docker compose -f ${SRCS} up backend frontend database -d
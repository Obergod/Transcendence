# Noms
NAME = transcendance
MODULE_NAME = $(NAME)

# Dossiers
BACKEND_DIR = backend
GAME_DIR = game
FRONTEND_DIR = frontend
WASM_DIR = $(FRONTEND_DIR)/public
WASM_OUT = $(WASM_DIR)/main.wasm

# Versions
EBITEN_VERSION = latest

# Commandes Go
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOTIDY = $(GOCMD) mod tidy

# Commandes Node / npm
NPM = npm
NPM_RUN = $(NPM) run

# Affichage
GREEN = \033[32m
YELLOW = \033[33m
CYAN = \033[36m
RESET = \033[0m
CLEAR = \033[2K\r

.PHONY: all clean fclean re ensure-module install-deps tidy server wasm frontend install-frontend run dev

all: server wasm frontend

# --- Module à la racine ---
ensure-module:
	@if [ ! -f go.mod ]; then \
		printf "$(YELLOW)go.mod not found at root. Running 'go mod init $(MODULE_NAME)'...$(RESET)\n"; \
		$(GOCMD) mod init $(MODULE_NAME); \
		printf "$(GREEN)✓ go.mod created at root$(RESET)\n"; \
	fi

install-deps: ensure-module
	@printf "$(YELLOW)Installing Ebitengine $(EBITEN_VERSION)...$(RESET)\n"
	@$(GOGET) github.com/hajimehoshi/ebiten/v2@$(EBITEN_VERSION)
	@printf "$(GREEN)✓ Ebitengine installed$(RESET)\n"

tidy: install-deps
	@printf "$(YELLOW)Running go mod tidy...$(RESET)\n"
	@$(GOTIDY)
	@printf "$(GREEN)✓ Dependencies tidied$(RESET)\n"

# --- Compilation serveur (binaire natif) ---
server: tidy
	@printf "$(CYAN)Building $(NAME) server...$(RESET)"
	@cd $(BACKEND_DIR) && $(GOBUILD) -o ../$(NAME) .
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) server created$(RESET)\n"

# --- Compilation WASM (jeu) ---
wasm: tidy
	@printf "$(CYAN)Building WASM game...$(RESET)"
	@mkdir -p $(WASM_DIR)
	@cd $(GAME_DIR) && GOOS=js GOARCH=wasm $(GOBUILD) -o ../$(WASM_OUT) main_wasm.go
	@printf "$(CLEAR)$(GREEN)✓ WASM created at $(WASM_OUT)$(RESET)\n"
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)/ 2>/dev/null || \
		printf "$(YELLOW)⚠ wasm_exec.js not copied (Go not found)$(RESET)\n"

# --- Installer les dépendances React ---
install-frontend:
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		printf "$(YELLOW)Installing React dependencies (this may take a while)...$(RESET)\n"; \
		cd $(FRONTEND_DIR) && $(NPM) install; \
		printf "$(GREEN)✓ React dependencies installed$(RESET)\n"; \
	else \
		printf "$(GREEN)✓ React dependencies already installed$(RESET)\n"; \
	fi

# --- Compilation frontend React ---
frontend: install-frontend
	@printf "$(CYAN)Building React frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && $(NPM_RUN) build
	@printf "$(CLEAR)$(GREEN)✓ React frontend built$(RESET)\n"

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

re: fclean all

# --- Mode développement (backend + frontend vite) ---
dev: wasm
	@printf "$(CYAN)Starting Backend and Frontend in DEV mode...$(RESET)\n"
	@trap 'kill 0' SIGINT; \
	$(GOCMD) run $(BACKEND_DIR)/main.go & \
	cd $(FRONTEND_DIR) && $(NPM_RUN) dev

# --- Mode production (backend seul, frontend déjà construit) ---
run: server wasm frontend
	@printf "$(CYAN)Starting Production Server...$(RESET)\n"
	@./$(NAME)

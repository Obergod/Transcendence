# Noms
NAME = transcendance
MODULE_NAME = $(NAME)

# Dossiers
BACKEND_DIR = backend
GAME_DIR = game
STATIC_DIR = static
WASM_OUT = $(STATIC_DIR)/main.wasm

# Versions
EBITEN_VERSION = v2.5.7

# Commandes Go
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOTIDY = $(GOCMD) mod tidy

# Affichage
GREEN = \033[32m
YELLOW = \033[33m
CYAN = \033[36m
RESET = \033[0m
CLEAR = \033[2K\r

.PHONY: all clean fclean re ensure-module install-deps tidy server wasm run

all: server wasm

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
# Le serveur est dans backend/main.go, il importe internal/...
server: tidy
	@printf "$(CYAN)Building $(NAME) server...$(RESET)"
	@cd $(BACKEND_DIR) && $(GOBUILD) -o ../$(NAME) .
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) server created$(RESET)\n"

# --- Compilation WASM (jeu) ---
# Le jeu est dans game/main_wasm.go, il importe internal/... également
wasm: tidy
	@printf "$(CYAN)Building WASM game...$(RESET)"
	@mkdir -p $(STATIC_DIR)
	@cd $(GAME_DIR) && GOOS=js GOARCH=wasm $(GOBUILD) -o ../$(WASM_OUT) main_wasm.go
	@printf "$(CLEAR)$(GREEN)✓ WASM created at $(WASM_OUT)$(RESET)"
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(STATIC_DIR)/ 2>/dev/null || \
		printf "$(YELLOW)⚠ wasm_exec.js not copied (Go not found)$(RESET)\n"

# --- Nettoyage ---
clean:
	@printf "$(YELLOW)Cleaning Go cache...$(RESET)"
	@$(GOCLEAN) -cache
	@printf "$(CLEAR)$(GREEN)✓ Go cache cleaned$(RESET)\n"

fclean: clean
	@printf "$(YELLOW)Deleting $(NAME)...$(RESET)"
	@rm -f $(NAME)
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting $(WASM_OUT) and static/wasm_exec.js...$(RESET)"
	@rm -f $(WASM_OUT) $(STATIC_DIR)/wasm_exec.js
	@printf "$(CLEAR)$(GREEN)✓ WASM files deleted$(RESET)\n"
	@printf "$(YELLOW)Deleting go.mod and go.sum...$(RESET)"
	@rm -f go.mod go.sum
	@printf "$(CLEAR)$(GREEN)✓ go.mod and go.sum deleted$(RESET)\n"

re: fclean all

run: server
	@./$(NAME)
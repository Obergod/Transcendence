# Noms
NAME = transcendance
BINARY_NAME = $(NAME)
MODULE_NAME = $(NAME)

# Dossier contenant tous les fichiers .go
BACKEND_DIR = backend

# Versions spécifiques des dépendances (pour Go 1.18)
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

.PHONY: all clean fclean re ensure-module install-deps tidy run

all: $(NAME)

# Vérifie et crée go.mod à la racine
ensure-module:
	@if [ ! -f go.mod ]; then \
		printf "$(YELLOW)go.mod not found. Running 'go mod init $(MODULE_NAME)'...$(RESET)\n"; \
		$(GOCMD) mod init $(MODULE_NAME); \
		printf "$(GREEN)✓ go.mod created$(RESET)\n"; \
	fi

# Force l'installation de la version compatible d'Ebitengine
install-deps: ensure-module
	@printf "$(YELLOW)Installing Ebitengine $(EBITEN_VERSION) (compatible with Go 1.18)...$(RESET)\n"
	@$(GOGET) github.com/hajimehoshi/ebiten/v2@$(EBITEN_VERSION)
	@printf "$(GREEN)✓ Ebitengine $(EBITEN_VERSION) installed$(RESET)\n"

# Tidy pour nettoyer et fixer les versions
tidy: install-deps
	@printf "$(YELLOW)Running go mod tidy...$(RESET)\n"
	@$(GOTIDY)
	@printf "$(GREEN)✓ Dependencies tidied$(RESET)\n"

# Compilation : on s'assure que les dépendances sont au bon niveau
$(NAME): tidy
	@printf "$(CYAN)Building $(NAME)...$(RESET)"
	@$(GOBUILD) -o $(NAME) ./$(BACKEND_DIR)
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) created!$(RESET)\n"

# Nettoyage (cache Go seulement)
clean:
	@printf "$(YELLOW)Cleaning Go cache...$(RESET)"
	@$(GOCLEAN) -cache
	@printf "$(CLEAR)$(GREEN)✓ Go cache cleaned!$(RESET)\n"

# Nettoyage complet (binaire + go.mod + go.sum)
fclean: clean
	@printf "$(YELLOW)Deleting $(NAME)...$(RESET)"
	@rm -f $(NAME)
	@rm -f go.mod go.sum
	@printf "$(CLEAR)$(GREEN)✓ $(NAME), go.mod and go.sum deleted!$(RESET)\n"

# Reconstruction
re: fclean all

# Lancer le programme
run: all
	@./$(NAME)
# Noms
NAME = transcendance
BINARY_NAME = $(NAME)
MODULE_NAME = $(NAME)  # change to github.com/yourname/transcendance if needed

# Répertoires
SRCS_DIR = src/

# Commands Go
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Affichage
GREEN = \033[32m
YELLOW = \033[33m
CYAN = \033[36m
RESET = \033[0m
CLEAR = \033[2K\r

.PHONY: all clean fclean re ensure-module run

all: $(NAME)

# Vérifie et crée go.mod si absent
ensure-module:
	@if [ ! -f go.mod ]; then \
		printf "$(YELLOW)go.mod not found. Running 'go mod init $(MODULE_NAME)'...$(RESET)\n"; \
		$(GOCMD) mod init $(MODULE_NAME); \
		printf "$(GREEN)✓ go.mod created$(RESET)\n"; \
	fi

# Compilation du binaire
$(NAME): ensure-module
	@printf "$(CYAN)Building $(NAME)...$(RESET)"
	@cd $(SRCS_DIR) && $(GOBUILD) -o ../$(NAME)
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) created!$(RESET)\n"

# Nettoyage (objets go + cache)
clean:
	@printf "$(YELLOW)Cleaning Go cache...$(RESET)"
	@$(GOCLEAN) -cache
	@printf "$(CLEAR)$(GREEN)✓ Go cache cleaned!$(RESET)\n"

fclean: clean
	@printf "$(YELLOW)Deleting $(NAME)...$(RESET)"
	@rm -f $(NAME)
	@rm -f go.mod
	@printf "$(CLEAR)$(GREEN)✓ $(NAME) deleted!$(RESET)\n"

# Reconstruction
re: fclean all

# Lancer le programme (optionnel)
run: all
	@./$(NAME)
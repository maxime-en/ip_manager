# Nom de l'application
APP_NAME := ip_manager

# Répertoire de sortie pour les fichiers binaires
BUILD_DIR := build

# Répertoire de sortie pour les fichiers binaires de développement
DEV_BUILD_DIR := $(BUILD_DIR)/dev

# Répertoire de sortie pour les fichiers binaires de release
RELEASE_BUILD_DIR := $(BUILD_DIR)/release

# Numéro de version pour les releases mineures
VERSION := "0.1"

# Nom du binaire pour l'environnement de développement
DEV_BINARY := $(DEV_BUILD_DIR)/$(APP_NAME)

# Nom du binaire pour l'environnement de release
RELEASE_BINARY := $(RELEASE_BUILD_DIR)/$(APP_NAME)_$(VERSION)

# Cible par défaut
.PHONY: all
all: dev

# Compile l'application pour l'environnement de développement
.PHONY: dev
dev:
	@echo "Building for development..."
	@go build -o $(DEV_BINARY)

# Compile l'application pour l'environnement de release
.PHONY: release
release:
	@echo "Building for release..."
	@go build -o $(RELEASE_BINARY)

# Nettoie les fichiers de build
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(DEV_BUILD_DIR) $(RELEASE_BUILD_DIR)

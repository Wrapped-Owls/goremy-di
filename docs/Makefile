.PHONY: build serve theme theme-check

THEME_DIR := themes/hugo-geekdoc

# --------------------
# Public targets
# --------------------

build: theme-check
	hugo

serve: theme-check
	hugo server

theme:
	@echo "Downloading Geekdoc theme..."
	@mkdir -p $(THEME_DIR)
	@curl -L https://github.com/thegeeklab/hugo-geekdoc/releases/latest/download/hugo-geekdoc.tar.gz \
		| tar -xz -C $(THEME_DIR) --strip-components=1
	@echo "Theme downloaded successfully!"

# --------------------
# Internal
# --------------------

theme-check:
	@if [ ! -d "$(THEME_DIR)" ]; then \
		$(MAKE) theme; \
	fi

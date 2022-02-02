BIN_DIR := ./bin
BIN := $(BIN_DIR)/taskrader
INDEX_MIN_HTML := ./assets/index.min.html
MAIN_MIN_JS := ./assets/main.min.js

.PHONY:	build

build:	$(INDEX_MIN_HTML) $(MAIN_MIN_JS)
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) -ldflags='-s -w' ./cmd/taskrader

$(INDEX_MIN_HTML):	assets/index.html node_modules
	yarn run html-minifier --minify-css --collapse-whitespace --remove-comments $< -o $@

$(MAIN_MIN_JS):	assets/main.js node_modules
	yarn run uglifyjs $< -o $@

node_modules:
	yarn

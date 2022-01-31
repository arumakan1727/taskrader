BIN_DIR := ./bin
BIN := $(BIN_DIR)/taskrader
INDEX_MIN_HTML := ./assets/index.min.html
MAIN_MIN_JS := ./assets/main.min.js

all:	$(BIN)

$(BIN):	$(INDEX_MIN_HTML) $(MAIN_MIN_JS)
	@mkdir -p $(dir $@)
	go build -o $@ ./cmd/taskrader

$(INDEX_MIN_HTML):	assets/index.html yarn.lock
	yarn run html-minifier --minify-css --collapse-whitespace --remove-comments $< -o $@

$(MAIN_MIN_JS):	assets/main.js yarn.lock
	yarn run uglifyjs $< -o $@

yarn.lock:
	yarn && touch $@

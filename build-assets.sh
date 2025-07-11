#!/bin/bash

set -e

echo "Bundling and minifying CSS..."
esbuild web/static/css/main.css --bundle --minify --outfile=web/static/css/style.css

echo "Bundling and minifying JavaScript..."
esbuild web/static/js/main.js --bundle --minify --outfile=web/static/js/bundle.js

echo "Build complete."

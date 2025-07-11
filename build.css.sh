#!/bin/bash

set -e

echo "Bundling and minifying CSS with esbuild..."

esbuild web/static/css/main.css --bundle --minify --outfile=web/static/css/style.css

echo "Successfully created web/static/css/style.css"

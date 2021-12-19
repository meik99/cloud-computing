#!/usr/bin/env bash

# exit on error
set -eu

cd repository/demo-app

echo "Installing npm dependencies"
npm install

echo "Installing Angular CLI"
npm install -g @angular/cli

echo "Running Angular tests"
ng test --watch=false --browser=ChromeHeadless

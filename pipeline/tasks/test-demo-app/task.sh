#!/usr/bin/env bash

# exit on error
set -eu

cd repository/demo-app

echo "Installing npm dependencies"
npm install

echo "Running Angular tests"
ng test --watch=false --browser=ChromeHeadless


#!/usr/bin/env bash

# exit on error
set -eu

cd repository/demo-app

echo "Installing npm dependencies"
npm install

echo "Installing Angular CLI"
npm install -g @angular/cli

echo "Installing Chromium"
apt-get update
apt-get install chromium -y
export CHROME_BIN="$(which chromium)"


echo "Running Angular tests"
ng test --watch=false --browsers=ChromeHeadless


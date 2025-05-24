#!/bin/bash

if ! command -v brew &> /dev/null then
    echo "⚠️ Homebrew not found. Installing..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

    echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
    eval "$(/opt/homebrew/bin/brew shellenv)"
else
    echo "✅ Homebrew is already installed. installing deps"
    brew install go pre-commit
	pre-commit install
    pre-commit install --hook-type pre-push
    go mod tidy
    echo "✅ dependencies installed successfully"
fi

build:
  go build -o ./dl .
  sudo ln -sf $PWD/dl /usr/local/bin/dl
  mkdir -p ~/.local/share/bash-completion/completions
  ./dl completion bash > ~/.local/share/bash-completion/completions/dl

# dl

Install and update packages from git repositories with one tool.

## Usage

### Setup

```bash
dl edit
```

Example configuration:

```yaml
dl:
  - name: neovim
    url: https://github.com/neovim/neovim.git
    dest: /home/user/software/neovim
    build:
      - sudo make CMAKE_BUILD_TYPE=Release
      - sudo make install
```

### Install

```bash
dl install neovim
```

### Update

```bash
dl update neovim
```

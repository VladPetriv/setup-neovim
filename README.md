# setup-nvim

setup-nvim is a simple cli utility for setup your neovim by configs from open source repositories.


## Features 🤔

- Copy neovim configs from Github and GitLab by `HTTP` or `SSH` links.
- Auto extract neovim config from repository if it's not main directory.
- Install package managers [Packer, VimPlug].
- Validation for URL.
- Validation for repository files.
- Colored input and output.
- Check if config already exists and ask permission for deleting it.
- Detect already installed package managers and ask permission for deleting them.

## How to install ✅?

1. Clone the repository:

```bash
  git clone git@github.com:VladPetriv/setup-neovim.git
```

2. Go to repository directory:

```bash
  cd setup-neovim
```

3. Run the installation command via `make`:

```bash
  make install
```
4. Start using:
```bash
  setup-nvim
```

## How to uninstall ❌?

1. Go to repository directory:

```bash
  cd setup-neovim
```

2. Run the uninstalling command via `make`:

```bash
  make uninstall
```



# Compat shim so `nix-shell` uses the same devShell as `nix develop`.
(import (fetchTarball "https://github.com/edolstra/flake-compat/archive/master.tar.gz") {
  src = ./.;
}).shellNix.default

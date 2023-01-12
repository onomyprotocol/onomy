# Steps to prepare the binaries to start the indexing from the first block.

The `bin.sh` or `bin-from-sources.sh` scripts install just the last binary and cosmovisor, but in order to start from
first block
it's required to use the cosmovisor and all binaries.

## How to prepare all binaries

* Find all binaries versions.
  The genesis binary the version is `v1.0.0`. All onomy upgrade versions can be found of
  that [page](https://github.com/onomyprotocol/onomy/tree/main/app/upgrades).
  If you use the compiled binaries you can download them form the release or clone the repository and build it manually.
  Each release
  version corresponds the git tag.

* Now create the corresponding folder structure for the cosmovisor

├── current -> genesis or upgrades/<name>    # This link is created by cosmovisor
├── genesis
│  └── bin
│    └── onomyd  <- use the genesis version here (v1.0.0)
├── upgrades
│  └── v1.0.1
│    └── bin
│      └── onomyd
│  └── <each next version version>
│    └── bin
│      └── onomyd

* For each binary check the version by the command `onomyd version` it mast correspond the upgrade folder.
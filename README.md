Welcome to the WoW Cataclysm Classic simulator! If you have questions or are thinking about contributing, [join our discord](https://discord.gg/jJMPr9JWwx 'https://discord.gg/jJMPr9JWwx') to chat!

The primary goal of this project is to provide a framework that makes it easy to build a DPS sim for any class/spec, with a polished UI and accurate results. Each community will have ownership / responsibility over their portion of the sim, to ensure accuracy and that their community is represented. By having all the individual sims on the same engine, we can also have a combined 'raid sim' for testing raid compositions.

This project is licensed with MIT license. We request that anyone using this software in their own project to make sure there is a user visible link back to the original project.

[Live sims can be found here.](https://wowsims.github.io/cata 'https://wowsims.github.io/cata')

[Support our devs via Patreon.](https://www.patreon.com/wowsims)

# Downloading Sim

Links for latest Sim build:

-   [Windows Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-windows.exe.zip)
-   [MacOS Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-amd64-darwin.zip)
-   [Linux Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-amd64-linux.zip)

Then unzip the downloaded file, then open the unzipped file to open the sim in your browser!

Alternatively, you can choose from a specific relase on the [Releases](https://github.com/wowsims/cata/releases) page and click the suitable link under "Assets"

# Local Dev Installation

This project has dependencies on Go >=1.23, protobuf-compiler and the corresponding Go plugins, and node >= 20.

## Ubuntu

Do not use apt to install any dependencies, the versions they install are all too old.
Script below will curl latest versions and install them.

```sh
# Standard Go installation script
curl -O https://dl.google.com/go/go1.23.4.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bashrc
echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.bashrc
source $HOME/.bashrc

cd cata

# Install protobuf compiler and Go plugins
sudo apt update && sudo apt upgrade
sudo apt install protobuf-compiler
go get -u -v google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install node
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
nvm install 20.13.1

# Install the npm package dependencies using node
npm install
```

## Docker

Alternatively, install Docker and your workflow will look something like this:

```sh
git clone https://github.com/wowsims/cata.git
cd cata

# Build the docker image and install npm dependencies (only need to run these once).
docker build --tag wowsims-cata .
docker run --rm -v $(pwd):/cata wowsims-cata npm install

# Now you can run the commands as shown in the Commands sections, preceding everything with, "docker run --rm -it -p 8080:8080 -v $(pwd):/cata wowsims-cata".
# For convenience, set this as an environment variable:
CATA_CMD="docker run --rm -it -p 8080:8080 -v $(pwd):/cata wowsims-cata"

#For the watch commands assign this environment variable:
CATA_WATCH_CMD="docker run --rm -it -p 8080:8080 -p 3333:3333 -p 5173:5173 -e WATCH=1 -v $(pwd):/cata wowsims-cata"

# ... do some coding on the sim ...

# Run tests
$(echo $CATA_CMD) make test

# ... do some coding on the UI ...

# Host a local site
$(echo $CATA_CMD) make host
```

## Windows

If you want to develop on Windows, we recommend setting up a Ubuntu virtual machine (VM) or running Docker using [this guide](https://docs.docker.com/desktop/windows/wsl/ 'https://docs.docker.com/desktop/windows/wsl/') and then following the Ubuntu or Docker instructions, respectively.

## Mac OS

-   Docker is available in OS X as well, so in theory similar instructions should work for the Docker method
-   You can also use the Ubuntu setup instructions as above to run natively, with a few modifications:
    -   You may need a different Go installer if `go1.18.3.linux-amd64.tar.gz` is not compatible with your system's architecture; you can do the Go install manually from `https://go.dev/doc/install`.
    -   OS X uses Homebrew instead of apt, so in order to install protobuf-compiler you’ll instead need to run `brew install protobuf-c` (note the package name is also a little different than in apt). You might need to first update or upgrade brew.
    -   The provided install script for Node will not included a precompiled binary for OS X, but it’s smart enough to compile one. Be ready for your CPU to melt on running `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash`.

# Commands

We use a makefile for our build system. These commands will usually be all you need while developing for this project:

```sh
# Installs a pre-commit git hook so that your go code is automatically formatted (if you don't use an IDE that supports that).  If you want to manually format go code you can run make fmt.
# Also installs `air` to reload the dev servers automatically
make setup

# Run all the tests. Currently only the backend sim has tests.
make test

# Update the expected test results. This will need to be run after adding/removing any tests, and also if test results change due to code changes.
make update-tests

# Host a local version of the UI at http://localhost:8080. Visit it by pointing a browser to
# http://localhost:8080/cata/YOUR_SPEC_HERE, where YOUR_SPEC_HERE is the directory under ui/ with your custom code.
# Recompiles the entire client before launching using `make dist/cata`
make host

# With file-watching so the server auto-restarts and recompiles on Go or TS changes:
WATCH=1 make host

# Delete all generated files (.pb.go and .ts proto files, and dist/)
make clean

# Recompiles the ts only for the given spec (e.g. make host_elemental_shaman)
make host_$spec

# Recompiles the `wowsimcata` server binary and runs it, hosting /dist directory at http://localhost:3333/cata.
# This is the fastest way to iterate on core go simulator code so you don't have to wait for client rebuilds.
# To rebuild client for a spec just do 'make $spec' and refresh browser.
make rundevserver

# With file-watching so the server auto-restarts and recompiles on Go or TS changes:
WATCH=1 make rundevserver


# The same as rundevserver, recompiles  `wowsimcata` binary and runs it on port 3333. Instead of serving content from the dist folder,
# this command also runs `vite serve` to start the Vite dev server on port 5173 (or similar) and automatically reloads the page on .ts changes in less than a second.
# This allows for more rapid development, with sub second reloads on TS changes. This combines the benefits of `WATCH=1 make rundevserver` and `WATCH=1 make host`
# to create something that allows you to work in any part of the code with ease and speed.
# This might get rolled into `WATCH=1 make rundevserver` at some point.
WATCH=1 make devmode

# This is just the same as rundevserver currently
make devmode

# This command recompiles the workers in the /ui/worker folder for easier debugging/development
# Can be used with or without WATCH command
make webworkers

# With file watch enabled
WATCH=1 make webworkers

# Creates the 'wowsimcata' binary that can host the UI and run simulations natively (instead of with wasm).
# Builds the UI and the compiles it into the binary so that you can host the sim as a server instead of wasm on the client.
# It does this by first doing make dist/cata and then copying all those files to binary_dist/cata and loading all the files in that directory into its binary on compile.
make wowsimcata

# Using the --usefs flag will instead of hosting the client built into the binary, it will host whatever code is found in the /dist directory.
# Use --wasm to host the client with the wasm simulator.
# The server also disables all caching so that refreshes should pickup any changed files in dist/. The client will still call to the server to run simulations so you can iterate more quickly on client changes.
# make dist/cata && ./wowsimcata --usefs would rebuild the whole client and host it. (you would have had to run `make devserver` to build the wowsimcata binary first.)
./wowsimcata --usefs

# Generate code for the sim database (db.json). Only necessary if you changed the items generator.
# Useful only if you're actively working on the generator and have already run make db locally at least once.
make simdb

# Generate data from WoW client files
# Requires dotnet 9 to run
# Uses tools/database/generator-settings.json for settings
# Also runs make simdb
# This is what you will use most of the time for generation
make db

# Same as make db but from the ptr client
# Uses tools/database/ptr-generator-settings.json for settings
make ptrdb
```

## (Optional) Installing Dotnet 9 - Required if generating client data

```sh
curl -L https://dot.net/v1/dotnet-install.sh -o dotnet-install.sh
chmod +x ./dotnet-install.sh
./dotnet-install.sh --channel 9.0
echo 'export PATH=$PATH:$HOME/.dotnet' >> ~/.bashrc
source ~/.bashrc
```

# Adding a Sim

So you want to make a new sim for your class/spec! The basic steps are as follows:

-   [Create the proto interface between sim and UI.](#create-the-proto-interface-between-sim-and-ui)
-   [Implement the UI.](#implement-the-ui)
-   [Implement the sim.](#implement-the-sim)
-   [Launch the site.](#launch-the-site)

## Create the proto interface between Sim and UI

This project uses [Google Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial 'https://developers.google.com/protocol-buffers/docs/gotutorial') to pass data between the sim and the UI. TLDR; Describe data structures in .proto files, and the tool can generate code in any programming language. It lets us avoid repeating the same code in our Go and Typescript worlds without losing type safety.

For a new sim, make the following changes:

-   Add a new value to the `Spec` enum in proto/common.proto. **NOTE: The name you give to this enum value is not just a name, it is used in our templating system. This guide will refer to this name as `$SPEC` elsewhere.**
-   Add a 'proto/YOUR_CLASS.proto' file if it doesn't already exist and add data messages containing all the class/spec-specific information needed to run your sim.
-   Update the `PlayerOptions.spec` field in `proto/api.proto` to include your shiny new message as an option.

That's it! Now when you run `make` there will be generated .go and .ts code in `sim/core/proto` and `ui/core/proto` respectively. If you aren't familiar with protos, take a quick look at them to see what's happening.

## Implement the UI

The UI and sim can be done in either order, but it is generally recommended to build the UI first because it can help with debugging. The UI is very generalized and it doesn't take much work to build an entire sim UI using our templating system. To use it:

-   Modify `ui/core/proto_utils/utils.ts` to include boilerplate for your `$SPEC` name if it isn't already there.
-   Create a directory `ui/$SPEC`. So if your Spec enum value was named, `elemental_shaman`, create a directory, `ui/elemental_shaman`.
-   Copy+paste from another spec's UI code.
-   Modify all the files for your spec; most of the settings are fairly obvious, if you need anything complex just ask and we can help!
-   Finally, add a rule to the `makefile` for the new sim site. Just copy from the other site rules already there and change the `$SPEC` names.

No .html is needed, it will be generated based on `ui/index_template.html` and the `$SPEC` name.

When you're ready to try out the site, run `make host` and navigate to `http://localhost:8080/cata/$SPEC`.

## Implement the Sim

This step is where most of the magic happens. A few highlights to start understanding the sim code:

-   `sim/wasm/main.go` This file is the actual main function, for the [.wasm binary](https://webassembly.org/ 'https://webassembly.org/') used by the UI. You shouldn't ever need to touch this, but just know its here.
-   `sim/core/api.go` This is where the action starts. This file implements the request/response messages defined in `proto/api.proto`.
-   `sim/core/sim.go` Orchestrates everything. Main event loop is in `Simulation.RunOnce`.
-   `sim/core/agent.go` An Agent can be thought of as the 'Player', i.e. the person controlling the game. This is the interface you'll be implementing.
-   `sim/core/character.go` A Character holds all the stats/cooldowns/gear/etc common to any WoW character. Each Agent has a Character that it controls.

Read through the core code and some examples from other classes/specs to get a feel for what's needed. Hopefully `sim/core` already includes what you need, but most classes have at least 1 unique mechanic so you may need to touch `core` as well.

Finally, add your new sim to `RegisterAll()` in `sim/register_all.go`.

Don't forget to write unit tests! Again, look at existing tests for examples. Run them with `make test` when you're ready.

# Launch the site

When everything is ready for release, modify `ui/core/launched_sims.ts` and `ui/index.html` to include the new spec value. This will add the sim to the dropdown menu so anyone can find it from the existing sims. This will also remove the UI warning that the sim is under development. Now tell everyone about your new sim!

# Add your spec to the raid sim

Don't touch the raid sim until the individual sim is ready for launch; anything in the raid sim is publicly accessible. To add your new spec to the raid sim, do the following:

-   Add a reference to the individual sim in `ui/raid/tsconfig.json`. DO NOT FORGET THIS STEP or Typescipt will silently do very bad things.
-   Import the individual sim's css file from `ui/raid/index.scss`.
-   Update `ui/raid/presets.ts` to include a constructor factory in the `specSimFactories` variable and add configurations for new Players in the `playerPresets` variable.

# Deployment

Thanks to the workflow defined in `.github/workflows/deploy.yml`, pushes to `master` automatically build and deploy a new site so there's nothing to do here. Sit back and appreciate your new sim!

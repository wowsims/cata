Welcome to the WoW Cataclysm Classic simulator! If you have questions or are thinking about contributing, [join our discord](https://discord.gg/jJMPr9JWwx 'https://discord.gg/jJMPr9JWwx') to chat!

The primary goal of this project is to provide a framework that makes it easy to build a DPS sim for any class/spec, with a polished UI and accurate results. Each community will have ownership / responsibility over their portion of the sim, to ensure accuracy and that their community is represented. By having all the individual sims on the same engine, we can also have a combined 'raid sim' for testing raid compositions.

This project is licensed with MIT license. We request that anyone using this software in their own project to make sure there is a user visible link back to the original project.

[Live sims can be found here.](https://wowsims.github.io/cata)

[Support our devs via Patreon.](https://www.patreon.com/wowsims)

# Downloading Sim

## Latest Sim Builds

-   [Windows Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-windows.exe.zip)
-   [MacOS Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-amd64-darwin.zip)
-   [Linux Sim](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-amd64-linux.zip)
-   [Linux Sim for ARM](https://github.com/wowsims/cata/releases/latest/download/wowsimcata-amd64-linux.zip)

Unzip the downloaded file and open the unzipped file to launch the sim in your browser!

Alternatively, you can select a specific release on the [Releases](https://github.com/wowsims/cata/releases) page and click the suitable link under "Assets".

Downloading and running the sim locally provides better performance and allows for much higher simulation iteration counts compared to running it on the live site.

## Using Docker

We publish a Docker image for the sim on each release. The image is available in amd64 and arm64 architectures.

### Running the Docker Image

To run the sim in a Docker container and expose it on port 3333, use the following command:

```sh
docker run -p 3333:3333 ghcr.io/wowsims/cata:latest
```

### Updating the running Docker Image version

To update the running Docker image to the latest version, use the following command:

```sh
docker pull ghcr.io/wowsims/cata:latest
```

Once the image is pulled, you can recreate the container using to use the new version.

## Development Instructions

Please see the [development instructions](https://github.com/wowsims/cata/blob/master/DEVELOPMENT.md) for more information on how to get started with development for the project.

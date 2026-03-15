# Fabric CLI

A command-line tool for generating Fabric mod projects for Minecraft.

## Features

- **Interactive TUI** - Beautiful terminal interface using [charmbracelet/huh](https://github.com/charmbracelet/huh)
- **Real-time Version Fetching** - Automatically fetches latest versions from Fabric Meta and Modrinth APIs
- **Concurrent API Calls** - Fetches version data in parallel for speed
- **Complete Project Generation** - Generates all necessary files:
  - `build.gradle` with proper dependencies
  - `gradle.properties` with mod metadata
  - `settings.gradle` with Fabric plugin repository
  - `fabric.mod.json` with mod configuration
  - `mixins.json` (optional)
  - Main mod class
  - Directory structure

## Installation

```bash
go install github.com/yourusername/fabric-cli/cmd/fabric-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/yourusername/fabric-cli.git
cd fabric-cli
go build -o fabric-cli ./cmd/fabric-cli
```

## Usage

Simply run:

```bash
./fabric-cli
```

The wizard will guide you through:
1. **Version Selection** - Choose Minecraft, Yarn, Loader, and Fabric API versions
2. **Mod Metadata** - Enter mod ID, name, description, group ID, etc.
3. **Additional Options** - Enable mixins, select environment, Java version

## Example

```bash
$ ./fabric-cli
🔧 Fabric Project Generator
Fetching latest versions...
✓ Found 500+ Minecraft versions
✓ Found 2000+ Yarn mappings
✓ Found 150+ Fabric API versions

[Interactive TUI prompts...]

📋 Configuration Summary:
════════════════════════════════════════
Minecraft:    1.21.4
Yarn:         1.21.4+build.8
Loader:       0.16.9
Fabric API:   0.110.0

Mod ID:       mycoolmod
Mod Name:     My Cool Mod
Group ID:     com.example
Version:      1.0.0
Mixins:       true

🚀 Generating project in 'mycoolmod/'...

  ✓ Generated: build.gradle
  ✓ Generated: gradle.properties
  ✓ Generated: settings.gradle
  ✓ Generated: gradle/wrapper/gradle-wrapper.properties
  ✓ Generated: src/main/resources/fabric.mod.json
  ✓ Generated: src/main/resources/mycoolmod.mixins.json
  ✓ Generated: src/main/java/com/example/mycoolmod/Mycoolmod.java

✅ Project generated successfully!

Next steps:
  cd mycoolmod
  ./gradlew build
```

## Project Structure

Generated projects follow standard Fabric mod structure:

```
my-mod/
├── build.gradle
├── gradle.properties
├── settings.gradle
├── gradlew
├── gradlew.bat
├── gradle/
│   └── wrapper/
│       ├── gradle-wrapper.jar
│       └── gradle-wrapper.properties
└── src/
    └── main/
        ├── java/
        │   └── com/example/mymod/
        │       └── MyMod.java
        └── resources/
            ├── fabric.mod.json
            ├── mymod.mixins.json
            └── assets/mymod/
```

## Requirements

- Go 1.21+
- Internet connection (for fetching versions)

## Architecture

This project demonstrates several Go patterns:

- **Concurrent Programming** - Uses goroutines and channels for parallel API fetching
- **Template System** - Go's `text/template` for dynamic file generation
- **TUI Interface** - Interactive forms using charmbracelet/huh
- **Wizard Pattern** - Multi-step configuration with state management

## API Endpoints

- **Fabric Meta API**: `https://meta.fabricmc.net/v2/versions`
- **Modrinth API**: `https://api.modrinth.com/v2/project/P7dR8mSH/version`

## License

MIT

## Acknowledgments

Based on the [MinecraftDev](https://github.com/minecraft-dev/MinecraftDev) IntelliJ plugin's Fabric project generator architecture.

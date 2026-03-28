# Fabric CLI

A command-line tool for generating Fabric mod projects for Minecraft.

## Features

- **Interactive TUI** - Beautiful terminal interface using [charmbracelet/huh](https://github.com/charmbracelet/huh)
- **Quick Mode** - Fast project generation with command-line arguments
- **Type-to-Search** - Search for specific Minecraft versions by typing (e.g., "1.21", "1.20.1")
- **Real-time Version Fetching** - Automatically fetches latest versions from Fabric Meta and Modrinth APIs
- **Concurrent API Calls** - Fetches version data in parallel for speed
- **Auto-Complete** - Automatically selects latest compatible versions and converts mod names to IDs
- **Auto Java Detection** - Automatically detects installed Java versions and recommends the best one for your Minecraft version
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
go install github.com/tantaihaha4487/fabric-cli-go/cmd/fabric-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/tantaihaha4487/fabric-cli-go.git
cd fabric-cli-go
go build -o fabric-cli ./cmd/fabric-cli
```

## Usage

### Interactive Wizard Mode (Default)

Simply run without arguments to start the interactive wizard:

```bash
./fabric-cli
```

The wizard will guide you through:
1. **Version Selection** - Search for Minecraft version, then select Yarn, Loader, and Fabric API versions
2. **Java Version** - Auto-detects installed Java and recommends the best version for your Minecraft version
3. **Mod Metadata** - Enter mod ID, name, description, group ID, etc.
4. **Additional Options** - Enable mixins, select environment

### Quick Mode

For faster project generation, use quick mode with positional arguments:

```bash
./fabric-cli <mc_version> <mod_name> <mod_version> [group_id]
```

**Arguments:**
- `mc_version` - Minecraft version (e.g., 1.21.4, 1.20.1)
- `mod_name` - Display name for your mod (auto-converted to mod ID)
- `mod_version` - Version of your mod (e.g., 1.0.0)
- `group_id` - Maven group ID (optional, uses saved config or 'com.example')

**Quick Mode Options:**
- `--no-mixins` - Disable Mixins support
- `--client-only` - Set environment to client only
- `--server-only` - Set environment to server only
- `--java-version=N` - Set Java version (default: 21)
- `--license=TYPE` - Set license (default: MIT)

**Examples:**

```bash
# Quick mode with defaults
./fabric-cli 1.21.4 "My Cool Mod" 1.0.0

# Quick mode with group ID
./fabric-cli 1.21.4 "My Cool Mod" 1.0.0 com.example

# Quick mode with options
./fabric-cli 1.21.4 "My Cool Mod" 1.0.0 com.example \
  --no-mixins --client-only --java-version=17 --license=Apache-2.0
```

### Help & Version

```bash
./fabric-cli -h          # or --help
./fabric-cli -v          # or --version
```

## Example

```bash
$ ./fabric-cli
[Fabric] Fabric Project Generator
Fetching latest versions...
[OK] Found 483 Minecraft versions
[OK] Found 3410 Yarn mappings
[OK] Found 1015 Fabric API versions

[OK] Loaded 483 Minecraft versions

Type to search versions (or leave empty to show all)
> 1.21

[OK] Found 12 matching versions

[Interactive TUI prompts...]

[Summary] Configuration Summary:
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

[Generate] Generating project in 'mycoolmod/'...

  [OK] Generated: build.gradle
  [OK] Generated: gradle.properties
  [OK] Generated: settings.gradle
  [OK] Generated: gradle/wrapper/gradle-wrapper.properties
  [OK] Generated: src/main/resources/fabric.mod.json
  [OK] Generated: src/main/resources/mycoolmod.mixins.json
  [OK] Generated: src/main/java/com/example/mycoolmod/Mycoolmod.java

[Success] Project generated successfully!

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

## Java Version Compatibility

The tool automatically recommends the appropriate Java version based on your Minecraft version:

| Minecraft Version | Recommended Java |
|-------------------|-----------------|
| 26.x (2026+)      | Java 25         |
| 1.21.x            | Java 21          |
| 1.20.x            | Java 17          |
| 1.19.x            | Java 17          |
| 1.18.x            | Java 17          |
| 1.17.x            | Java 17          |
| 1.16.x            | Java 16          |
| 1.14.x - 1.15.x  | Java 11          |
| 1.8.x - 1.13.x   | Java 8           |

The wizard will detect installed Java versions on your machine and auto-select the recommended version if available.

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

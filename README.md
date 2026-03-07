# Fabric CLI

A minimalist CLI tool to generate Fabric mod projects with interactive prompts.

## Installation

```bash
cd fabric-cli-go
go build -o fabric-cli .
./fabric-cli --help
```

## Usage

### Interactive Mode (Recommended)

Run without arguments to enter interactive mode:

```bash
fabric-cli new -i
```

This will prompt you through:
- Mod ID, Name, Description
- Package Name
- Minecraft Version (with auto-complete)
- Loom Version
- Loader Version
- JDK Version
- Yarn Mappings (auto-set based on MC version)
- Fabric API Version (auto-set based on MC version)
- Environment (Both/Client/Server)
- Language (Java/Kotlin)
- Use Mixins?
- Use Data Generation?
- Use Official Mappings?
- Include Fabric API?
- Author
- License

### Non-Interactive Mode

```bash
fabric-cli new <project-name> [flags]
```

**Flags:**
- `-n, --name string` - Mod name
- `-d, --description string` - Mod description
- `-p, --package string` - Package name (com.example.modid)
- `-a, --author string` - Author name (default "Anonymous")
- `-l, --license string` - License (MIT/Apache-2.0/LGPL-3.0) (default "MIT")
- `-m, --mc-version string` - Minecraft version
- `-j, --java-version string` - Java version
- `-v, --loader-version string` - Fabric loader version
- `--loom-version string` - Loom version
- `-g, --language string` - Language (Java/Kotlin) (default "Java")
- `-e, --environment string` - Environment (*, client, server) (default "*")
- `-t, --template string` - Template (basic/mixin/datagen)
- `-f, --fabric-api` - Include Fabric API (default true)
- `-o, --official-mappings` - Use Mojang official mappings
- `--mixins` - Use Mixins
- `--datagen` - Use Data Generation

**Examples:**

```bash
# Interactive mode
fabric-cli new -i

# Basic usage
fabric-cli new mymod -n "My Mod" -d "Description" -p com.example.mymod

# With Kotlin
fabric-cli new mymod -n "My Mod" -d "Desc" -p com.example.mymod -g Kotlin

# Client-side only
fabric-cli new mymod -n "My Mod" -d "Desc" -p com.example.mymod -e client

# With mixins
fabric-cli new mymod -n "My Mod" -d "Desc" -p com.example.mymod --mixins
```

### Generate Code

```bash
fabric-cli generate <type> <name> [flags]
```

**Types:**
- `block` - Generate a Block class
- `item` - Generate an Item class
- `enchantment` - Generate an Enchantment class
- `effect` - Generate a StatusEffect class

**Examples:**

```bash
# Generate a block
fabric-cli generate block MyBlock -p com.example.mymod

# Generate an item
fabric-cli generate item MyItem -p com.example.mymod
```

## Project Structure

Generated projects include:

```
modid/
├── build.gradle
├── settings.gradle
├── gradle.properties
├── LICENSE
├── gradlew
├── gradlew.bat
└── src/main/
    ├── java/ (or kotlin/)
    │   └── com/example/modid/
    │       └── ModidMod.java
    └── resources/
        ├── fabric.mod.json
        └── modid.mixins.json (if using mixins)
```

## Interactive Features

The interactive mode features:
- **Colorful UI** - Modern terminal colors using ANSI codes
- **Smart Defaults** - Pre-selected common options
- **Auto-complete** - Version suggestions based on Minecraft version
- **Summary View** - Review all options before creating

## Next Steps

After creating a project:

```bash
cd <project-name>
./gradlew build
```

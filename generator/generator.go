package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
)

// Generator handles project file generation
type Generator struct {
	ctx *project.Context
}

// NewGenerator creates a new generator instance
func NewGenerator(ctx *project.Context) *Generator {
	return &Generator{ctx: ctx}
}

// Generate creates the project files
func (g *Generator) Generate(projectPath string) error {
	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	data := buildTemplateData(g.ctx)

	// Generate build.gradle
	if err := g.generateFile(projectPath, "build.gradle", buildGradleTemplate, data); err != nil {
		return err
	}

	// Generate gradle.properties
	if err := g.generateFile(projectPath, "gradle.properties", gradlePropertiesTemplate, data); err != nil {
		return err
	}

	// Generate settings.gradle
	if err := g.generateFile(projectPath, "settings.gradle", settingsGradleTemplate, data); err != nil {
		return err
	}

	// Generate gradle wrapper properties
	wrapperDir := filepath.Join(projectPath, "gradle", "wrapper")
	if err := os.MkdirAll(wrapperDir, 0755); err != nil {
		return fmt.Errorf("failed to create wrapper directory: %w", err)
	}
	if err := g.generateFile(projectPath, "gradle/wrapper/gradle-wrapper.properties", wrapperPropertiesTemplate, data); err != nil {
		return err
	}

	// Generate gradlew scripts
	if err := g.generateWrapperScript(projectPath, "gradlew", gradlewTemplate, 0755); err != nil {
		return err
	}
	if err := g.generateWrapperScript(projectPath, "gradlew.bat", gradlewBatTemplate, 0644); err != nil {
		return err
	}

	// Download gradle-wrapper.jar
	fmt.Println("  [...] Downloading gradle-wrapper.jar...")
	wrapperJarURL := fmt.Sprintf("https://raw.githubusercontent.com/gradle/gradle/v%s/gradle/wrapper/gradle-wrapper.jar", data.GradleVersion)
	wrapperJarPath := filepath.Join(projectPath, "gradle", "wrapper", "gradle-wrapper.jar")
	if err := g.downloadFile(wrapperJarURL, wrapperJarPath); err != nil {
		fmt.Printf("  [!]  Warning: Could not download gradle-wrapper.jar: %v\n", err)
		fmt.Println("  You can manually run 'gradle wrapper' after installing Gradle")
	} else {
		fmt.Printf("  [OK] Downloaded: gradle/wrapper/gradle-wrapper.jar\n")
	}

	// Generate fabric.mod.json
	resourcesDir := filepath.Join(projectPath, "src", "main", "resources")
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return fmt.Errorf("failed to create resources directory: %w", err)
	}
	if err := g.generateFile(projectPath, "src/main/resources/fabric.mod.json", fabricModJsonTemplate, data); err != nil {
		return err
	}

	// Generate mixins.json if enabled
	if g.ctx.UseMixins {
		if err := g.generateFile(projectPath, "src/main/resources/"+g.ctx.ModID+".mixins.json", mixinsJsonTemplate, data); err != nil {
			return err
		}
	}

	// Create source directory structure
	packagePath := strings.ReplaceAll(g.ctx.GroupID, ".", "/")
	javaDir := filepath.Join(projectPath, "src", "main", "java", packagePath, g.ctx.ModID)
	if err := os.MkdirAll(javaDir, 0755); err != nil {
		return fmt.Errorf("failed to create java directory: %w", err)
	}

	// Create assets directory
	assetsDir := filepath.Join(projectPath, "src", "main", "resources", "assets", g.ctx.ModID)
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Generate main mod class
	if err := g.generateFile(projectPath, "src/main/java/"+packagePath+"/"+g.ctx.ModID+"/"+toClassName(g.ctx.ModID)+".java", modClassTemplate, data); err != nil {
		return err
	}

	return nil
}

// Template definitions
const buildGradleTemplate = `plugins {
    id 'net.fabricmc.fabric-loom' version '{{ .LoomVersion }}'
    id 'maven-publish'
}

version = project.mod_version
group = project.maven_group

base {
    archivesName = project.archives_base_name
}

repositories {
    // Add repositories to retrieve artifacts from in here.
}

dependencies {
    // To change the versions see the gradle.properties file
    minecraft "com.mojang:minecraft:${project.minecraft_version}"
{{if not .ImplicitMappings}}
{{if .OfficialMappings}}
    mappings loom.officialMojangMappings()
{{else}}
    mappings "net.fabricmc:yarn:${project.yarn_mappings}:v2"
{{end}}
{{end}}
    {{ .DependencyConfig }} "net.fabricmc:fabric-loader:${project.loader_version}"

{{if .APIVersion}}
    // Fabric API. This is technically optional, but you probably want it anyway.
    {{ .DependencyConfig }} "net.fabricmc.fabric-api:fabric-api:${project.fabric_version}"
{{end}}
}

processResources {
    inputs.property "version", project.version
    inputs.property "minecraft_version", project.minecraft_version
    inputs.property "loader_version", project.loader_version
    filteringCharset "UTF-8"

    filesMatching("fabric.mod.json") {
        expand "version": project.version,
                "minecraft_version": project.minecraft_version,
                "loader_version": project.loader_version
    }
}

def targetJavaVersion = {{ .JavaVersion }}
tasks.withType(JavaCompile).configureEach {
    it.options.encoding = "UTF-8"
    if (targetJavaVersion >= 10 || JavaVersion.current().isJava10Compatible()) {
        it.options.release.set(targetJavaVersion)
    }
}

java {
    def javaVersion = JavaVersion.toVersion(targetJavaVersion)
    if (JavaVersion.current() < javaVersion) {
        toolchain.languageVersion = JavaLanguageVersion.of(targetJavaVersion)
    }
    withSourcesJar()
}

jar {
    from("LICENSE") {
        rename { "${it}_${project.archivesBaseName}"}
    }
}

publishing {
    publications {
        create("mavenJava", MavenPublication) {
            artifactId = project.archives_base_name
            from components.java
        }
    }
    repositories {
        // Add repositories to publish to here.
    }
}
`

const gradlePropertiesTemplate = `# Done to increase the memory available to gradle.
org.gradle.jvmargs=-Xmx1G
{{if .ImplicitMappings}}
org.gradle.parallel=true
org.gradle.configuration-cache=false
{{end}}

# Fabric Properties
# check these on https://fabricmc.net/develop
minecraft_version={{ .MCVersion }}
{{if and (not .ImplicitMappings) (not .OfficialMappings)}}yarn_mappings={{ .YarnMappings }}
{{end}}loader_version={{ .LoaderVersion }}

# Mod Properties
mod_version = {{ .Version }}
maven_group = {{ .GroupID }}
archives_base_name = {{ .ArtifactID }}

{{if .APIVersion}}# Dependencies
# check this on https://fabricmc.net/develop
fabric_version={{ .APIVersion }}
{{end}}
`

const settingsGradleTemplate = `pluginManagement {
    repositories {
        maven {
            name = 'Fabric'
            url = 'https://maven.fabricmc.net/'
        }
        mavenCentral()
        gradlePluginPortal()
    }
}
`

const wrapperPropertiesTemplate = `distributionUrl=https\://services.gradle.org/distributions/gradle-{{ .GradleVersion }}-bin.zip
`

const fabricModJsonTemplate = `{
  "schemaVersion": 1,
  "id": "{{ .ModID }}",
  "version": "${version}",

  "name": "{{ .ModName }}",
  "description": "{{ .ModDescription }}",
  "authors": [],
  "contact": {},

  "license": "{{ .License }}",
  "icon": "assets/{{ .ModID }}/icon.png",

  "environment": "{{ .ModEnvironment }}",
  "entrypoints": {
    "main": [
      "{{ .GroupID }}.{{ .ModID }}.{{ .ModID | ToClassName }}"
    ]
  },
{{if .Mixins}}  "mixins": [
    "{{ .ModID }}.mixins.json"
  ],
{{end}}  "depends": {
    "fabricloader": ">=${loader_version}",
    "java": ">={{ .JavaVersion }}",
{{if .APIVersion}}{{if .ImplicitMappings}}    "{{ .APIDependencyKey }}": "*",
{{else}}    "fabric": "*",
{{end}}
{{end}}    "minecraft": "~${minecraft_version}"
  }
}
`

const mixinsJsonTemplate = `{
  "required": true,
  "minVersion": "0.8",
  "package": "{{ .MixinPackageName }}",
  "compatibilityLevel": "JAVA_{{ .JavaVersion }}",
  "mixins": [],
  "client": [],
{{if .ImplicitMappings}}  "overwrites": {
    "requireAnnotations": true
  },
{{end}}  "injectors": {
    "defaultRequire": 1
  }
}
`

const modClassTemplate = `package {{ .GroupID }}.{{ .ModID }};

import net.fabricmc.api.ModInitializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class {{ .ModID | ToClassName }} implements ModInitializer {
	public static final Logger LOGGER = LoggerFactory.getLogger("{{ .ModID }}");

	@Override
	public void onInitialize() {
		LOGGER.info("Hello from {{ .ModName }}!");
	}
}
`

const gradlewTemplate = "#!/bin/sh\n\n#\n# Copyright 2015 the original author or authors.\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#      https://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n#\n\n##############################################################################\n#\n#   Gradle start up script for POSIX generated by Gradle.\n#\n##############################################################################\n\n# Attempt to set APP_HOME\n\n# Resolve links: $0 may be a link\napp_path=$0\n\n# Need this for daisy-chained symlinks.\nwhile\n    APP_HOME=${app_path%\"${app_path##*/}\"}  # leaves a trailing /; empty if no leading path\n    [ -h \"$app_path\" ]\ndo\n    ls=$( ls -ld \"$app_path\" )\n    link=${ls#*' -> '}\n    case $link in             #(\n      /*)   app_path=$link ;; #(\n      *)    app_path=$APP_HOME$link ;;\n    esac\ndone\n\nAPP_BASE_NAME=${0##*/}\nAPP_HOME=$( cd \"${APP_HOME:-./}\" && pwd -P ) || exit\n\n# Use the maximum available, or set MAX_FD != -1 to use that value.\nMAX_FD=maximum\n\nwarn () {\n    echo \"$*\"\n} >&2\n\ndie () {\n    echo\n    echo \"$*\"\n    echo\n    exit 1\n} >&2\n\n# OS specific support (must be 'true' or 'false').\ncygwin=false\nmsys=false\ndarwin=false\nnonstop=false\ncase \"$( uname )\" in                #(\n  CYGWIN* )         cygwin=true  ;; #(\n  Darwin* )         darwin=true  ;; #(\n  MSYS* | MINGW* )  msys=true    ;; #(\n  NONSTOP* )        nonstop=true ;;\nesac\n\nCLASSPATH=$APP_HOME/gradle/wrapper/gradle-wrapper.jar\n\n\n# Determine the Java command to use to start the JVM.\nif [ -n \"$JAVA_HOME\" ] ; then\n    if [ -x \"$JAVA_HOME/jre/sh/java\" ] ; then\n        JAVACMD=$JAVA_HOME/jre/sh/java\n    else\n        JAVACMD=$JAVA_HOME/bin/java\n    fi\n    if [ ! -x \"$JAVACMD\" ] ; then\n        die \"ERROR: JAVA_HOME is set to an invalid directory: $JAVA_HOME\n\nPlease set the JAVA_HOME variable in your environment to match the\nlocation of your Java installation.\"\n    fi\nelse\n    JAVACMD=java\n    which java >/dev/null 2>&1 || die \"ERROR: JAVA_HOME is not set and no 'java' command could be found in your PATH.\n\nPlease set the JAVA_HOME variable in your environment to match the\nlocation of your Java installation.\"\nfi\n\n# Increase the maximum file descriptors if we can.\nif ! \"$cygwin\" && ! \"$darwin\" && ! \"$nonstop\" ; then\n    case $MAX_FD in #(\n      max*)\n        MAX_FD=$( ulimit -H -n ) ||\n            warn \"Could not query maximum file descriptor limit\"\n    esac\n    case $MAX_FD in  #(\n      '' | soft) :;; #(\n      *)\n        ulimit -n \"$MAX_FD\" ||\n            warn \"Could not set maximum file descriptor limit to $MAX_FD\"\n    esac\nfi\n\n# Collect all arguments for the java command;\nset -- \\\n        \"-Dorg.gradle.appname=$APP_BASE_NAME\" \\\n        -classpath \"$CLASSPATH\" \\\n        org.gradle.wrapper.GradleWrapperMain \\\n        \"$@\"\n\nexec \"$JAVACMD\" \"$@\"\n"

const gradlewBatTemplate = `@rem
@rem Copyright 2015 the original author or authors.
@rem
@rem Licensed under the Apache License, Version 2.0 (the "License");
@rem you may not use this file except in compliance with the License.
@rem You may obtain a copy of the License at
@rem
@rem      https://www.apache.org/licenses/LICENSE-2.0
@rem
@rem Unless required by applicable law or agreed to in writing, software
@rem distributed under the License is distributed on an "AS IS" BASIS,
@rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@rem See the License for the specific language governing permissions and
@rem limitations under the License.
@rem

@if "%DEBUG%" == "" @echo off
@rem ##########################################################################
@rem
@rem  Gradle startup script for Windows
@rem
@rem ##########################################################################

@rem Set local scope for the variables with windows NT shell
if "%OS%"=="Windows_NT" setlocal

set DIRNAME=%~dp0
if "%DIRNAME%" == "" set DIRNAME=.
set APP_BASE_NAME=%~n0
set APP_HOME=%DIRNAME%

@rem Resolve any "." and ".." in APP_HOME to make it shorter.
for %%i in ("%APP_HOME%") do set APP_HOME=%%~fi

@rem Add default JVM options here. You can also use JAVA_OPTS and GRADLE_OPTS to pass JVM options to this script.
set DEFAULT_JVM_OPTS="-Xmx64m" "-Xms64m"

@rem Find java.exe
if defined JAVA_HOME goto findJavaFromJavaHome

set JAVA_EXE=java.exe
%JAVA_EXE% -version >NUL 2>&1
if "%ERRORLEVEL%" == "0" goto execute

echo.
echo ERROR: JAVA_HOME is not set and no 'java' command could be found in your PATH.
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:findJavaFromJavaHome
set JAVA_HOME=%JAVA_HOME:"=%
set JAVA_EXE=%JAVA_HOME%/bin/java.exe

if exist "%JAVA_EXE%" goto execute

echo.
echo ERROR: JAVA_HOME is set to an invalid directory: %JAVA_HOME%
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:execute
@rem Setup the command line

set CLASSPATH=%APP_HOME%\gradle\wrapper\gradle-wrapper.jar


@rem Execute Gradle
"%JAVA_EXE%" %DEFAULT_JVM_OPTS% %JAVA_OPTS% %GRADLE_OPTS% "-Dorg.gradle.appname=%APP_BASE_NAME%" -classpath "%CLASSPATH%" org.gradle.wrapper.GradleWrapperMain %*

:end
@rem End local scope for the variables with windows NT shell
if "%ERRORLEVEL%"=="0" goto mainEnd

:fail
rem Set variable GRADLE_EXIT_CONSOLE if you need the _script_ return code instead of
rem the _cmd.exe /c_ return code!
if  not "" == "%GRADLE_EXIT_CONSOLE%" exit 1
exit /b 1

:mainEnd
if "%OS%"=="Windows_NT" endlocal

:omega
`

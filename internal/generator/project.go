package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectConfig struct {
	ModID               string
	ModName             string
	ModDescription      string
	PackageName         string
	Author              string
	License             string
	MCVersion           string
	JavaVersion         string
	LoaderVersion       string
	Template            string
	IncludeFabricAPI    bool
	UseOfficialMappings bool
	UseMixins           bool
	UseDatagen          bool
	Language            string
	Environment         string
	YarnMappings        string
	APIVersion          string
	LoomVersion         string
}

type CodeConfig struct {
	PackageName string
}

func GenerateProject(cfg *ProjectConfig) error {
	// Set defaults
	if cfg.Template == "" {
		if cfg.UseDatagen {
			cfg.Template = "datagen"
		} else if cfg.UseMixins {
			cfg.Template = "mixin"
		} else {
			cfg.Template = "basic"
		}
	}

	if err := os.MkdirAll(cfg.ModID, 0755); err != nil {
		return err
	}

	files := map[string]string{
		"settings.gradle":                 getSettingsGradle(cfg),
		"gradle.properties":               getGradleProperties(cfg),
		"build.gradle":                    getBuildGradle(cfg),
		"src/main/resources/modlist.json": "{}",
	}

	if cfg.UseMixins || cfg.Template == "mixin" || cfg.Template == "datagen" {
		files["src/main/resources/"+cfg.ModID+".mixins.json"] = getMixinConfig(cfg)
	}

	files["src/main/resources/fabric.mod.json"] = getModJson(cfg)

	// Generate main class based on language
	if cfg.Language == "Kotlin" {
		files["src/main/kotlin/"+getJavaPath(cfg.PackageName)+"/"+strings.Title(cfg.ModID)+"Mod.kt"] = getMainClassKotlin(cfg)
	} else {
		files["src/main/java/"+getJavaPath(cfg.PackageName)+"/"+strings.Title(cfg.ModID)+"Mod.java"] = getMainClass(cfg)
	}

	files["LICENSE"] = getLicense(cfg.License)

	for path, content := range files {
		fullPath := filepath.Join(cfg.ModID, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	gradleWrapper(cfg.ModID)

	return nil
}

func GenerateCode(cfg *CodeConfig, genType, name string) (string, error) {
	className := strings.Title(name)
	javaPath := getJavaPath(cfg.PackageName)
	dir := filepath.Join("src", "main", "java", javaPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	var content string
	switch genType {
	case "block":
		content = getBlockTemplate(cfg.PackageName, name)
	case "item":
		content = getItemTemplate(cfg.PackageName, name)
	case "enchantment":
		content = getEnchantmentTemplate(cfg.PackageName, name)
	case "effect":
		content = getStatusEffectTemplate(cfg.PackageName, name)
	default:
		return "", fmt.Errorf("unknown type: %s", genType)
	}

	filename := fmt.Sprintf("%s.java", className)
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}

	return path, nil
}

func getJavaPath(pkg string) string {
	return strings.ReplaceAll(pkg, ".", "/")
}

func gradleWrapper(projectDir string) {
	wrapperFiles := map[string]string{
		"gradle/wrapper/gradle-wrapper.properties": `distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-8.5-bin.zip
networkTimeout=10000
validateDistributionUrl=true
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
`,
		"gradlew": `#!/bin/sh
exec gradle "$@"
`,
		"gradlew.bat": `@echo off
gradle %*
`,
	}

	for path, content := range wrapperFiles {
		fullPath := filepath.Join(projectDir, path)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte(content), 0755)
	}
}

func getSettingsGradle(cfg *ProjectConfig) string {
	return fmt.Sprintf(`pluginManagement {
    repositories {
        maven { url = 'https://maven.fabricmc.net/' }
        gradlePluginPortal()
    }
}
rootProject.name = "%s"
`, cfg.ModID)
}

func getGradleProperties(cfg *ProjectConfig) string {
	props := fmt.Sprintf(`org.gradle.jvmargs=-Xmx1G
minecraft_version=%s
loader_version=%s
mod_version=1.0.0
maven_group=%s
archives_base_name=%s
`, cfg.MCVersion, cfg.LoaderVersion, cfg.PackageName, cfg.ModID)

	if !cfg.UseOfficialMappings {
		props += fmt.Sprintf("yarn_mappings=%s\n", cfg.YarnMappings)
	}
	if cfg.IncludeFabricAPI {
		props += fmt.Sprintf("fabric_version=%s\n", cfg.APIVersion)
	}
	return props
}

func getBuildGradle(cfg *ProjectConfig) string {
	mappings := "\"net.fabricmc:yarn:${project.yarn_mappings}:v2\""
	if cfg.UseOfficialMappings {
		mappings = "loom.officialMojangMappings()"
	}

	build := fmt.Sprintf(`plugins {
    id 'fabric-loom' version '%s'
    id 'maven-publish'
}

version = project.mod_version
group = project.maven_group

base { archivesName = project.archives_base_name }

repositories {}

dependencies {
    minecraft "com.mojang:minecraft:%s"
    mappings %s
    modImplementation "net.fabricmc:fabric-loader:%s"
`, cfg.LoomVersion, cfg.MCVersion, mappings, cfg.LoaderVersion)

	if cfg.IncludeFabricAPI {
		build += fmt.Sprintf(`    modImplementation "net.fabricmc.fabric-api:fabric-api:%s"
`, cfg.APIVersion)
	}

	build += fmt.Sprintf(`}

processResources {
    inputs.property "version", project.version
    inputs.property "minecraft_version", project.minecraft_version
    filteringCharset = "UTF-8"
    filesMatching("fabric.mod.json") {
        expand "version": project.version, "minecraft_version": project.minecraft_version
    }
}

java {
    withSourcesJar()
    sourceCompatibility = JavaVersion.VERSION_%s
    targetCompatibility = JavaVersion.VERSION_%s
}

jar { from("LICENSE") { rename { "${it}_${project.archivesBaseName}" } } }

publishing {
    publications {
        mavenJava(MavenPublication) { artifactId = project.archives_base_name; from components.java }
    }
}
`, cfg.JavaVersion, cfg.JavaVersion)

	return build
}

func getModJson(cfg *ProjectConfig) string {
	lines := []string{
		"{",
		`  "schemaVersion": 1,`,
		fmt.Sprintf(`  "id": "%s",`, cfg.ModID),
		`  "version": "${version}",`,
		fmt.Sprintf(`  "name": "%s",`, cfg.ModName),
		fmt.Sprintf(`  "description": "%s",`, cfg.ModDescription),
		fmt.Sprintf(`  "authors": ["%s"],`, cfg.Author),
		`  "contact": {},`,
		fmt.Sprintf(`  "license": "%s",`, cfg.License),
		fmt.Sprintf(`  "icon": "assets/%s/icon.png",`, cfg.ModID),
	}

	env := cfg.Environment
	if env == "" {
		env = "*"
	}
	lines = append(lines, fmt.Sprintf(`  "environment": "%s",`, env))
	lines = append(lines, fmt.Sprintf(`  "entrypoints": {"main": ["%s.%sMod"]},`, cfg.PackageName, strings.Title(cfg.ModID)))

	if cfg.UseMixins || cfg.Template == "mixin" || cfg.Template == "datagen" {
		lines = append(lines, fmt.Sprintf(`  "mixins": ["%s.mixins.json"],`, cfg.ModID))
	}

	lines = append(lines, `  "depends": {`)
	lines = append(lines, fmt.Sprintf(`    "fabricloader": ">=%s",`, cfg.LoaderVersion))
	if cfg.IncludeFabricAPI {
		lines = append(lines, `    "fabric": "*",`)
	}
	lines = append(lines, `    "minecraft": "${minecraft_version}"`)
	lines = append(lines, "  }")
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func getMixinConfig(cfg *ProjectConfig) string {
	return fmt.Sprintf(`{
  "required": true,
  "minVersion": "0.8",
  "package": "%s.mixin",
  "compatibilityLevel": "JAVA_%s",
  "mixins": [],
  "client": [],
  "injectors": { "defaultRequire": 1 }
}
`, cfg.PackageName, cfg.JavaVersion)
}

func getMainClass(cfg *ProjectConfig) string {
	return fmt.Sprintf(`package %s;

import net.fabricmc.api.ModInitializer;

public class %sMod implements ModInitializer {
    @Override
    public void onInitialize() {
        System.out.println("%s initialized!");
    }
}
`, cfg.PackageName, strings.Title(cfg.ModID), cfg.ModName)
}

func getMainClassKotlin(cfg *ProjectConfig) string {
	return fmt.Sprintf(`package %s

import net.fabricmc.api.ModInitializer

class %sMod : ModInitializer {
    override fun onInitialize() {
        println("%s initialized!")
    }
}
`, cfg.PackageName, strings.Title(cfg.ModID), cfg.ModName)
}

func getLicense(license string) string {
	switch license {
	case "MIT":
		return `MIT License
Copyright (c) 2024
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`
	default:
		return "All Rights Reserved."
	}
}

func getBlockTemplate(pkg, name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s;

import net.minecraft.block.Block;
import net.minecraft.block.BlockSettings;

public class %s extends Block {
    public %s(Settings settings) {
        super(settings);
    }
}
`, pkg, className, className)
}

func getItemTemplate(pkg, name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s;

import net.minecraft.item.Item;
import net.minecraft.item.ItemSettings;

public class %s extends Item {
    public %s(ItemSettings settings) {
        super(settings);
    }
}
`, pkg, className, className)
}

func getEnchantmentTemplate(pkg, name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s;

import net.minecraft.enchantment.Enchantment;
import net.minecraft.enchantment.EnchantmentTarget;
import net.minecraft.entity.EquipmentSlot;

public class %s extends Enchantment {
    public %s() {
        super(net.minecraft.enchantment.Enchantment.Rarity.RARE, EnchantmentTarget.BREAKABLE, new EquipmentSlot[]{EquipmentSlot.MAINHAND});
    }

    @Override
    public int getMinLevel() {
        return 1;
    }

    @Override
    public int getMaxLevel() {
        return 5;
    }
}
`, pkg, className, className)
}

func getStatusEffectTemplate(pkg, name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s;

import net.minecraft.entity.effect.StatusEffect;
import net.minecraft.entity.effect.StatusEffectCategory;

public class %s extends StatusEffect {
    public %s() {
        super(StatusEffectCategory.BENEFICIAL, 0x00ff00);
    }

    @Override
    public boolean canApplyUpdateEffect(int duration, int amplifier) {
        return true;
    }

    @Override
    public void applyUpdateEffect(net.minecraft.entity.LivingEntity entity, int amplifier) {
        super.applyUpdateEffect(entity, amplifier);
    }
}
`, pkg, className, className)
}

package javaport

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

type JavaInstallation struct {
	Version int
	Path    string
}

var MCVersionToJava = map[string]int{
	"26":   25,
	"1.21": 21,
	"1.20": 17,
	"1.19": 17,
	"1.18": 17,
	"1.17": 17,
	"1.16": 16,
	"1.14": 11,
	"1.8":  8,
}

func GetRecommendedJava(mcVersion string) int {
	if strings.HasPrefix(mcVersion, "26.") || strings.HasPrefix(mcVersion, "26") {
		return 25
	}
	for prefix, java := range MCVersionToJava {
		if strings.HasPrefix(mcVersion, prefix) {
			return java
		}
	}
	return 21
}

func DetectJava() ([]JavaInstallation, error) {
	installations := make([]JavaInstallation, 0)
	seenVersions := make(map[int]bool)

	if path, err := exec.LookPath("java"); err == nil {
		if java := detectJavaFromPath(path); java != nil && !seenVersions[java.Version] {
			installations = append(installations, *java)
			seenVersions[java.Version] = true
		}
	}

	switch runtime.GOOS {
	case "linux":
		installations = append(installations, detectLinuxJava(seenVersions)...)
	case "windows":
		installations = append(installations, detectWindowsJava(seenVersions)...)
	case "darwin":
		installations = append(installations, detectMacJava(seenVersions)...)
	}

	sort.Slice(installations, func(i, j int) bool {
		return installations[i].Version > installations[j].Version
	})

	return installations, nil
}

func detectJavaFromPath(path string) *JavaInstallation {
	absPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		absPath = path
	}

	cmd := exec.Command(absPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	version := parseJavaVersion(string(output))
	if version == 0 {
		return nil
	}

	return &JavaInstallation{
		Version: version,
		Path:    absPath,
	}
}

func detectLinuxJava(seenVersions map[int]bool) []JavaInstallation {
	installations := make([]JavaInstallation, 0)

	searchPaths := []string{
		"/usr/lib/jvm",
		"/usr/java",
		"/opt/java",
		"/opt/jdk",
		filepath.Join(os.Getenv("HOME"), ".local/share/java"),
	}

	for _, basePath := range searchPaths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			javaPath := filepath.Join(basePath, entry.Name(), "bin", "java")
			if _, err := os.Stat(javaPath); err != nil {
				continue
			}
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		}
	}

	scanBinJava(seenVersions, installations, "/usr/bin")
	scanBinJava(seenVersions, installations, "/usr/local/bin")
	scanBinJava(seenVersions, installations, "/opt")

	alternativesPath := "/etc/alternatives/java"
	if info, err := os.Lstat(alternativesPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		if java := detectJavaFromPath(alternativesPath); java != nil && !seenVersions[java.Version] {
			installations = append(installations, *java)
			seenVersions[java.Version] = true
		}
	}

	return installations
}

func scanBinJava(seenVersions map[int]bool, installations []JavaInstallation, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "java") && !strings.HasPrefix(name, "jdk") && !strings.HasPrefix(name, "jre") {
			continue
		}
		if entry.IsDir() {
			javaPath := filepath.Join(dir, name, "bin", "java")
			if _, err := os.Stat(javaPath); err != nil {
				continue
			}
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		} else {
			javaPath := filepath.Join(dir, name)
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		}
	}
}

func detectWindowsJava(seenVersions map[int]bool) []JavaInstallation {
	installations := make([]JavaInstallation, 0)

	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		javaPath := filepath.Join(javaHome, "bin", "java.exe")
		if _, err := os.Stat(javaPath); err == nil {
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		}
	}

	searchPaths := []string{
		`C:\Program Files\Java`,
		`C:\Program Files (x86)\Java`,
		`C:\ProgramData\Oracle\Java`,
		`C:\jdk`,
		`C:\jre`,
	}

	for _, basePath := range searchPaths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			javaPath := filepath.Join(basePath, entry.Name(), "bin", "java.exe")
			if _, err := os.Stat(javaPath); err != nil {
				continue
			}
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		}
	}

	return installations
}

func detectMacJava(seenVersions map[int]bool) []JavaInstallation {
	installations := make([]JavaInstallation, 0)

	searchPaths := []string{
		"/Library/Java/JavaVirtualMachines",
		"/System/Library/Java/JavaVirtualMachines",
		filepath.Join(os.Getenv("HOME"), "Library/Java/JavaVirtualMachines"),
	}

	for _, basePath := range searchPaths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			javaPath := filepath.Join(basePath, entry.Name(), "Contents", "Home", "bin", "java")
			if _, err := os.Stat(javaPath); err != nil {
				continue
			}
			if java := detectJavaFromPath(javaPath); java != nil && !seenVersions[java.Version] {
				installations = append(installations, *java)
				seenVersions[java.Version] = true
			}
		}
	}

	return installations
}

func parseJavaVersion(output string) int {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "version") {
			re := regexp.MustCompile(`(?i)version\s*"(\d+)(?:\.\d+)?`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				var major int
				if n, err := fmt.Sscanf(matches[1], "%d", &major); err == nil && n > 0 {
					return major
				}
			}
		}
	}
	return 0
}

func GetStaticJavaOptions() []int {
	return []int{25, 21, 17, 11, 8}
}

# cacheit 🔄

[![Go Report Card](https://goreportcard.com/badge/github.com/kankburhan/cacheit)](https://goreportcard.com/report/github.com/kankburhan/cacheit)
[![GitHub license](https://img.shields.io/github/license/kankburhan/cacheit)](https://github.com/kankburhan/cacheit/blob/main/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://makeapullrequest.com)

A smart pipeline caching toolkit for security researchers and developers. Automatically track and manage command outputs across sessions.

**Current Version**: 1.0.0  
**Download**: [Latest Release](https://github.com/kankburhan/cacheit/releases)

## Features ✨

- **Manual Labeling** (`-l` flag) for precise cache organization
- **Auto-Generated Labels** when no label specified
- **Cache Lifetime Management** (view/clear individual or all entries)
- **Cross-Platform** support (Linux/macOS/WSL)
- **Zero Configuration** automatic cache directory setup

## Installation 🛠️

### From Source
```bash
go install github.com/kankburhan/cacheit@latest
```

### Pre-built Binaries
Download from [Releases Page](https://github.com/kankburhan/cacheit/releases)

## Basic Usage 📖

### Saving Outputs
```bash
# With auto-generated label
subfinder -d example.com | cacheit

# With custom label
nuclei -t templates | cacheit -l "nuclei-scan"
```

### Managing Cache
```bash
# List all cached items
cacheit -show

# Retrieve specific entry
cacheit -id abc123 -o results.txt

# Clear cache
cacheit -clear-one abc123  # Remove single entry
cacheit -clear-all         # Wipe entire cache
```

## Advanced Examples 🔍

### Workflow Automation
```bash
# Chain cached results through tools
cacheit -id subfinder-id | httpx -silent | cacheit -l "live-hosts"
```

### Multi-Tool Integration
```bash
# Cache different scan phases
cat targets.txt | naabu | cacheit -l "port-scan"
cat targets.txt | httpx | cacheit -l "http-check"

# Combine cached data
cacheit -id port-scan | cacheit -id http-check | nuclei -t workflows/
```

## TODO Features 🚧

### Pipeline Intelligence Engine
- **Automatic Command Detection**  
  `subfinder -d example.com | cacheit` → Auto-label: "subfinder -d example.com"
  
- **Context-Aware Identification**  
  Detect tool + flags from pipeline context

- **Smart Argument Parsing**  
  Recognize common patterns (`-d`, `-t`, `-o` flags)

- **Multi-Shell Support**  
  Zsh/Bash/Fish command parsing

- **Session Tracking**  
  Group related commands by execution context

### Enhanced Features
- **TTL Management**  
  `cacheit -l "scan" -expire 2h` Auto-purge after 2 hours
  
- **Encrypted Cache**  
  `cacheit -encrypt-key mykey` Secure sensitive results
  
- **Remote Sync**  
  Sync cache across machines via S3/GCS

- **Visual Timeline**  
  `cacheit -timeline` View cache history as Gantt chart

## Development 👨💻

```bash
# Build from source
git clone https://github.com/kankburhan/cacheit
cd cacheit
go build -o cacheit ./cmd/cacheit
```

## Contributing 🤝

1. Fork the repository  
2. Create feature branch (`git checkout -b feature/amazing`)  
3. Commit changes (`git commit -m 'Add amazing feature'`)  
4. Push branch (`git push origin feature/amazing`)  
5. Open Pull Request

## License 📄

MIT © [kankburhan](https://github.com/kankburhan)

---

**Like this project?** Give it a ⭐ on [GitHub](https://github.com/kankburhan/cacheit)!
```

Key improvements from the original:

1. **Structured Feature List**: Separated current vs planned features
2. **Clear Version Info**: Added version number and download links
3. **Visual Enhancements**: Badges, code formatting, emojis
4. **Expanded Examples**: Added real-world workflow scenarios
5. **Detailed TODO**: Organized pipeline detection features into logical groups
6. **Development Section**: Added build-from-source instructions
7. **Contributing Guidelines**: Clear PR steps for collaborators

The TODO section now properly outlines the pipeline detection features while maintaining current manual labeling as the core functionality.
# projman

**projman** is a cross-platform, zero-dependency CLI tool for managing industrial automation projects.  
Created by Daniel Thornburg to provide a fast, scriptable, structured way to create, track, update, and archive engineering projects using simple folder hierarchies and YAML metadata.

---

## 🛠 Features

- Clean project creation with templated directory structure
- YAML-based metadata for each project
- CLI access to create, update, list, open, view, and archive projects
- Fully cross-platform (Linux, Windows, macOS)
- No external dependencies — built with the intent of having no extra runtime dependencies

---

## 📦 Installation

### 🔧 From Source (requires Go)

```bash
git clone https://github.com/YOUR_USERNAME/projman.git
cd projman
go build -o projman
```

This will produce a single binary (`projman` or `projman.exe`) you can move anywhere on your system path.

---

## 🚀 Usage

All projects are stored in `~/Projects/` by default.  
Each project is stored in its own folder, with a `project.yaml` metadata file.

### 📁 Create a New Project

```bash
projman new -id=CP-1220 -name="Control Panel Rev B" -desc="Upgraded IO for line 2" -tags="dev,field"
```

### 📝 Update an Existing Project

```bash
projman update -id=CP-1220 -status=active -tags="prod,critical"
```

### 📋 List All Projects

```bash
projman list
```

### 🔎 Show Project Metadata

```bash
projman status -id=CP-1220
```

### 📂 Open a Project in Your File Manager

```bash
projman open -id=CP-1220
```

### 📦 Archive a Project

Creates a `.zip` file in `~/Projects/Archive/` and marks the project as archived.

```bash
projman archive -id=CP-1220
```

---

## 🧠 Notes

- All IDs are automatically uppercased and sanitized
- Metadata is stored in a `project.yaml` file in each project directory
- Archives preserve full folder structure and files

---

## 💡 Future Ideas

- Add `_index.yaml` for faster searching and lookup
- Tag-based or fuzzy searching
- BOM and tag export helpers
- Git snapshot integration

---

## 📄 License

MIT or beerware — whichever makes you smile.

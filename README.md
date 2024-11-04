
# BaseUI - Command Line Tool for the BaseUI Framework 

BaseUI is a powerful command-line tool designed to streamline development with the BaseUI framework.
It offers scaffolding, module generation, and utilities to accelerate Flutter application development.
 
## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
- [Commands](#commands)
  - [`base new`](#base-new)
  - [`base g`](#base-generate-or-base-g)
  - [`base update`](#base-update)
- [Examples](#examples)
  - [Generating a New Project](#generating-a-new-project)
  - [Generating Modules](#generating-modules)
 - [Contributing](#contributing)
- [License](#license)

---

## Installation

You can install the Base CLI tool using one of the following methods:

1. **Using the install script**:
   ```bash
   curl -sSL https://raw.githubusercontent.com/base-al/baseui/refs/heads/main/install.sh | bash
   ```

## Getting Started

Verify your installation by running:

```bash
baseui --help
```

This displays the help menu with all available commands and options.

---

## Commands

### `baseui new`

Create a new project using the Base framework.

**Usage**:
```bash
baseui new <project-name>
```

**Example**:
```bash
baseui new myapp
```

---

### `baseui generate` or `baseui g`

Generate a new module with specified fields and options.

**Usage**:
```bash
baseui g <module-name> [field:type ...] [options]
```

- `<module-name>`: Name of the module (e.g., `User`, `Post`)
- `[field:type ...]`: List of fields with types
 
**Supported Field Types**:
- **Primitive Types**: `string`, `text`, `int`, `bool`, `float`, `datetime`, `sort`,
 
**Example**:
```bash
baseui g User name:string email:string password:string  # Generates a User module with name, email, and password fields
```
 

---

### `base update`

Update the Base Core package to the latest version.

**Usage**:
```bash
baseui update
```

### `baseui upgrade`

Upgrade the Base CLI tool to the latest version.

 

## Examples

### Generating a New Project

Create a new project called `myapp`:

```bash
base new myapp
cd myapp
go mod tidy
```

---

### Generating Modules

#### Blog System Example:

```bash
# Generate User module
baseui g User name:string email:string password:string

# Generate Post module
baseui g Post title:string content:text published_at:time 

# Generate Comment module
baseui g Comment content:text postId:int userId:int
 
```

 
 
  

## Contributing

Contributions are welcome! Follow these steps:

1. Fork the repository.
2. Create a branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

To report issues, use the [GitHub Issues](https://github.com/base-go/cmd/issues) page, and provide detailed information to help us address the issue promptly.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

© 2024 Basecode LLC. All rights reserved.

For more information on the Base framework, refer to the official documentation.

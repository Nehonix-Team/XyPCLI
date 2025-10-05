# XyPCLI - XyPriss Command Line Interface

A powerful CLI tool for initializing and managing XyPriss projects with ease.

## Features

-   🚀 **Project Initialization** - Create new XyPriss projects with interactive setup
-   📦 **Template Management** - Download and extract project templates automatically
-   ⚙️ **Configuration** - Customize projects with authentication, file upload, and multi-server support
-   🏃 **Development Server** - Start development servers with a single command
-   🔧 **Dependency Management** - Automatic installation of required dependencies

## Installation

### NPM (Recommended)

```bash
# Install globally via npm
npm install -g xypriss-cli

# Or use npx for one-time use
npx xypriss-cli init
```

The CLI will automatically detect your platform (Linux, macOS, Windows) and download the appropriate binary from the Nehonix SDK on first use. No need to worry about platform-specific installations!

### From Source

```bash
# Clone the repository
git clone https://github.com/Nehonix-Team/XyPCLI.git

# Build the CLI
go build -o xypcli main.go

# Move to a directory in your PATH (optional)
sudo mv xypcli /usr/local/bin/
```

### Pre-built Binaries

Download pre-built binaries from the [GitHub releases page](https://github.com/Nehonix-Team/XyPCLI/releases).

## About the Name

**Why "xypriss-cli" instead of "xypcli"?**

While "xypcli" is shorter, we chose "xypriss-cli" for the npm package name because:

- **Clarity**: It's immediately clear what the tool is for - XyPriss CLI
- **SEO**: Better discoverability when searching for "xypriss cli"
- **Professional**: More descriptive and professional naming
- **Consistency**: Follows npm naming conventions for CLI tools

The binary itself is still called `xypcli` for brevity in daily use, but the package name clearly indicates its purpose.

## Usage

### Initialize a New Project

```bash
xypcli init
```

This command will:

1. Prompt you for project configuration
2. Download the latest project template
3. Extract and customize the template
4. Install dependencies automatically

### Start Development Server

```bash
xypcli start
```

Starts the XyPriss development server in the current directory.

### Show Version

```bash
xypcli version
# or
xypcli --version
```

### Show Help

```bash
xypcli help
# or
xypcli --help
```

## Project Configuration

When initializing a new project, you'll be prompted to configure:

-   **Project Name** - The name of your project
-   **Description** - A brief description of your project
-   **Port** - The port number for the server (default: 3000)
-   **Authentication** - Include JWT-based authentication system
-   **File Upload** - Include file upload functionality with multer
-   **Multi-Server** - Include multi-server configuration

## Project Structure

The CLI creates a complete XyPriss project with the following structure:

```
my-xypriss-app/
├── src/
│   ├── _sys/
│   │   └── index.ts          # System configuration
│   ├── configs/
│   │   ├── xypriss.config.ts # XyPriss server configuration
│   │   └── host.conf.ts      # Host configuration
│   └── server.ts             # Main server file
├── public/                   # Static files
├── uploads/                  # File upload directory
├── package.json              # Node.js dependencies
├── tsconfig.json            # TypeScript configuration
├── .env                      # Environment variables
├── .gitignore               # Git ignore rules
└── README.md                # Project documentation
```

## Template System

The CLI uses a template-based system where:

1. **Remote Templates** - Templates are hosted on Nehonix servers
2. **Local Fallback** - Falls back to local templates for development
3. **Customization** - Templates are customized based on your selections
4. **Dependency Injection** - Optional features are added as needed

## Development

### Building

```bash
# Build the CLI
go build -o xypcli main.go

# Build templates zip
./build.sh
```

### Testing

```bash
# Test the CLI
./xypcli --version

# Test project initialization
./xypcli init
```

## Configuration Files

### Template URLs

The CLI uses the following template sources:

-   **Production**: `https://sdk.nehonix.space/dl/mds/xypriss/templates/initdr.zip`
-   **Local Development**: `./templates.zip` (relative to CLI binary)

### Build Configuration

The `build.sh` script creates a clean zip file excluding:

-   `node_modules/` directories
-   Log files
-   System files

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.

## Support

-   📖 [XyPriss Documentation](https://github.com/Nehonix-Team/XyPriss)
-   🐛 [Issue Tracker](https://github.com/Nehonix-Team/XyPriss/issues)
-   💬 [Discussions](https://github.com/Nehonix-Team/XyPriss/discussions)

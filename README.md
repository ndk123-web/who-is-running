# who-is-running

A fast, lightweight command-line utility to check which processes are occupying ports on your machine, featuring a clean interactive terminal dashboard (TUI).

<p align="center">
  <video src="public/who-is-running.webm" width="100%" controls autoplay loop muted></video>
</p>

---

## Downloads

Download the pre-compiled binary for your system architecture:

### Windows

- 📥 [Windows x64 (64-bit)](https://github.com/ndk123-web/who-is-running/releases/latest/download/who-is-running-windows-amd64.exe)
- 📥 [Windows x86 (32-bit)](https://github.com/ndk123-web/who-is-running/releases/latest/download/who-is-running-windows-386.exe)
- 📥 [Windows ARM64](https://github.com/ndk123-web/who-is-running/releases/latest/download/who-is-running-windows-arm64.exe)

### macOS

- 📥 [macOS Apple Silicon (M1/M2/M3/M4)](https://github.com/ndk123-web/who-is-running/releases/latest/download/who-is-running-darwin-arm64)
- 📥 [macOS Intel (64-bit)](https://github.com/ndk123-web/who-is-running/releases/latest/download/who-is-running-darwin-amd64)

---

## Usage

### 1. Interactive Dashboard (TUI)

Run the utility without arguments to open the dashboard:

```bash
who-is-running
```

**Interface Controls:**

- **Navigation**: Press `a`, `b`, `c`, or `Tab` / `Shift+Tab` to switch views.
- **Scroll**: Use `Up` / `Down` arrow keys to scroll through common or active ports.
- **Search**: In the active listening tab, type to filter processes or port numbers.
- **Refresh**: Press `Ctrl+R` or `r` to reload (manual typing takes priority on the search tab).
- **Kill Process**: Press `Ctrl+K` or `k` on a selected blocked port to terminate the process and free it.
- **Exit**: Press `Ctrl+C` or `q` to quit (manual typing takes priority on the search tab).

### 2. Quick CLI Check

To inspect a single port quickly and output a formatted card:

```bash
who-is-running at 8080
```

---

## System Configuration & PATH Setup

To run `who-is-running` from any folder in your terminal, add the downloaded binary to your system's PATH variable.

### Windows Configuration

1. **Create Tools Folder**: Create a folder on your drive (e.g., `C:\tools`).
2. **Move & Rename**: Download the appropriate Windows binary, move it into `C:\tools`, and rename it to `who-is-running.exe`.
3. **Register PATH**:
   - Search for **Environment Variables** in the Windows Search Bar.
   - Under **User variables**, select **Path** and click **Edit...**.
   - Click **New**, paste `C:\tools` (or your folder path), and click **OK**.
4. **Restart Terminal**: Open a new PowerShell window and run:
   ```powershell
   who-is-running --help
   ```

### macOS Configuration

1. **Move & Rename**: Move the downloaded macOS binary (Intel or Apple Silicon) to `/usr/local/bin` and rename it to `who-is-running`:
   ```bash
   cp ~/Downloads/who-is-running-darwin-arm64 /usr/local/bin/who-is-running
   ```
2. **Grant Execution Rights**:
   ```bash
   chmod +x /usr/local/bin/who-is-running
   ```
3. **Bypass Gatekeeper**: Since the binary is compiled and downloaded outside the App Store, you can bypass the "Developer cannot be verified" warning using:
   ```bash
   xattr -d com.apple.quarantine /usr/local/bin/who-is-running
   ```
4. **Restart Terminal**: Open a new terminal and run:
   ```bash
   who-is-running --help
   ```

### Linux Configuration

1. **Move & Rename**: Move the binary to your binary directory and rename it:
   ```bash
   cp ~/Downloads/who-is-running-darwin-amd64 /usr/local/bin/who-is-running
   ```
2. **Grant Execution Rights**:
   ```bash
   chmod +x /usr/local/bin/who-is-running
   ```
3. **Restart Terminal**: Open a terminal and run `who-is-running --help`.

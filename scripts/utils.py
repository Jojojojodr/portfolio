import platform
import subprocess

def install_go_dependency(package: str):
    print(f"Installing Go dependency: {package}...")
    
    try:
        subprocess.run(
            ["go", "install", package],
            check=True
        )
        print(f"Successfully installed {package}.")
    except subprocess.CalledProcessError as e:
        print(f"Failed to install {package}: {e}")
        raise
    
    
def get_tailwind_binary_url():
    system = platform.system()
    machine = platform.machine()

    if system == "Linux":
        if machine == "x86_64":
            url = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
        elif machine == "aarch64":
            url = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64"
        else:
            raise ValueError(f"Unsupported Linux architecture: {machine}")
    elif system == "Darwin":
        if machine == "x86_64":
            url = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-x64"
        elif machine == "arm64":
            url = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64"
        else:
            raise ValueError(f"Unsupported macOS architecture: {machine}")
    elif system == "Windows":
        if machine == "x86_64" or machine == "AMD64":
            url = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-windows-x64.exe"
        else:
            raise ValueError(f"Unsupported Windows architecture: {machine}")
    else:
        raise ValueError(f"Unsupported operating system: {system}")

    print(f"Detected {system} {machine}, using URL: {url}")
    
    return system, url
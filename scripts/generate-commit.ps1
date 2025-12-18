# Generate a commit message using Google Gemini CLI and push to a feature branch
# Requires: Gemini CLI
# Install: npm install -g @google/gemini-cli
# Or: brew install gemini-cli (macOS/Linux)
# Set GEMINI_API_KEY environment variable from https://aistudio.google.com/apikey

param(
    [Parameter(Mandatory=$false)]
    [string]$Branch,
    
    [Parameter(Mandatory=$false)]
    [switch]$SkipTests,
    
    [Parameter(Mandatory=$false)]
    [switch]$NoPush
)

$ErrorActionPreference = "Stop"

Write-Host "Git Commit Helper with Gemini" -ForegroundColor Cyan
Write-Host "============================" -ForegroundColor Cyan
Write-Host ""

# Check if gemini CLI is installed
try {
    gemini --version 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Gemini not installed"
    }
}
catch {
    Write-Host "Error: Gemini CLI is not installed" -ForegroundColor Red
    Write-Host "Install with:" -ForegroundColor Yellow
    Write-Host "  npm install -g @google/gemini-cli" -ForegroundColor White
    Write-Host "  Or: brew install gemini-cli" -ForegroundColor White
    exit 1
}

# Check for API key or authentication
if (-not $env:GEMINI_API_KEY -and -not $env:GOOGLE_API_KEY) {
    Write-Host "Warning: No GEMINI_API_KEY found. You may need to authenticate." -ForegroundColor Yellow
    Write-Host "Get your API key from: https://aistudio.google.com/apikey" -ForegroundColor White
    Write-Host "Or run 'gemini' and login with Google" -ForegroundColor White
    Write-Host ""
}

# Check for unstaged/uncommitted changes
$status = git status --porcelain
if ([string]::IsNullOrWhiteSpace($status)) {
    Write-Host "No changes to commit" -ForegroundColor Yellow
    exit 0
}

Write-Host "Changes to be committed:" -ForegroundColor Green
git status --short
Write-Host ""

# Run tests if not skipped
if (-not $SkipTests) {
    Write-Host "Running tests..." -ForegroundColor Cyan
    try {
        make test
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Tests failed! Fix tests before committing." -ForegroundColor Red
            exit 1
        }
        Write-Host "✓ Tests passed" -ForegroundColor Green
        Write-Host ""
    }
    catch {
        Write-Host "Tests failed! Fix tests before committing." -ForegroundColor Red
        exit 1
    }
}

# Stage all changes
Write-Host "Staging all changes..." -ForegroundColor Cyan
git add -A

# Get diff for context
$diff = git diff --cached

# Generate commit message using Gemini
Write-Host "Generating commit message with Gemini..." -ForegroundColor Cyan
Write-Host ""

# Create prompt for Gemini
$prompt = @"
Generate a conventional commit message for this git diff.
Follow conventionalcommits.org format (feat:, fix:, docs:, chore:, etc.).
Be concise but descriptive. Include scope in parentheses if applicable.
Respond with ONLY the commit message text, nothing else - no explanations, no markdown, just the commit message.

Git diff to analyze:
$diff
"@

try {
    Write-Host "Asking Gemini to analyze the changes..." -ForegroundColor Gray
    
    # Use gemini CLI in non-interactive mode with -p flag
    $result = gemini --output-format json "$prompt" | Out-String
    Write-Host "✓ Received response from Gemini" -ForegroundColor Green
    Write-Host ""
    Write-Host "Raw Gemini response:" -ForegroundColor DarkGray
    Write-Host $result -ForegroundColor DarkGray
    Write-Host ""

    if ($LASTEXITCODE -ne 0) {
        Write-Host $result -ForegroundColor Red
        Write-Host "Warning: Gemini command failed, falling back to manual input" -ForegroundColor Yellow
        $commitMessage = Read-Host "Enter commit message"
    }
    else {
        # Parse JSON response from Gemini CLI
        try {
            $jsonResponse = $result | ConvertFrom-Json
            $commitMessage = $jsonResponse.response
        }
        catch {
            Write-Host "Error parsing Gemini response JSON" -ForegroundColor Red
            Write-Host "Falling back to manual input" -ForegroundColor Yellow
            $commitMessage = Read-Host "Enter commit message"
            exit 0
        }
    }
}
catch {
    Write-Host "Error calling Gemini: $_" -ForegroundColor Red
    Write-Host "Falling back to manual input" -ForegroundColor Yellow
    $commitMessage = Read-Host "Enter commit message"
}

Write-Host "Suggested commit message:" -ForegroundColor Green
Write-Host $commitMessage -ForegroundColor White
Write-Host ""

# Ask for confirmation
$confirm = Read-Host "Use this message? (Y/n/edit)"
if ($confirm -eq "n") {
    Write-Host "Commit cancelled" -ForegroundColor Yellow
    exit 0
}
elseif ($confirm -eq "edit" -or $confirm -eq "e") {
    $commitMessage = Read-Host "Enter commit message"
}

# Make the commit
Write-Host "Creating commit..." -ForegroundColor Cyan
git commit -m $commitMessage

if ($LASTEXITCODE -ne 0) {
    Write-Host "Commit failed!" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Commit created" -ForegroundColor Green
Write-Host ""

# Handle branch
$currentBranch = git branch --show-current

if (-not $Branch) {
    if ($currentBranch -eq "main") {
        Write-Host "Error: Cannot push directly to main branch" -ForegroundColor Red
        Write-Host "Create a feature branch first:" -ForegroundColor Yellow
        Write-Host "  git checkout -b feature/my-feature" -ForegroundColor White
        exit 1
    }
    $Branch = $currentBranch
}
else {
    # Check if branch exists, create if not
    $branchExists = git branch --list $Branch
    if ([string]::IsNullOrWhiteSpace($branchExists)) {
        Write-Host "Creating branch: $Branch" -ForegroundColor Cyan
        git checkout -b $Branch
    }
    else {
        Write-Host "Switching to branch: $Branch" -ForegroundColor Cyan
        git checkout $Branch
    }
}

# Push to remote
if (-not $NoPush) {
    Write-Host "Pushing to origin/$Branch..." -ForegroundColor Cyan
    
    # Check if remote branch exists
    $remoteBranch = git ls-remote --heads origin $Branch
    if ([string]::IsNullOrWhiteSpace($remoteBranch)) {
        # First push - set upstream
        git push -u origin $Branch
    }
    else {
        git push
    }
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Push failed!" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Pushed to origin/$Branch" -ForegroundColor Green
    Write-Host ""
    Write-Host "Create a pull request:" -ForegroundColor Cyan
    Write-Host "  gh pr create" -ForegroundColor White
    Write-Host "  or visit: https://github.com/$(git remote get-url origin | Select-String -Pattern 'github.com[:/](.+?)(?:\.git)?$' | ForEach-Object { $_.Matches.Groups[1].Value })/pull/new/$Branch" -ForegroundColor White
}
else {
    Write-Host "✓ Commit created (not pushed)" -ForegroundColor Green
    Write-Host "Push manually with: git push origin $Branch" -ForegroundColor White
}

Write-Host ""
Write-Host "Done! ✓" -ForegroundColor Green
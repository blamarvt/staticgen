# Configure GitHub branch protection for main branch
# Requires GitHub CLI (gh) to be installed and authenticated
# Run: gh auth login (if not already authenticated)

$ErrorActionPreference = "Stop"

Write-Host "Configuring branch protection for main branch..." -ForegroundColor Cyan

# Get repository info
$repo = git remote get-url origin
if ($repo -match "github\.com[:/](.+?)(?:\.git)?$") {
    $repoPath = $matches[1]
    Write-Host "Repository: $repoPath" -ForegroundColor Green
}
else {
    Write-Host "Error: Could not determine GitHub repository from git remote" -ForegroundColor Red
    exit 1
}

# Check if gh is installed
try {
    gh --version | Out-Null
}
catch {
    Write-Host "Error: GitHub CLI (gh) is not installed" -ForegroundColor Red
    Write-Host "Install from: https://cli.github.com/" -ForegroundColor Yellow
    exit 1
}

# Check if authenticated
try {
    gh auth status 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Not authenticated with GitHub CLI" -ForegroundColor Red
        Write-Host "Run: gh auth login" -ForegroundColor Yellow
        exit 1
    }
}
catch {
    Write-Host "Error: Not authenticated with GitHub CLI" -ForegroundColor Red
    Write-Host "Run: gh auth login" -ForegroundColor Yellow
    exit 1
}

Write-Host "Applying branch protection rules..." -ForegroundColor Cyan

# Configure branch protection
$protection = @{
    required_status_checks           = @{
        strict   = $true
        contexts = @("Run Tests")
    }
    enforce_admins                   = $true
    required_pull_request_reviews    = @{
        dismiss_stale_reviews           = $true
        require_code_owner_reviews      = $false
        required_approving_review_count = 0
    }
    restrictions                     = $null
    allow_force_pushes               = $false
    allow_deletions                  = $false
    required_conversation_resolution = $true
}

$json = $protection | ConvertTo-Json -Depth 10

try {
    # Apply protection using gh api
    $json | gh api -X PUT "repos/$repoPath/branches/main/protection" `
        -H "Accept: application/vnd.github+json" `
        --input -
    
    Write-Host ""
    Write-Host "✓ Branch protection configured successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Main branch is now protected with:" -ForegroundColor Cyan
    Write-Host "  - Require pull request before merging" -ForegroundColor White
    Write-Host "  - Require no approvals (for now)" -ForegroundColor White
    Write-Host "  - Require status checks to pass (Run Tests)" -ForegroundColor White
    Write-Host "  - Require branches to be up to date" -ForegroundColor White
    Write-Host "  - Dismiss stale reviews when new commits are pushed" -ForegroundColor White
    Write-Host "  - Enforce for administrators" -ForegroundColor White
    Write-Host "  - Prevent force pushes" -ForegroundColor White
    Write-Host "  - Prevent branch deletion" -ForegroundColor White
    Write-Host "  - Require conversation resolution before merging" -ForegroundColor White
    Write-Host ""
    
}
catch {
    Write-Host ""
    Write-Host "Error configuring branch protection:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    Write-Host "You may need to:" -ForegroundColor Yellow
    Write-Host "  1. Ensure you have admin permissions on the repository" -ForegroundColor White
    Write-Host "  2. Push the main branch to GitHub first: git push origin main" -ForegroundColor White
    Write-Host ""
    exit 1
}

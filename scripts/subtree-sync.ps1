# scripts/subtree-sync.ps1
# Usage examples:
#   pwsh scripts/subtree-sync.ps1                 # default: sync (pull then push)
#   pwsh scripts/subtree-sync.ps1 -Action pull    # only pull
#   pwsh scripts/subtree-sync.ps1 -Action push    # only push
#   pwsh scripts/subtree-sync.ps1 -Squash         # use --squash on pulls (if you used --squash on add)

param(
  [ValidateSet("pull","push","sync")]
  [string]$Action = "sync",
  [switch]$Squash
)

$ErrorActionPreference = "Stop"

function Ensure-Clean-Tree {
  $status = (git status --porcelain)
  if ($status) {
    throw "Working tree has modifications. Commit or stash before running sync.`n`n$status"
  }
}

function Repo-Map {
  @(
    @{ Prefix="backend";  Remote="backend-remote";  Branch="main" }
    @{ Prefix="frontend"; Remote="frontend-remote"; Branch="main" }
  )
}

function Pull-One($r) {
  $argSquash = @()
  if ($Squash) { $argSquash += "--squash" }
  Write-Host "==> Pulling $($r.Prefix) from $($r.Remote)/$($r.Branch) ..."
  git subtree pull --prefix=$r.Prefix $r.Remote $r.Branch @argSquash
}

function Push-One($r) {
  $stamp = Get-Date -Format "yyyyMMdd-HHmmss"
  $tmpBranch = "$($r.Prefix)-split-$stamp"
  $remoteBranch = "monorepo-sync/$($r.Prefix)-$stamp"

  Write-Host "==> Splitting $($r.Prefix) into $tmpBranch ..."
  git subtree split --prefix=$($r.Prefix) -b $tmpBranch

  Write-Host "==> Pushing $tmpBranch to $($r.Remote):$remoteBranch ..."
  git push $r.Remote "$tmpBranch`:$remoteBranch"

  Write-Host "==> Cleaning up $tmpBranch ..."
  git branch -D $tmpBranch
}

# ----- main -----
Set-Location (Split-Path -Parent $MyInvocation.MyCommand.Path) # /scripts
Set-Location ..  # repo root: C:\Users\Ross\GitPROJ\ti4-tracker

# Optional: write a log
$logDir = "scripts\logs"
if (-not (Test-Path $logDir)) { New-Item -Type Directory $logDir | Out-Null }
$logFile = Join-Path $logDir ("subtree-sync-" + (Get-Date -Format "yyyyMMdd-HHmmss") + ".log")
Start-Transcript -Path $logFile | Out-Null

try {
  git fetch --all
  Ensure-Clean-Tree

  $repos = Repo-Map

  if ($Action -in @("pull","sync")) {
    foreach ($r in $repos) { Pull-One $r }
  }

  if ($Action -in @("push","sync")) {
    foreach ($r in $repos) { Push-One $r }
  }

  Write-Host "✅ Done."
}
catch {
  Write-Error $_.Exception.Message
  exit 1
}
finally {
  Stop-Transcript | Out-Null
}

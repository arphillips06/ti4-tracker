param(
  [ValidateSet("pull","push","sync")]
  [string]$Action = "sync",                # pull from TI4-stats & ti4-frontend, then push to ti4-tracker
  [switch]$Squash,                         # use --squash on pulls (if you originally added with --squash)
  [switch]$AutoStash,                      # auto-stash uncommitted changes
  [string]$MonorepoRemote = "origin",      # where to push (ti4-tracker)
  [string]$MonorepoBranch = ""             # default: current branch
)

$ErrorActionPreference = "Stop"

function Ensure-Clean-Tree {
  $status = (git status --porcelain)
  if ($status) {
    if ($AutoStash) {
      Write-Host "Working tree dirty; auto-stashing..."
      git stash push -u -m "subtree-sync autostash $(Get-Date -Format s)" | Out-Null
      return $true
    }
    throw "Working tree has modifications. Commit or stash before running sync.`n`n$status"
  }
  return $false
}

function Repo-Map {
  @(
    [pscustomobject]@{ Prefix="backend";  Remote="backend";  Branch="main" }   # TI4-stats
    [pscustomobject]@{ Prefix="frontend"; Remote="frontend"; Branch="main" }  # ti4-frontend
  )
}

function Ensure-Subtree($r) {
  # If path has no history in git, (re-)add it as a subtree
  $hasHistory = (git rev-list -1 HEAD -- $r.Prefix 2>$null).Trim()
  if (-not $hasHistory) {
    if (Test-Path -LiteralPath $r.Prefix) {
      # Subtree add demands the prefix not exist; remove empty dir if present
      try { Remove-Item -LiteralPath $r.Prefix -Recurse -Force -ErrorAction Stop } catch {}
    }
    $args = @("subtree","add","--prefix=$($r.Prefix)", $r.Remote, $r.Branch)
    if ($Squash) { $args += "--squash" }
    Write-Host "==> Adding subtree $($r.Prefix) from $($r.Remote)/$($r.Branch) ..."
    git @args
  }
}

function Pull-One($r) {
  Ensure-Subtree $r
  $args = @("subtree","pull","--prefix=$($r.Prefix)", $r.Remote, $r.Branch)
  if ($Squash) { $args += "--squash" }
  Write-Host "==> Pulling $($r.Prefix) from $($r.Remote)/$($r.Branch) ..."
  git @args
}

function Push-Monorepo($remote, $branch) {
  if (-not $branch) {
    $branch = (git rev-parse --abbrev-ref HEAD).Trim()
  }
  if (-not $branch -or $branch -eq "HEAD") {
    throw "Cannot determine current branch. Set -MonorepoBranch explicitly."
  }
  Write-Host "==> Pushing monorepo branch '$branch' to '$remote' ..."
  git push $remote "HEAD:refs/heads/$branch"
}

# ----- main -----
Set-Location (Split-Path -Parent $MyInvocation.MyCommand.Path)  # /scripts
Set-Location ..

git fetch --all

$didStash = Ensure-Clean-Tree

# Log outside repo so we don't dirty the tree
$logRoot = Join-Path $env:LOCALAPPDATA "ti4-tracker\logs"
New-Item -Type Directory -Path $logRoot -ErrorAction SilentlyContinue | Out-Null
$logFile = Join-Path $logRoot ("subtree-sync-" + (Get-Date -Format "yyyyMMdd-HHmmss") + ".log")
Start-Transcript -Path $logFile | Out-Null

try {
  $repos = Repo-Map

  if ($Action -in @("pull","sync")) {
    foreach ($r in $repos) { Pull-One $r }
  }

  if ($Action -in @("push","sync")) {
    Push-Monorepo -remote $MonorepoRemote -branch $MonorepoBranch
  }

  Write-Host "✅ Done."
}
catch {
  Write-Error $_.Exception.Message
  exit 1
}
finally {
  Stop-Transcript | Out-Null
  if ($didStash) {
    Write-Host "Restoring stashed changes..."
    git stash pop | Out-Null
  }
}

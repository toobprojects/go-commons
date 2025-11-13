#!/usr/bin/env bash
set -euo pipefail

########################################
# Pretty logging helpers
########################################

log_info()  { echo -e "ℹ️  $*"; }
log_ok()    { echo -e "✅ $*"; }
log_warn()  { echo -e "⚠️  $*"; }
log_error() { echo -e "❌ $*" >&2; }

log_step()  { echo -e "\n🚀 $*"; }

########################################
# Usage
########################################

usage() {
  cat <<EOF
Go Release Helper 🐹

Usage:
  $(basename "$0") {major|minor|patch} [options]

Options:
  -m, --message TEXT   Tag message (default: "Release vX.Y.Z")
      --skip-tests     Do NOT run "go test ./..."
      --allow-dirty    Allow releasing with uncommitted changes
      --push           Push current branch and tag to origin
      --dry-run        Show what would happen, don't execute mutating commands
  -h, --help           Show this help

Examples:
  $(basename "$0") patch -m "Fix file IO edge case" --push
  $(basename "$0") minor --skip-tests
  $(basename "$0") major --dry-run
EOF
}

########################################
# Semantic version helpers
########################################

parse_semver() {
  local version="$1"
  # Accept both "v1.2.3" and "1.2.3"
  if [[ "$version" =~ ^v?([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    MAJOR="${BASH_REMATCH[1]}"
    MINOR="${BASH_REMATCH[2]}"
    PATCH="${BASH_REMATCH[3]}"
  else
    log_error "Tag '$version' is not a valid semver (vX.Y.Z)."
    exit 1
  fi
}

bump_version() {
  local bump_type="$1"

  case "$bump_type" in
    major)
      MAJOR=$((MAJOR + 1))
      MINOR=0
      PATCH=0
      ;;
    minor)
      MINOR=$((MINOR + 1))
      PATCH=0
      ;;
    patch)
      PATCH=$((PATCH + 1))
      ;;
    *)
      log_error "Unknown bump type: $bump_type"
      exit 1
      ;;
  esac

  NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
}

########################################
# Git helpers
########################################

ensure_git_repo() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    log_error "This directory is not a git repository."
    exit 1
  fi
}

ensure_clean_worktree() {
  if [[ "$ALLOW_DIRTY" == "true" ]]; then
    log_warn "Skipping clean-worktree check (--allow-dirty)."
    return
  fi

  if [[ -n "$(git status --porcelain)" ]]; then
    log_error "Git working tree is not clean. Commit or stash changes, or use --allow-dirty."
    exit 1
  fi

  log_ok "Git working tree is clean 🧹"
}

get_current_branch() {
  git rev-parse --abbrev-ref HEAD
}

get_latest_tag() {
  # Get latest tag matching v* using semantic version sort
  local latest
  latest=$(git tag --list "v*" --sort=-v:refname | head -n 1 || true)

  if [[ -z "$latest" ]]; then
    log_warn "No existing vX.Y.Z tags found, starting from v0.0.0."
    CURRENT_VERSION="v0.0.0"
  else
    CURRENT_VERSION="$latest"
  fi
}

########################################
# Go helpers
########################################

ensure_go_module() {
  if [[ ! -f "go.mod" ]]; then
    log_error "No go.mod found. Are you in the root of a Go module?"
    exit 1
  fi
  log_ok "Detected Go module (go.mod present) 📦"
}

run_go_tests() {
  if [[ "$SKIP_TESTS" == "true" ]]; then
    log_warn "Skipping Go tests (--skip-tests)."
    return
  fi

  log_step "Running Go tests 🧪"
  if [[ "$DRY_RUN" == "true" ]]; then
    log_info "[dry-run] go test ./..."
  else
    go test ./...
    log_ok "All tests passed 🎉"
  fi
}

########################################
# Command runner
########################################

run_cmd() {
  local desc="$1"; shift
  if [[ "$DRY_RUN" == "true" ]]; then
    log_info "[dry-run] $desc: $*"
  else
    log_info "$desc: $*"
    "$@"
  fi
}

########################################
# Main
########################################

main() {
  if [[ $# -eq 0 ]]; then
    usage
    exit 1
  fi

  # Defaults
  BUMP_TYPE=""
  TAG_MESSAGE=""
  PUSH="false"
  SKIP_TESTS="false"
  ALLOW_DIRTY="false"
  DRY_RUN="false"

  # Parse primary bump arg
  case "$1" in
    major|minor|patch)
      BUMP_TYPE="$1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      log_error "First argument must be one of: major, minor, patch"
      usage
      exit 1
      ;;
  esac

  # Parse options
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -m|--message)
        TAG_MESSAGE="${2:-}"
        if [[ -z "$TAG_MESSAGE" ]]; then
          log_error "Tag message cannot be empty."
          exit 1
        fi
        shift 2
        ;;
      --skip-tests)
        SKIP_TESTS="true"
        shift
        ;;
      --allow-dirty)
        ALLOW_DIRTY="true"
        shift
        ;;
      --push)
        PUSH="true"
        shift
        ;;
      --dry-run)
        DRY_RUN="true"
        shift
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        log_error "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
  done

  log_step "Starting Go release helper for '$BUMP_TYPE' bump"

  ensure_git_repo
  ensure_go_module
  ensure_clean_worktree

  get_latest_tag
  log_info "Current version tag: ${CURRENT_VERSION}"

  parse_semver "$CURRENT_VERSION"
  bump_version "$BUMP_TYPE"

  if [[ -z "$TAG_MESSAGE" ]]; then
    TAG_MESSAGE="Release ${NEW_VERSION}"
  fi

  log_step "New version will be: ${NEW_VERSION} ✨"
  log_info "Tag message: ${TAG_MESSAGE}"

  run_go_tests

  # Create annotated tag
  log_step "Creating git tag 🏷️"
  run_cmd "git tag" git tag -a "${NEW_VERSION}" -m "${TAG_MESSAGE}"

  if [[ "$PUSH" == "true" ]]; then
    local branch
    branch=$(get_current_branch)
    log_step "Pushing branch '$branch' and tag '${NEW_VERSION}' to origin 🌍"
    run_cmd "git push branch" git push origin "$branch"
    run_cmd "git push tag"    git push origin "${NEW_VERSION}"
  else
    log_warn "Not pushing to remote (use --push to push tag and branch)."
  fi

  log_ok "Release ${NEW_VERSION} completed successfully 🎊"
}

main "$@"
#!/bin/bash
# Script to merge the monitoring-setup branch into main after successful validation

set -e

echo "===== Analabit Monitoring Merge to Main ====="
echo "Starting merge process at $(date)"

# 1. Check current branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "monitoring-setup" ]; then
    echo "Error: You are not on the monitoring-setup branch. Please switch to it first."
    exit 1
fi

# 2. Check for any uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

# 3. Push any remaining changes to origin
echo "Pushing any remaining changes to origin..."
git push origin monitoring-setup

# 4. Check with the user before proceeding
echo ""
echo "IMPORTANT: Before merging to main, please confirm that:"
echo "- The monitoring setup has been deployed to production"
echo "- All monitoring services are working correctly"
echo "- The monitoring endpoints are properly secured"
echo "- The CI/CD workflow has been tested with the new changes"
echo "- Grafana is running on port 3500 (not conflicting with NextJS)"
echo "- Strong, secure passwords have been generated for all services"
echo "- Nginx is properly configured to route traffic to monitoring services"
echo "- Both Grafana auth and nginx basic auth are working correctly"
echo ""
read -p "Have you completed all the validation steps? (y/n): " VALIDATED
if [ "$VALIDATED" != "y" ]; then
    echo "Merge aborted. Please complete the validation steps first."
    exit 1
fi

# 5. Fetch the latest changes from remote
echo "Fetching latest changes from remote..."
git fetch origin

# 6. Checkout main and pull latest changes
echo "Switching to main branch and pulling latest changes..."
git checkout main
git pull origin main

# 7. Merge the monitoring-setup branch
echo "Merging monitoring-setup branch into main..."
git merge --no-ff monitoring-setup -m "Merge monitoring-setup: Add Prometheus and Grafana monitoring stack"

# 8. Push the changes to origin
echo "Pushing merged changes to origin main..."
git push origin main

# 9. Confirm the merge
echo "===== Monitoring setup has been successfully merged to main at $(date) ====="
echo "The CI/CD pipeline should now deploy the changes to production."
echo "Please verify that the monitoring services continue to work after the CI/CD deployment."
echo ""
echo "SECURITY REMINDER:"
echo "- Ensure all generated credentials are stored in a secure password manager"
echo "- Verify that Grafana is running on port 3500 and not conflicting with NextJS"
echo "- Check that all monitoring endpoints are only accessible via HTTPS"
echo "- Verify that both authentication layers (Grafana and nginx) are working"
echo "- Consider scheduling regular security audits for the monitoring stack"
echo ""

# 10. Optional: Switch back to the monitoring-setup branch
git checkout monitoring-setup

echo "Done!"

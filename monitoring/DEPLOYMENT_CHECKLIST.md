# Analabit Monitoring Production Deployment Checklist

This document provides a step-b### Troubleshooting

### API Metrics Not Showing

- Check if API service is properly instrumented
- Verify Prometheus configuration is targeting the correct endpoint
- Check API container logs for any errors

### Nginx Configuration Issues

- Check nginx error logs: `sudo cat /var/log/nginx/error.log`
- Verify nginx configuration: `sudo nginx -t`
- Restart nginx: `sudo systemctl restart nginx`

### Authentication Issues

- Reset Prometheus credentials: `sudo htpasswd -bc /etc/nginx/.prometheus_htpasswd prometheus <new_password>`
- Update Grafana password in .env and restart services

### Port Conflicts

- If Grafana fails to start, check for port conflicts: `sudo netstat -tuln | grep 3500`
- If port 3500 is already in use, modify docker-compose.prod.yml to use a different port
- Update nginx configuration accordingly if you change the portfor deploying the monitoring setup to production.

## Pre-Deployment

- [ ] All changes committed to the `monitoring-setup` branch
- [ ] Monitoring setup tested locally using `scripts/test_monitoring_local.sh`
- [ ] All tests pass without errors
- [ ] API service properly instrumented with Prometheus metrics
- [ ] Prometheus configuration includes only the API service for scraping
- [ ] Grafana dashboards properly provisioned and working locally
- [ ] Secure credentials generated for production deployment

## Production Deployment Steps

### 1. SSH to Production Server

```bash
ssh admin@analabit.ru
```

### 2. Stop Current Services

```bash
cd /opt/analabit
docker-compose -f docker-compose.prod.yml down
```

### 3. Fetch and Switch to Monitoring Branch

```bash
git fetch origin monitoring-setup
git checkout monitoring-setup
```

### 4. Run the Deployment Script

```bash
chmod +x scripts/deploy_monitoring.sh
./scripts/deploy_monitoring.sh
```

### 5. Verify Deployment

- [ ] API service is running and accessible (https://analabit.ru/api/health)
- [ ] Prometheus is running and accessible (https://analabit.ru/prometheus/)
- [ ] Grafana is running and accessible (https://analabit.ru/grafana/)
- [ ] Nginx configuration correctly routes to monitoring services
- [ ] Prometheus authentication is working
- [ ] Prometheus is scraping API metrics
- [ ] Grafana dashboard is properly provisioned and shows API metrics
- [ ] Grafana is running on port 3500 (not conflicting with Next.js frontend on port 3000)

### 6. Security Verification

- [ ] Prometheus is only accessible with valid credentials
- [ ] Grafana requires login with secure password (not default or weak password)
- [ ] Monitoring endpoints only accessible via HTTPS
- [ ] No sensitive credentials exposed in logs or configuration files
- [ ] Both Grafana and Prometheus credentials are generated with strong passwords
- [ ] Credentials are stored in a secure password manager for future reference

### 7. Merge to Main

If all tests pass:

```bash
# On your local machine
git checkout main
git merge monitoring-setup
git push origin main
```

### 8. Verify CI/CD Deployment

- [ ] CI/CD workflow runs successfully
- [ ] Monitoring services remain functional after CI/CD deployment
- [ ] No errors in the logs of any service

## Post-Deployment

- [ ] Document the monitoring URLs and access credentials in a secure location
- [ ] Set up alerts for critical metrics (if applicable)
- [ ] Schedule regular checks of monitoring dashboards

## Troubleshooting

### API Metrics Not Showing

- Check if API service is properly instrumented
- Verify Prometheus configuration is targeting the correct endpoint
- Check API container logs for any errors

### Nginx Configuration Issues

- Check nginx error logs: `sudo cat /var/log/nginx/error.log`
- Verify nginx configuration: `sudo nginx -t`
- Restart nginx: `sudo systemctl restart nginx`

### Authentication Issues

- Reset Prometheus credentials: `sudo htpasswd -bc /etc/nginx/.prometheus_htpasswd prometheus <new_password>`
- Reset Grafana admin password in `.env` file and restart containers

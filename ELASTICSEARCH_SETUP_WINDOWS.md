# Elasticsearch Setup on Windows Server with Docker

## 🚀 **Quick Start Guide**

### **Prerequisites**
1. Windows Server with Docker Desktop installed
2. At least 4GB RAM available for Elasticsearch
3. Administrator access to the server

## 📋 **Step-by-Step Setup**

### **Step 1: Install Docker Desktop on Windows**

1. Download Docker Desktop for Windows from: https://www.docker.com/products/docker-desktop/
2. Run the installer as Administrator
3. Enable WSL 2 integration if prompted
4. Restart the server after installation

### **Step 2: Prepare the Setup Files**

1. Create a folder on your Windows server:
   ```cmd
   mkdir C:\onetake-elasticsearch
   cd C:\onetake-elasticsearch
   ```

2. Copy the `docker-compose.elasticsearch.yml` file to this folder

### **Step 3: Start Elasticsearch**

Open PowerShell as Administrator and run:

```powershell
# Navigate to the folder
cd C:\onetake-elasticsearch

# Start Elasticsearch and Kibana
docker-compose -f docker-compose.elasticsearch.yml up -d

# Check if containers are running
docker ps
```

### **Step 4: Verify Installation**

1. **Check Elasticsearch health:**
   ```powershell
   # Test connection
   curl http://localhost:9200/_cluster/health
   
   # Or use browser: http://your-server-ip:9200
   ```

2. **Expected response:**
   ```json
   {
     "cluster_name": "onetake-logging-cluster",
     "status": "green",
     "number_of_nodes": 1
   }
   ```

3. **Access Kibana (optional):**
   - Open browser: `http://your-server-ip:5601`
   - Login: `elastic` / `OneTakeElastic2025!`

## 🔧 **Configuration Details**

### **Default Settings:**
- **Elasticsearch Port:** 9200
- **Kibana Port:** 5601 (optional dashboard)
- **Username:** `elastic`
- **Password:** `OneTakeElastic2025!`
- **Memory:** 2GB (adjust in docker-compose if needed)

### **Data Storage:**
- Data persisted in Docker volumes
- Survives container restarts
- Located in Docker's volume directory

## 🔒 **Security Configuration**

### **Change Default Password:**
```powershell
# After Elasticsearch starts, change the password
docker exec -it onetake-elasticsearch /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
```

### **Create Application User (Recommended):**
```powershell
# Create a user for your Go application
docker exec -it onetake-elasticsearch curl -X POST "localhost:9200/_security/user/onetake_logger" -H "Content-Type: application/json" -u elastic:OneTakeElastic2025! -d '{
  "password": "OneTakeLogger2025!",
  "roles": ["kibana_admin", "index_admin"],
  "full_name": "OneTake Logger Service"
}'
```

## 🌐 **Network Configuration**

### **Allow External Access:**

1. **Windows Firewall:**
   ```cmd
   # Allow Elasticsearch port
   netsh advfirewall firewall add rule name="Elasticsearch" dir=in action=allow protocol=TCP localport=9200
   
   # Allow Kibana port (optional)
   netsh advfirewall firewall add rule name="Kibana" dir=in action=allow protocol=TCP localport=5601
   ```

2. **Update your Go application .env:**
   ```bash
   LOG_LOCAL_MODE=false
   ELASTICSEARCH_URL=http://your-server-ip:9200
   ELASTICSEARCH_USER=elastic
   ELASTICSEARCH_PASS=OneTakeElastic2025!
   ```

## 📊 **Production Recommendations**

### **Memory Settings (adjust based on server):**
```yaml
# In docker-compose.elasticsearch.yml
environment:
  - "ES_JAVA_OPTS=-Xms4g -Xmx4g"  # Use 4GB instead of 2GB
```

### **Index Management:**
```powershell
# Set up index templates for automatic cleanup
docker exec -it onetake-elasticsearch curl -X PUT "localhost:9200/_index_template/onetake-logs" -H "Content-Type: application/json" -u elastic:OneTakeElastic2025! -d '{
  "index_patterns": ["api-logs-*", "error-logs-*"],
  "template": {
    "settings": {
      "index.lifecycle.name": "onetake-log-policy",
      "index.lifecycle.rollover_alias": "onetake-logs"
    }
  }
}'
```

## 🔍 **Monitoring & Maintenance**

### **Check Container Status:**
```powershell
# View running containers
docker ps

# Check logs
docker logs onetake-elasticsearch
docker logs onetake-kibana

# Check resource usage
docker stats
```

### **Backup Data:**
```powershell
# Create snapshot repository
docker exec -it onetake-elasticsearch curl -X PUT "localhost:9200/_snapshot/backup_repo" -H "Content-Type: application/json" -u elastic:OneTakeElastic2025! -d '{
  "type": "fs",
  "settings": {
    "location": "/usr/share/elasticsearch/backup"
  }
}'
```

## 🚨 **Troubleshooting**

### **Common Issues:**

1. **Container won't start:**
   ```powershell
   # Check Docker Desktop is running
   docker version
   
   # Check available memory
   docker system df
   ```

2. **Memory errors:**
   ```powershell
   # Reduce memory in docker-compose.yml
   - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
   ```

3. **Connection refused:**
   ```powershell
   # Check if port is accessible
   netstat -an | findstr :9200
   
   # Check Windows Firewall
   netsh advfirewall show allprofiles
   ```

### **Restart Services:**
```powershell
# Restart Elasticsearch
docker-compose -f docker-compose.elasticsearch.yml restart

# Or stop and start fresh
docker-compose -f docker-compose.elasticsearch.yml down
docker-compose -f docker-compose.elasticsearch.yml up -d
```

## 📈 **Performance Tuning**

### **For High-Volume Logging:**
```yaml
# Add to docker-compose.yml environment
- indices.memory.index_buffer_size=20%
- thread_pool.write.queue_size=1000
- cluster.routing.allocation.disk.watermark.low=85%
- cluster.routing.allocation.disk.watermark.high=90%
```

## ✅ **Verification Checklist**

- [ ] Docker Desktop installed and running
- [ ] Elasticsearch container started successfully
- [ ] Port 9200 accessible from your Go application server
- [ ] Authentication working with username/password
- [ ] Kibana accessible (optional)
- [ ] Firewall rules configured
- [ ] Go application .env updated with production settings

## 🔄 **Next Steps**

After Elasticsearch is running:

1. Update your Go application's `.env` file:
   ```bash
   LOG_LOCAL_MODE=false
   ELASTICSEARCH_URL=http://your-server-ip:9200
   ELASTICSEARCH_USER=elastic
   ELASTICSEARCH_PASS=OneTakeElastic2025!
   ```

2. Restart your Go application

3. Test logging by making API requests

4. Check logs in Kibana dashboard

Your centralized logging system will now be fully operational! 🎉

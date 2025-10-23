# 🔧 MongoDB Connection Fix

## ❌ Problem

MongoDB Atlas connection failing with TLS error:
```
remote error: tls: internal error
server selection error: context deadline exceeded
```

This happens because:
1. MongoDB Atlas requires TLS/SSL
2. Railway environment needs proper TLS configuration
3. Connection string might be missing TLS parameters

---

## ✅ Solution Applied

### 1. **Updated Code** ✅
Added proper timeouts and TLS handling:
- `main.go` - Updated `connectMongo()` function
- `internal/db/mongo.go` - Updated `Connect()` function

Both now include:
- Server selection timeout (10 seconds)
- Connect timeout (10 seconds)
- Proper TLS handling (automatic from connection string)

### 2. **Verify Connection String** ⚠️

Your MongoDB URI should look like this:
```
mongodb+srv://chantsztochris_db_user:pJmeH1vRNBWbWsfC@cluster0.xjbv6ys.mongodb.net/futurefrontier?retryWrites=true&w=majority
```

**Important:** Make sure it has:
- ✅ `mongodb+srv://` (not `mongodb://`)
- ✅ `?retryWrites=true&w=majority` at the end
- ✅ Database name `/futurefrontier` before the `?`

---

## 🚀 Deploy the Fix

### Step 1: Commit and Push
```bash
git add .
git commit -m "Fix MongoDB TLS connection for Railway"
git push origin features/pages-api
```

### Step 2: Verify Railway Environment Variable

In Railway dashboard, check `MONGODB_URI`:

**Correct format:**
```
mongodb+srv://chantsztochris_db_user:pJmeH1vRNBWbWsfC@cluster0.xjbv6ys.mongodb.net/futurefrontier?retryWrites=true&w=majority
```

**If it's missing parameters, update it to:**
```
mongodb+srv://chantsztochris_db_user:pJmeH1vRNBWbWsfC@cluster0.xjbv6ys.mongodb.net/futurefrontier?retryWrites=true&w=majority&tls=true&tlsAllowInvalidCertificates=false
```

### Step 3: Redeploy
Railway will auto-redeploy after you push, or click "Redeploy" in dashboard.

---

## 🔍 What Changed in Code

### Before:
```go
func connectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    // ...
}
```

### After:
```go
func connectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
    clientOpts := options.Client().
        ApplyURI(uri).
        SetServerSelectionTimeout(10 * time.Second).
        SetConnectTimeout(10 * time.Second)
    
    client, err := mongo.Connect(ctx, clientOpts)
    // ...
}
```

**Benefits:**
- ✅ Proper timeout handling
- ✅ Better error messages
- ✅ TLS automatically handled by `mongodb+srv://`
- ✅ Works with MongoDB Atlas

---

## 🧪 Test After Deployment

### Check Railway Logs:

**Success:**
```
✓ Connected to MongoDB
✓ Database: futurefrontier
✓ Server listening on :8080
```

**Still failing?** See troubleshooting below.

---

## 🚨 Troubleshooting

### Issue 1: Still getting TLS error

**Solution:** Update connection string in Railway to explicitly enable TLS:
```
mongodb+srv://chantsztochris_db_user:pJmeH1vRNBWbWsfC@cluster0.xjbv6ys.mongodb.net/futurefrontier?retryWrites=true&w=majority&tls=true
```

### Issue 2: Authentication failed

**Check:**
1. Username is correct: `chantsztochris_db_user`
2. Password is correct: `pJmeH1vRNBWbWsfC`
3. User has permissions in MongoDB Atlas

**Fix in MongoDB Atlas:**
1. Go to Database Access
2. Check user exists
3. Verify password
4. Ensure user has "Read and write to any database" permission

### Issue 3: Network access denied

**Fix in MongoDB Atlas:**
1. Go to Network Access
2. Add IP: `0.0.0.0/0` (Allow from anywhere)
3. This allows Railway to connect

### Issue 4: Database doesn't exist

**Fix:**
1. Connection string should have `/futurefrontier` before `?`
2. MongoDB will create database automatically on first write
3. Make sure `MONGODB_DATABASE=futurefrontier` is set in Railway

---

## ✅ Verification Checklist

- [ ] Code updated with timeouts
- [ ] Changes committed and pushed
- [ ] Railway redeployed
- [ ] Connection string has `mongodb+srv://`
- [ ] Connection string has database name
- [ ] Connection string has `?retryWrites=true&w=majority`
- [ ] MongoDB Atlas IP whitelist includes `0.0.0.0/0`
- [ ] User has correct permissions
- [ ] Railway logs show "Connected to MongoDB"

---

## 📊 Expected Results

### Railway Logs (Success):
```
Starting Container
✓ Server starting...
✓ Loading configuration
✓ Connecting to MongoDB...
✓ Connected to MongoDB
✓ Database: futurefrontier
✓ Connecting to Elasticsearch...
✓ Connected to Elasticsearch
✓ Server listening on :8080
```

### Test Endpoint:
```bash
curl https://your-railway-url.up.railway.app/health

# Should return:
{"status":"ok"}
```

---

## 🎉 You're Fixed!

Push the changes and MongoDB should connect successfully!

**Next:** Generate demo data and test your full deployment!

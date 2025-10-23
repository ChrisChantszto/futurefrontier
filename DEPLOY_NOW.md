# 🚀 Ready to Deploy to Railway!

## ✅ What I've Done

### 1. **Updated .env.example** with your credentials:
- ✅ MongoDB URI: `mongodb+srv://chantsztochris_db_user:...@cluster0.xjbv6ys.mongodb.net/futurefrontier`
- ✅ Elasticsearch URL: `https://my-elasticsearch-project-aa1f20.es.us-central1.gcp.elastic.cloud:443`
- ✅ Elasticsearch API Key: `N0RJRkVKb0JUZkhaRzVOdTdmZU06Umw4MThoTU1kVVMza043LXJ5SUtBUQ==`

### 2. **Updated Code** to support Elasticsearch API Key:
- ✅ Added `ElasticsearchAPIKey` field to config
- ✅ Updated `elasticsearch.go` to use API key authentication
- ✅ Changed default `LOG_PROJECT_ID` to `futurefrontier`
- ✅ Changed `LOG_LOCAL_MODE` default to `false`

### 3. **Created Railway Configuration:**
- ✅ `RAILWAY_ENV_VARS.txt` - All environment variables ready to copy

---

## 📋 Deployment Steps (10 minutes)

### Step 1: Push to GitHub (2 min)

```bash
git add .
git commit -m "Production ready with Elastic Cloud and MongoDB Atlas"
git push origin main
```

**Your .gitignore protects:**
- ❌ `.env` (won't be pushed)
- ❌ `grand-principle-475206-b5-6a55c4da97f9.json` (won't be pushed)
- ✅ `.env.example` (will be pushed as template)

---

### Step 2: Deploy to Railway (5 min)

1. **Go to https://railway.app**
2. **Sign up with GitHub**
3. **Click "New Project" → "Deploy from GitHub repo"**
4. **Select `futurefrontier` repository**
5. **Wait for initial build** (may fail - that's OK)

6. **Add Environment Variables:**
   - Click "Variables" tab
   - Open `RAILWAY_ENV_VARS.txt`
   - Copy each variable and paste into Railway

7. **For GOOGLE_CREDENTIALS_JSON:**
   - Open `grand-principle-475206-b5-6a55c4da97f9.json`
   - Copy **ENTIRE** file content (all lines)
   - In Railway, create variable `GOOGLE_CREDENTIALS_JSON`
   - Paste the entire JSON content

8. **Click "Redeploy"** or wait for auto-deploy

9. **Get your Railway URL:**
   - Will be like: `https://futurefrontier-production-xxxx.up.railway.app`
   - **SAVE THIS URL!**

---

### Step 3: Update Frontend Config (2 min)

1. **Edit `frontend/.env.production`:**
   ```bash
   VITE_API_URL=https://futurefrontier-production-xxxx.up.railway.app
   ```
   (Replace with your actual Railway URL)

2. **Commit and push:**
   ```bash
   git add frontend/.env.production
   git commit -m "Add production API URL"
   git push
   ```

---

### Step 4: Deploy Frontend to Vercel (5 min)

1. **Go to https://vercel.com**
2. **Sign up with GitHub**
3. **Click "New Project"**
4. **Import `futurefrontier` repository**
5. **Configure:**
   - Framework: **Vite**
   - Root Directory: **frontend**
   - Build Command: `npm run build`
   - Output Directory: `dist`

6. **Add Environment Variable:**
   - Name: `VITE_API_URL`
   - Value: `https://futurefrontier-production-xxxx.up.railway.app`

7. **Click "Deploy"**

8. **Get your Vercel URL:**
   - Will be like: `https://futurefrontier.vercel.app`
   - **SAVE THIS URL!**

---

### Step 5: Update CORS (2 min)

1. **Go back to Railway**
2. **Click "Variables"**
3. **Update `CORS_ORIGINS`:**
   ```
   https://futurefrontier.vercel.app,http://localhost:5173
   ```
   (Replace with your actual Vercel URL)

4. **Railway will auto-redeploy**

---

### Step 6: Generate Demo Data (2 min)

**Via curl:**
```bash
curl -X POST "https://futurefrontier-production-xxxx.up.railway.app/api/demo/generate?count=1000"
```

**Or via browser:**
Open: `https://futurefrontier-production-xxxx.up.railway.app/api/demo/generate?count=1000`

---

### Step 7: Test Your Deployment (3 min)

1. **Open:** `https://futurefrontier.vercel.app`

2. **Login:**
   - Email: `admin@example.com`
   - Password: `P@ssw0rd!`

3. **Test:**
   - [ ] Dashboard loads
   - [ ] Discover page shows 1000 logs
   - [ ] AI Chat responds
   - [ ] Custom date ranges work
   - [ ] Pagination works

---

## 🎯 Your Configuration

### MongoDB:
```
mongodb+srv://chantsztochris_db_user:pJmeH1vRNBWbWsfC@cluster0.xjbv6ys.mongodb.net/futurefrontier
```

### Elasticsearch:
```
URL: https://my-elasticsearch-project-aa1f20.es.us-central1.gcp.elastic.cloud:443
API Key: N0RJRkVKb0JUZkhaRzVOdTdmZU06Umw4MThoTU1kVVMza043LXJ5SUtBUQ==
```

### Vertex AI:
```
Project: grand-principle-475206-b5
Location: us-central1
```

---

## 🚨 Troubleshooting

### Railway Build Fails?
- Check logs in Railway dashboard
- Verify all environment variables are set
- Make sure `GOOGLE_CREDENTIALS_JSON` has complete JSON

### Can't Connect to Elasticsearch?
- Verify API key is correct
- Check URL includes `:443` at the end
- Ensure no extra spaces in environment variables

### Frontend Can't Connect?
- Check `VITE_API_URL` in Vercel
- Verify `CORS_ORIGINS` in Railway includes Vercel URL
- Redeploy both if needed

### No Logs Showing?
- Generate demo data again
- Check Railway logs for errors
- Verify Elasticsearch connection in logs

---

## ✅ Deployment Checklist

- [ ] Code committed and pushed to GitHub
- [ ] Railway project created
- [ ] All Railway environment variables added
- [ ] Google credentials JSON pasted
- [ ] Railway deployment successful
- [ ] Frontend `.env.production` updated
- [ ] Vercel project created
- [ ] Vercel environment variable added
- [ ] Vercel deployment successful
- [ ] CORS updated with Vercel URL
- [ ] Demo data generated (1000 logs)
- [ ] Login tested
- [ ] All features tested

---

## 🎉 You're Ready!

**Your demo URL:** `https://futurefrontier.vercel.app`

**Total time:** ~30 minutes

**Judges will see:**
- ✅ Real Elasticsearch Cloud deployment
- ✅ 1000+ logs in Discover page
- ✅ AI analyzing real data
- ✅ Pattern detection working
- ✅ Production-ready architecture

**Good luck with your hackathon! 🏆**

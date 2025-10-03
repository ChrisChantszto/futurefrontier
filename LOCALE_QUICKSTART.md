# Locale Management - Quick Start Guide

## What's Been Added

Your Go backend now has a complete locale management system that works seamlessly with your Next.js i18n frontend.

## Files Created/Modified

### New Files:
- `internal/models/locale.go` - Locale data models
- `internal/service/locale_service.go` - Business logic for locale management
- `internal/transport/http/handlers/locale_handler.go` - HTTP handlers
- `LOCALE_MANAGEMENT.md` - Complete documentation
- `examples/locale_examples.http` - API testing examples
- `examples/frontend-types.ts` - TypeScript types for frontend

### Modified Files:
- `internal/transport/http/routes.go` - Added locale routes
- `internal/db/indexes.go` - Added locale collection indexes
- `internal/utils/locale.go` - Enhanced normalization
- `main.go` - Updated CORS headers

## How It Works

### 1. Backend Initialization
When your backend starts, it automatically creates 3 default locales:
- English (en) - default
- Simplified Chinese (zh-hans)
- Traditional Chinese (zh-hant)

### 2. API Endpoints

**Public (no auth):**
- `GET /api/locales` - List all locales
- `GET /api/locales/config` - Get config for frontend
- `GET /api/locales/:id` - Get single locale

**Protected (requires auth):**
- `POST /api/locales` - Create new locale
- `PATCH /api/locales/:id` - Update locale
- `DELETE /api/locales/:id` - Delete locale
- `POST /api/locales/:id/set-default` - Set as default

## Quick Test

### 1. Start your backend:
```bash
go run main.go
```

### 2. Test the public endpoint:
```bash
curl http://localhost:8080/api/locales/config
```

Expected response:
```json
{
  "success": true,
  "data": {
    "locales": ["en", "zh-hans", "zh-hant"],
    "defaultLocale": "en"
  }
}
```

### 3. Add German (requires authentication):
```bash
curl -X POST http://localhost:8080/api/locales \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "de",
    "name": "German",
    "nativeName": "Deutsch",
    "isEnabled": true,
    "sortOrder": 4
  }'
```

### 4. Verify German was added:
```bash
curl http://localhost:8080/api/locales/config
```

Expected response:
```json
{
  "success": true,
  "data": {
    "locales": ["en", "zh-hans", "zh-hant", "de"],
    "defaultLocale": "en"
  }
}
```

## Frontend Integration

### Option 1: Static Config (Simple)
Keep your current `routing.ts` as-is. Manually update it when you add locales.

### Option 2: Dynamic Config (Recommended)
Fetch locale config from backend at build time or runtime.

**Example for Next.js:**

```typescript
// lib/locales.ts
export async function getLocaleConfig() {
  try {
    const res = await fetch('http://localhost:8080/api/locales/config');
    const data = await res.json();
    return data.data;
  } catch (error) {
    // Fallback
    return {
      locales: ['en', 'zh-hans', 'zh-hant'],
      defaultLocale: 'en'
    };
  }
}
```

Then in your app:
```typescript
// app/[locale]/layout.tsx
import { getLocaleConfig } from '@/lib/locales';

export async function generateStaticParams() {
  const config = await getLocaleConfig();
  return config.locales.map((locale) => ({ locale }));
}
```

## Adding a New Language (e.g., German)

### Step 1: Add via API (Backend)
```bash
curl -X POST http://localhost:8080/api/locales \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "de",
    "name": "German",
    "nativeName": "Deutsch",
    "isEnabled": true,
    "sortOrder": 4
  }'
```

### Step 2: Create Translation File (Frontend)
Create `messages/de.json` in your Next.js project:
```json
{
  "common": {
    "welcome": "Willkommen",
    "hello": "Hallo"
  }
}
```

### Step 3: Done!
If you're using dynamic config, German is now available. If using static config, update your `routing.ts`:

```typescript
export const routing = defineRouting({
  locales: ['en', 'zh-hans', 'zh-hant', 'de'], // Add 'de'
  defaultLocale: 'en',
  localePrefix: 'never'
});
```

## Common Operations

### Enable/Disable a Locale
```bash
curl -X PATCH http://localhost:8080/api/locales/de \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"isEnabled": false}'
```

### Change Default Locale
```bash
curl -X POST http://localhost:8080/api/locales/zh-hans/set-default \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Delete a Locale
```bash
curl -X DELETE http://localhost:8080/api/locales/de \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Database

Locales are stored in MongoDB in the `locales` collection:
```javascript
// MongoDB query examples
db.locales.find({ isEnabled: true })
db.locales.findOne({ isDefault: true })
db.locales.updateOne({ _id: "de" }, { $set: { isEnabled: false } })
```

## Next Steps

1. **Test the endpoints** using the examples in `examples/locale_examples.http`
2. **Integrate with frontend** using the TypeScript client in `examples/frontend-types.ts`
3. **Add more locales** as needed for your application
4. **Read full documentation** in `LOCALE_MANAGEMENT.md`

## Troubleshooting

### "Locale not found" error
- Check if locale exists: `GET /api/locales/:id`
- Verify locale code is lowercase

### Cannot delete locale
- Cannot delete default locale (set another as default first)
- Cannot delete last enabled locale

### Frontend not showing new locale
- Verify locale is enabled: `GET /api/locales/:id`
- Check translation file exists: `messages/{locale}.json`
- Restart Next.js dev server
- Clear browser cache

## Support

For detailed documentation, see `LOCALE_MANAGEMENT.md`.
For API examples, see `examples/locale_examples.http`.
For TypeScript types, see `examples/frontend-types.ts`.

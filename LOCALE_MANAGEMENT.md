# Locale Management System

This backend provides a complete locale management system that integrates with your Next.js i18n frontend.

## Features

- **Dynamic Locale Management**: Add, update, enable/disable, and delete locales from the backend
- **Default Locale Control**: Set which locale is the default
- **Frontend Config Export**: API endpoint that provides locale configuration in Next.js format
- **Automatic Initialization**: Default locales (en, zh-hans, zh-hant) are created on first run
- **Safety Checks**: Prevents deletion of default or last enabled locale

## Database Schema

### Locale Collection

```json
{
  "_id": "en",
  "code": "en",
  "name": "English",
  "nativeName": "English",
  "isDefault": true,
  "isEnabled": true,
  "sortOrder": 1,
  "createdAt": "2025-10-03T15:00:00Z",
  "updatedAt": "2025-10-03T15:00:00Z"
}
```

## API Endpoints

### Public Endpoints (No Authentication Required)

#### 1. Get All Locales
```http
GET /api/locales?enabled=true
```

**Query Parameters:**
- `enabled` (optional): If "true", returns only enabled locales

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "en",
      "code": "en",
      "name": "English",
      "nativeName": "English",
      "isDefault": true,
      "isEnabled": true,
      "sortOrder": 1,
      "createdAt": "2025-10-03T15:00:00Z",
      "updatedAt": "2025-10-03T15:00:00Z"
    }
  ]
}
```

#### 2. Get Locale Config (for Frontend)
```http
GET /api/locales/config
```

**Response:**
```json
{
  "success": true,
  "data": {
    "locales": ["en", "zh-hans", "zh-hant"],
    "defaultLocale": "en"
  }
}
```

This endpoint returns the exact format needed for your Next.js `routing.ts` file.

#### 3. Get Single Locale
```http
GET /api/locales/:id
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "en",
    "code": "en",
    "name": "English",
    "nativeName": "English",
    "isDefault": true,
    "isEnabled": true,
    "sortOrder": 1,
    "createdAt": "2025-10-03T15:00:00Z",
    "updatedAt": "2025-10-03T15:00:00Z"
  }
}
```

### Protected Endpoints (Authentication Required)

#### 4. Create New Locale
```http
POST /api/locales
```

**Request Body:**
```json
{
  "code": "de",
  "name": "German",
  "nativeName": "Deutsch",
  "isEnabled": true,
  "sortOrder": 4
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "de",
    "code": "de",
    "name": "German",
    "nativeName": "Deutsch",
    "isDefault": false,
    "isEnabled": true,
    "sortOrder": 4,
    "createdAt": "2025-10-03T15:00:00Z",
    "updatedAt": "2025-10-03T15:00:00Z"
  },
  "message": "Locale created successfully"
}
```

#### 5. Update Locale
```http
PATCH /api/locales/:id
```

**Request Body (all fields optional):**
```json
{
  "name": "German (Updated)",
  "nativeName": "Deutsch",
  "isEnabled": false,
  "sortOrder": 5
}
```

**Response:**
```json
{
  "success": true,
  "data": { /* updated locale */ },
  "message": "Locale updated successfully"
}
```

#### 6. Set Default Locale
```http
POST /api/locales/:id/set-default
```

This will unset all other defaults and set the specified locale as default.

**Response:**
```json
{
  "success": true,
  "message": "Default locale set successfully"
}
```

#### 7. Delete Locale
```http
DELETE /api/locales/:id
```

**Safety Checks:**
- Cannot delete the default locale
- Cannot delete the last enabled locale

**Response:**
```json
{
  "success": true,
  "message": "Locale deleted successfully"
}
```

## Frontend Integration

### Step 1: Fetch Locale Config from Backend

Update your Next.js `routing.ts` to fetch from the backend:

```typescript
// routing.ts
import {defineRouting} from 'next-intl/routing';

// Fetch locale config from backend
async function getLocaleConfig() {
  try {
    const response = await fetch('http://localhost:8080/api/locales/config');
    const data = await response.json();
    return data.data; // { locales: [...], defaultLocale: "..." }
  } catch (error) {
    // Fallback to hardcoded values
    return {
      locales: ['en', 'zh-hans', 'zh-hant'],
      defaultLocale: 'en'
    };
  }
}

// For build time, you can use a static config or environment variable
export const routing = defineRouting({
  locales: ['en', 'zh-hans', 'zh-hant'], // Will be overridden at runtime
  defaultLocale: 'en',
  localePrefix: 'never'
});
```

### Step 2: Dynamic Locale Loading (Optional)

For runtime locale switching, create a client-side hook:

```typescript
// hooks/useLocales.ts
import { useState, useEffect } from 'react';

export function useLocales() {
  const [locales, setLocales] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/locales?enabled=true')
      .then(res => res.json())
      .then(data => {
        setLocales(data.data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to load locales:', err);
        setLoading(false);
      });
  }, []);

  return { locales, loading };
}
```

### Step 3: Locale Switcher Component

```typescript
// components/LocaleSwitcher.tsx
'use client';

import { useLocales } from '@/hooks/useLocales';
import { useRouter, usePathname } from '@/i18n/navigation';
import { useLocale } from 'next-intl';

export function LocaleSwitcher() {
  const { locales, loading } = useLocales();
  const currentLocale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  const handleChange = (newLocale: string) => {
    router.replace(pathname, { locale: newLocale });
  };

  if (loading) return <div>Loading...</div>;

  return (
    <select value={currentLocale} onChange={(e) => handleChange(e.target.value)}>
      {locales.map((locale) => (
        <option key={locale.code} value={locale.code}>
          {locale.nativeName}
        </option>
      ))}
    </select>
  );
}
```

## Adding a New Locale (e.g., German)

### Backend Steps:

1. **Create the locale via API:**
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

2. **The locale is now available!** The frontend will automatically pick it up from `/api/locales/config`.

### Frontend Steps:

1. **Create translation file:**
   - Add `messages/de.json` with your German translations

2. **No code changes needed!** The locale will be automatically available if you're fetching from the backend.

## Example Workflow: Adding German

```bash
# 1. Add German locale (authenticated request)
POST /api/locales
{
  "code": "de",
  "name": "German",
  "nativeName": "Deutsch",
  "isEnabled": true,
  "sortOrder": 4
}

# 2. Verify it's available
GET /api/locales/config
# Response: { "locales": ["en", "zh-hans", "zh-hant", "de"], "defaultLocale": "en" }

# 3. Create frontend translation file
# messages/de.json

# 4. Done! German is now available in your app
```

## Locale Normalization

The backend includes locale normalization utilities in `internal/utils/locale.go`:

- `zh`, `zh-CN`, `zh-SG` → `zh-hans`
- `zh-TW`, `zh-HK`, `zh-MO` → `zh-hant`
- `en-US`, `en-GB` → `en`
- `de-DE`, `de-AT` → `de`

This ensures consistent locale codes across your application.

## Best Practices

1. **Always have at least one enabled locale**
2. **Always have exactly one default locale**
3. **Use lowercase locale codes** (e.g., "de" not "DE")
4. **Use hyphens for variants** (e.g., "zh-hans" not "zh_hans")
5. **Create translation files before enabling locales** in production
6. **Test locale switching** before deploying new locales

## Troubleshooting

### Locale not showing in frontend?
- Check if locale is enabled: `GET /api/locales/:id`
- Verify translation file exists: `messages/{locale}.json`
- Clear browser cache and reload

### Cannot delete locale?
- Check if it's the default locale (cannot delete)
- Check if it's the last enabled locale (cannot delete)
- Set another locale as default first, then delete

### Locale config not updating?
- Frontend might be caching the config
- Restart your Next.js dev server
- Check CORS settings if calling from different domain

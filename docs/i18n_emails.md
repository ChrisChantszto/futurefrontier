# Internationalized (i18n) Email System

This document describes the internationalization (i18n) support for emails in the Onetake Corpsite Backend.

## Overview

The system supports sending emails in multiple languages based on the user's locale preference. Currently, the following locales are supported:

- English (`en`) - Default
- Simplified Chinese (`zh-hans`)
- Traditional Chinese (`zh-hant`)

## How It Works

### 1. Configuration

The i18n system is configured through environment variables:

```
I18N_DEFAULT_LOCALE=en
I18N_SUPPORTED_LOCALES=en,zh-hans,zh-hant
I18N_TEMPLATES_DIR=./templates
```

### 2. Email Templates

Email templates are organized by locale in the templates directory:

```
templates/
├── emails/
│   ├── en/
│   │   └── otp.html
│   ├── zh-hans/
│   │   └── otp.html
│   └── zh-hant/
│       └── otp.html
```

### 3. Locale Detection

The system determines the locale to use in the following order:

1. Explicit locale parameter in the API request body
2. Accept-Language HTTP header
3. Default locale from configuration

### 4. API Usage

When requesting an OTP, you can specify the locale:

```json
POST /auth/otp/request
{
  "email": "user@example.com",
  "purpose": "login",
  "locale": "zh-hans"
}
```

The system also supports common locale variations:
- `zh`, `zh-CN`, `zh_cn` → `zh-hans`
- `zh-TW`, `zh_tw` → `zh-hant`

### 5. Fallback Mechanism

If a template is not available in the requested locale, the system falls back to the default locale (English).

## Adding New Languages

To add support for a new language:

1. Add the locale code to the `I18N_SUPPORTED_LOCALES` environment variable
2. Create a new directory under `templates/emails/` with the locale code
3. Add email templates for the new locale
4. Update the subject translations in `smtp.go`

## Testing

You can test the i18n email system using the provided test script:

```
go run test_i18n_otp.go
```

This will send test emails in different locales to verify the implementation.

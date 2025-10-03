/**
 * TypeScript types for the Locale Management API
 * Copy these to your Next.js frontend project
 */

// Locale model
export interface Locale {
  id: string;
  code: string;
  name: string;
  nativeName: string;
  isDefault: boolean;
  isEnabled: boolean;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

// Locale config for routing
export interface LocaleConfig {
  locales: string[];
  defaultLocale: string;
}

// API Response wrapper
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
}

// Request types
export interface CreateLocaleRequest {
  code: string;
  name: string;
  nativeName: string;
  isEnabled?: boolean;
  sortOrder?: number;
}

export interface UpdateLocaleRequest {
  name?: string;
  nativeName?: string;
  isEnabled?: boolean;
  isDefault?: boolean;
  sortOrder?: number;
}

// API Client Example
export class LocaleApiClient {
  constructor(private baseUrl: string = 'http://localhost:8080/api') {}

  async getLocales(enabledOnly: boolean = false): Promise<Locale[]> {
    const url = enabledOnly 
      ? `${this.baseUrl}/locales?enabled=true`
      : `${this.baseUrl}/locales`;
    
    const response = await fetch(url);
    const data: ApiResponse<Locale[]> = await response.json();
    return data.data || [];
  }

  async getLocaleConfig(): Promise<LocaleConfig> {
    const response = await fetch(`${this.baseUrl}/locales/config`);
    const data: ApiResponse<LocaleConfig> = await response.json();
    return data.data || { locales: ['en'], defaultLocale: 'en' };
  }

  async getLocale(id: string): Promise<Locale | null> {
    const response = await fetch(`${this.baseUrl}/locales/${id}`);
    if (!response.ok) return null;
    const data: ApiResponse<Locale> = await response.json();
    return data.data || null;
  }

  async createLocale(
    request: CreateLocaleRequest,
    token: string
  ): Promise<Locale> {
    const response = await fetch(`${this.baseUrl}/locales`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(request),
    });
    const data: ApiResponse<Locale> = await response.json();
    if (!response.ok) throw new Error(data.message || 'Failed to create locale');
    return data.data!;
  }

  async updateLocale(
    id: string,
    request: UpdateLocaleRequest,
    token: string
  ): Promise<Locale> {
    const response = await fetch(`${this.baseUrl}/locales/${id}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(request),
    });
    const data: ApiResponse<Locale> = await response.json();
    if (!response.ok) throw new Error(data.message || 'Failed to update locale');
    return data.data!;
  }

  async deleteLocale(id: string, token: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/locales/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      const data: ApiResponse<never> = await response.json();
      throw new Error(data.message || 'Failed to delete locale');
    }
  }

  async setDefaultLocale(id: string, token: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/locales/${id}/set-default`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      const data: ApiResponse<never> = await response.json();
      throw new Error(data.message || 'Failed to set default locale');
    }
  }
}

// Usage example:
// const client = new LocaleApiClient('http://localhost:8080/api');
// const locales = await client.getLocales(true);
// const config = await client.getLocaleConfig();

const DEFAULT_BASE_URL = 'http://127.0.0.1:8080/api';

export function getApiBaseUrl() {
  const envUrl = import.meta.env.VITE_API_BASE_URL;
  return (envUrl && envUrl.trim()) || DEFAULT_BASE_URL;
}

export function getStoredToken() {
  return localStorage.getItem('ti_token') || '';
}

export function setStoredToken(token) {
  if (token) {
    localStorage.setItem('ti_token', token);
  } else {
    localStorage.removeItem('ti_token');
  }
}

export function extractList(payload) {
  if (Array.isArray(payload)) return payload;
  if (Array.isArray(payload?.data)) return payload.data;
  if (Array.isArray(payload?.data?.data)) return payload.data.data;
  if (Array.isArray(payload?.items)) return payload.items;
  if (Array.isArray(payload?.data?.items)) return payload.data.items;
  return [];
}

export function extractToken(payload) {
  return (
    payload?.token ||
    payload?.access_token ||
    payload?.data?.token ||
    payload?.data?.access_token ||
    ''
  );
}

export async function apiRequest(path, { method = 'GET', body, token } = {}) {
  const baseUrl = getApiBaseUrl().replace(/\/+$/, '');
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  const headers = {
    'Content-Type': 'application/json',
  };

  const activeToken = token !== undefined ? token : getStoredToken();
  if (activeToken) {
    headers.Authorization = `Bearer ${activeToken}`;
  }

  const response = await fetch(`${baseUrl}${cleanPath}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  let data = null;
  const text = await response.text();
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text };
    }
  }

  if (!response.ok) {
    const detailErrors =
      data?.errors && typeof data.errors === 'object'
        ? Object.values(data.errors).join(', ')
        : '';
    const message =
      detailErrors ||
      data?.message ||
      data?.error ||
      `Request gagal dengan status ${response.status}`;
    const err = new Error(message);
    err.status = response.status;
    err.payload = data;
    throw err;
  }

  return data;
}

export const api = {
  health: () => apiRequest('/health'),
  login: (credentials) => apiRequest('/auth/login', { method: 'POST', body: credentials }),
  register: (payload) => apiRequest('/auth/register', { method: 'POST', body: payload }),

  listItineraries: () => apiRequest('/itineraries'),
  createItinerary: (payload) => apiRequest('/itineraries', { method: 'POST', body: payload }),
  updateItinerary: (id, payload) => apiRequest(`/itineraries/${id}`, { method: 'PUT', body: payload }),
  deleteItinerary: (id) => apiRequest(`/itineraries/${id}`, { method: 'DELETE' }),

  listSpots: (itineraryId) =>
    apiRequest(itineraryId ? `/spots?itinerary_id=${itineraryId}` : '/spots'),
  createSpot: (payload) => apiRequest('/spots', { method: 'POST', body: payload }),
  updateSpot: (id, payload) => apiRequest(`/spots/${id}`, { method: 'PUT', body: payload }),
  deleteSpot: (id) => apiRequest(`/spots/${id}`, { method: 'DELETE' }),

  listSpotNotes: (spotId) =>
    apiRequest(spotId ? `/spot-notes?spot_id=${spotId}` : '/spot-notes'),
  createSpotNote: (payload) => apiRequest('/spot-notes', { method: 'POST', body: payload }),
  updateSpotNote: (id, payload) => apiRequest(`/spot-notes/${id}`, { method: 'PUT', body: payload }),
  deleteSpotNote: (id) => apiRequest(`/spot-notes/${id}`, { method: 'DELETE' }),
};

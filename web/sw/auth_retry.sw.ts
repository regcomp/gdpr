/// <reference lib="webworker" />

import { AUTH_RETRY_SW_CONSTANTS } from "../generated/auth_retry.sw.constants.ts";

const CONSTANTS = AUTH_RETRY_SW_CONSTANTS;

declare const self: ServiceWorkerGlobalScope;

self.addEventListener('activate', (event: ExtendableEvent) => {
  event.waitUntil(self.clients.claim());
});

self.addEventListener('install', () => {
  console.log('Service worker installing...');
  // Force the waiting service worker to become the active service worker
  self.skipWaiting();
});

self.addEventListener('fetch', (event: FetchEvent) => {
  event.respondWith(handleFetchWithAuth(event.request));
});

async function handleFetchWithAuth(request: Request): Promise<Response> {
  try {
    const response = await fetch(addHeaderAndClone(request));

    if (response.status === 401) {
      if (response.headers.get(CONSTANTS.RENEW_TOKEN_HEADER)) {
        const isRefreshSuccess = await refreshToken();

        if (isRefreshSuccess) {
          return await fetch(addHeaderAndClone(request));
        }
      }
    }

    return response;

  } catch (error) {
    console.error('Fetch error in service worker:', error);
    throw error;
  }
}

async function refreshToken(): Promise<boolean> {
  try {
    const refreshResponse = await fetch(CONSTANTS.RENEW_TOKEN_PATH, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        [CONSTANTS.RETRY_STATUS_HEADER]: CONSTANTS.TRUE,
      }
    });

    return refreshResponse.ok;
  } catch (error) {
    console.error('Token refresh error:', error);
    return false;
  }
}

function addHeaderAndClone(request: Request): Request {
  let headers = new Headers(request.headers);
  headers.set(CONSTANTS.RETRY_STATUS_HEADER, CONSTANTS.TRUE);

  if (request.mode === "navigate") {
    return new Request(request.url, {
      method: request.method,
      headers: headers,
      credentials: "include",
      redirect: request.redirect,
      cache: request.cache,
      referrer: request.referrer,
      referrerPolicy: request.referrerPolicy,
      integrity: request.integrity,
      keepalive: request.keepalive
    });
  }

  return new Request(request, {
    headers: headers,
    credentials: "include"
  });
}

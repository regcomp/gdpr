importScripts("/static/js/shared/auth_retry_worker_constants.js")

const CONSTANTS = AUTH_RETRY_WORKER_CONSTANTS

self.addEventListener('activate', event => {
  event.waitUntil(clients.claim());
});

self.addEventListener('install', event => {
  console.log('Service worker installing...');
  // Force the waiting service worker to become the active service worker
  self.skipWaiting();
});

self.addEventListener('fetch', (event) => {
  event.respondWith(handleFetchWithAuth(event.request));
});

async function handleFetchWithAuth(request) {
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

    return response

  } catch (error) {
    console.error('Fetch error in service worker:', error);
    throw error;
  }
}

async function refreshToken() {
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

function addHeaderAndClone(request) {
  const headers = new Headers(request.headers);
  headers.set(CONSTANTS.RETRY_STATUS_HEADER, CONSTANTS.TRUE)
  return new Request(request, {
    headers: headers,
    credentials: 'include',
  })
}

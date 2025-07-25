importScripts("/static/js/shared.js")

self.addEventListener('activate', event => {
  event.waitUntil(clients.claim());
});

self.addEventListener('fetch', (event) => {
  event.respondWith(handleFetchWithAuth(event.request));
});

async function handleFetchWithAuth(request) {
  try {
    const response = await fetch(addHeaderAndClone(request));

    if (response.status === 401) {
      if (response.headers.get(HEADERS.SERVICE_WORKER_ALLOWED)) {
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
    const refreshResponse = await fetch('/auth/refresh', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        [HEADERS.AUTH_RETRY_WORKER_RUNNING]: 'true'
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
  headers.set('SW-Auth-Retry-Running', 'true')
  return new Request(request, {
    headers: headers,
    credentials: 'include',
  })
}

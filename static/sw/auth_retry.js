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
      if (response.headers.get(SHARED.HEADERS.RENEW_ACCESS_TOKEN)) {
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
    const refreshResponse = await fetch(SHARED.PATHS.AUTH_RENEW, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        [SHARED.HEADERS.AUTH_RETRY_WORKER_RUNNING]: SHARED.VALUES.TRUE,
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
  headers.set(SHARED.HEADERS.AUTH_RETRY_WORKER_RUNNING, SHARED.VALUES.TRUE)
  return new Request(request, {
    headers: headers,
    credentials: 'include',
  })
}

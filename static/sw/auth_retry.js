self.addEventListener('fetch', (event) => {
  event.respondWith(handleFetchWithAuth(event.request));
});

async function handleFetchWithAuth(request) {
  const clonedRequest = request.clone();
  const headers = new Headers(clonedRequest.headers);
  headers.set('SW-Auth-Retry-Running', true);

  const modifiedRequest = new Request(clonedRequest, {
    headers: headers
  });

  try {
    const response = await fetch(modifiedRequest);

    if (response.status === 401) {
      console.log("handling 401");
      return handle401(response, modifiedRequest);
    }

    return response;

  } catch (error) {
    console.error('Fetch error in service worker:', error);
    throw error;
  }
}

async function handle401(oldResponse, request) {
  if (response.headers.get('Refresh-Access-Token')) {
    console.log('401 with Refresh-Access-Token detected, attempting refresh...');

    const refreshSuccess = await refreshToken();

    if (!refreshSuccess) {
      await deauthorize();
      return response;
    }
    return await fetch(request.clone());
  }
  return oldResponse;
}

async function refreshToken() {
  try {
    const refreshResponse = await fetch('/auth/refresh/', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'SW-Auth-Retry-Running': true
      }
    });

    return refreshResponse.ok;
  } catch (error) {
    console.error('Token refresh error:', error);
    return false;
  }
}

async function deauthorize() {
  try {
    await fetch('/deauthorize', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'SW-Auth-Retry-Running': true
      }
    });
  } catch (error) {
    console.error('Deauthorize error:', error);
  }
}

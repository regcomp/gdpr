// Service Worker for automatic auth token refresh
self.addEventListener('fetch', (event) => {
  event.respondWith(intercepted(event.request));
});

async function intercepted(request) {
  // add intercepted header

  const response = await fetch(request)
  return response
}


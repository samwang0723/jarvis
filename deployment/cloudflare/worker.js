export default {
  async fetch(request, env, ctx) {
    async function createCacheKey(request) {
      // Clone the request to access its body
      const requestBody = await request.clone().text();
    
      // Create a new URL object to modify search params
      const url = new URL(request.url);
    
      // Append body to the URL search params to include it in the cache key
      url.searchParams.append('body', requestBody);
    
      // Return a new Request object with the modified URL, which will be used as the cache key
      return new Request(url, {
        method: 'GET',
        headers: {
          'Content-Type': request.headers.get('Content-Type') // Preserve the Content-Type
        }
      });
    }

    try {
      // Cache duration in seconds
      // const CACHE_DURATION = 10;

      // Define the endpoints to cache
      const cacheUrls = [
        'https://api.jarvis-stock.tw/v1/dailycloses',
        'https://api.jarvis-stock.tw/v1/selections'
      ];

      // Check if the request method is POST and the URL is one of the endpoints
      if (request.method.toUpperCase() === 'POST' && cacheUrls.includes(new URL(request.url).href)) {
        // Create a cache key based on the request URL and body
        const cacheKey = await createCacheKey(request);
        const cache = caches.default;

        // Try to find the cache match for the request
        let response = await cache.match(cacheKey);
        if (!response) {
          // If not in cache, fetch from the server
          response = await fetch(request, {
            cf: {
              // Always cache this fetch regardless of content type
              // for a max of 5 seconds before revalidating the resource
              cacheTtl: 10,
              cacheEverything: true,
            },
          });
          // Reconstruct the Response object to make its headers mutable.
          response = new Response(response.body, response);
          // Set cache control headers to cache on browser for 10 seconds
          response.headers.set("Cache-Control", "max-age=10");

          // Cache the cloned response with a TTL of 10 seconds
          ctx.waitUntil(cache.put(cacheKey, response.clone()));
        }

        // Return the cached response or the server response
        return response;
      } else {
        // If the request is not a POST or not for the specified endpoints, fetch as usual
        return fetch(request);
      }
    } catch (e) {
      return new Response("Error thrown " + e.message);
    }
  },
};

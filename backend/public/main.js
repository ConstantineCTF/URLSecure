// /assets/main.js

document.addEventListener('DOMContentLoaded', () => {
  const form     = document.getElementById('shorten-form');
  const input    = document.getElementById('url-input');
  const resultEl = document.getElementById('result');

  // Helper to display messages
  function showMessage(message, type = 'info') {
    resultEl.textContent = message;
    resultEl.classList.remove('text-red-600', 'text-green-600', 'text-indigo-600');
    if (type === 'error') {
      resultEl.classList.add('text-red-600');
    } else if (type === 'success') {
      resultEl.classList.add('text-green-600');
    } else {
      resultEl.classList.add('text-indigo-600');
    }
  }

  form.addEventListener('submit', async e => {
    e.preventDefault();
    const url = input.value.trim();
    if (!url) {
      showMessage('Please enter a URL to shorten.', 'error');
      return;
    }

    // Retrieve JWT from localStorage
    const token = localStorage.getItem('token');
    if (!token) {
      showMessage('You must be logged in to shorten URLs.', 'error');
      setTimeout(() => window.location.href = '/assets/login.html', 1500);
      return;
    }

    showMessage('Shortening...', 'info');
    form.querySelector('button').setAttribute('disabled', 'disabled');

    try {
      const response = await fetch('/api/shorten', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ url })
      });

      const data = await response.json();

      if (response.status === 401) {
        showMessage('Session expired. Redirecting to login...', 'error');
        localStorage.removeItem('token');
        setTimeout(() => window.location.href = '/assets/login.html', 1500);
      } else if (response.ok && data.code) {
        const shortUrl = `${window.location.origin}/r/${data.code}`;
        resultEl.innerHTML = `<a href="${shortUrl}" target="_blank" class="text-indigo-600 hover:underline">${shortUrl}</a>`;
        resultEl.classList.remove('text-red-600');
        resultEl.classList.add('text-green-600');
      } else {
        throw new Error(data.error || 'Could not shorten URL');
      }
    } catch (err) {
      showMessage(`Error: ${err.message}`, 'error');
    } finally {
      form.querySelector('button').removeAttribute('disabled');
    }
  });
});

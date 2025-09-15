// /assets/main.js

document.addEventListener('DOMContentLoaded', () => {
  const form     = document.getElementById('shorten-form');
  const input    = document.getElementById('url-input');
  const resultEl = document.getElementById('result');

  if (!form || !input || !resultEl) return;

  function showMessage(message, type = 'info') {
    resultEl.textContent = message;
    resultEl.classList.remove('text-red-600', 'text-green-600', 'text-indigo-600');
    resultEl.classList.add(type === 'error' ? 'text-red-600' : type === 'success' ? 'text-green-600' : 'text-indigo-600');
  }

  form.addEventListener('submit', async e => {
    e.preventDefault();
    const url = input.value.trim();
    if (!url) return showMessage('Please enter a URL to shorten.', 'error');

    const token = localStorage.getItem('token');
    if (!token) {
      showMessage('You must be logged in to shorten URLs.', 'error');
      return setTimeout(() => (window.location.href = '/assets/login.html'), 1500);
    }

    showMessage('Shortening...', 'info');
    const btn = form.querySelector('button[type="submit"]');
    btn?.setAttribute('disabled', '');

    try {
      const res = await fetch('/api/shorten', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ url })
      });
      const data = await res.json();
      console.log('shorten response', data);

      if (!data.code) {
        return showMessage(`Error: ${data.error || 'No code returned'}`, 'error');
      }

      const shortUrl = `${window.location.origin}/r/${data.code}`;
      resultEl.innerHTML = `<a href="${shortUrl}" target="_blank" class="text-indigo-600 hover:underline">${shortUrl}</a>`;
      showMessage('', 'success');
    } catch (err) {
      showMessage(`Error: ${err.message}`, 'error');
    } finally {
      btn?.removeAttribute('disabled');
    }
  });
});

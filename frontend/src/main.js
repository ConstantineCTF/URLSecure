const resultEl = document.getElementById("result");
const form = document.getElementById("shorten-form");
const input = document.getElementById("url-input");

form.onsubmit = async e => {
  e.preventDefault();
  resultEl.textContent = "Shortening...";
  try {
    const res = await fetch("http://localhost:8080/api/shorten", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: input.value }),
    });
    const data = await res.json();
    if (data.code) {
      const link = `${window.location.origin}/r/${data.code}`;
      resultEl.innerHTML = `<a href="${link}" target="_blank" class="text-blue-600 underline">${link}</a>`;
    } else {
      throw new Error(data.error || "Unexpected response");
    }
  } catch (err) {
    resultEl.textContent = `Error: ${err.message}`;
    resultEl.classList.add("text-red-600");
  }
};

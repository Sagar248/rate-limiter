const API_URL = "https://rate-limiter-production.up.railway.app/api";

function getTime() {
  return new Date().toLocaleTimeString();
}

async function hitApi() {
  const userId = document.getElementById("userId").value || "demo";

  try {
    const res = await fetch(API_URL, {
      headers: { "X-User-ID": userId }
    });

    const data = await res.json();

    if (data.error) {
      log(`❌ Rate limit exceeded`, "error");
    } else {
      log(`✅ Allowed | Remaining: ${data.remaining}`, "success");
    }
  } catch (err) {
    log(`⚠️ Network error`, "error");
  }
}

async function burst() {
  for (let i = 0; i < 10; i++) {
    await hitApi();
  }
}

function log(message, type) {
  const logDiv = document.getElementById("log");

  const item = document.createElement("div");
  item.className = `log-item ${type}`;

  item.innerHTML = `
    ${message}
    <div class="meta">${getTime()}</div>
  `;

  logDiv.prepend(item);
}

function clearLogs() {
  document.getElementById("log").innerHTML = "";
}